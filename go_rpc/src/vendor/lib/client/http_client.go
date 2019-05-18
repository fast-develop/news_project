package client

import (
	"bytes"
	"fmt"
	log "google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/grpclog/glogger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
	"lib/config"
)

const (
	HttpFormatJson     = "json"
	HttpFormatProtoBuf = "protobuf"
)

const (
	ContentTypeJson     = "application/json;charset=utf-8"
	ContentTypeProtoBuf = "application/x-protobuf"
	ContentTypeRawData  = "application/rawdata"
	ContentTypeForm     = "application/x-www-form-urlencoded"
)

const HttpDefaultTimeOut = 1000
const HttpDefaultIdleConn = 100
const HttpDefaultIdleKeepTime = 30

const AllowReqErrDistance = 0.005

const AutoAdjustLowestReqCount = uint32(100)

const WeightGrowStep = 10
const WeightGrowPercent = 0.2

const MinAdjustPeriod = 1
const MaxAdjustPeriod = 60

const MinOverload = 1.2
const MaxOverload = 2.0

const DefaultMaxWeight = 1000

const RangeMinMaxWeight = 100
const RangeMaxMaxWeight = 10000

type ClientReqCount struct {
	totalReq uint32
	failReq  uint32
}

type HttpClient struct {
	config            *config.ClientInfo
	httpclient        *http.Client
	clientWeight      []int
	clientReqCount    []ClientReqCount
	addrLen           int
	adjustPeriod      int
	clientMaxWeight   int
	clientMaxOverload float32
}

func NewHttpClient(cinfo *config.ClientInfo) *HttpClient {
	httpTimeOut := HttpDefaultTimeOut
	httpIdleConn := HttpDefaultIdleConn
	httpIdleKeepTime := HttpDefaultIdleKeepTime

	if cinfo.Timeout > 0 {
		httpTimeOut = cinfo.Timeout
	}

	if cinfo.MaxIdleConn > 0 {
		httpIdleConn = cinfo.MaxIdleConn
	}

	client := &HttpClient{config: cinfo,
		httpclient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        httpIdleConn,
				MaxIdleConnsPerHost: httpIdleConn,
				IdleConnTimeout:     time.Second * time.Duration(httpIdleKeepTime),
			},
			Timeout: time.Duration(httpTimeOut) * time.Millisecond,
		},
	}

	if client.config.WeightMax == 0 {
		client.clientMaxWeight = DefaultMaxWeight
	} else if client.config.WeightMax < RangeMinMaxWeight {
		client.clientMaxWeight = RangeMinMaxWeight
	} else if client.config.WeightMax > RangeMaxMaxWeight {
		client.clientMaxWeight = RangeMaxMaxWeight
	} else {
		client.clientMaxWeight = client.config.WeightMax
	}

	client.addrLen = len(client.config.Addrs)
	client.clientWeight = make([]int, client.addrLen)
	client.clientReqCount = make([]ClientReqCount, client.addrLen)

	for i := 0; i < client.addrLen; i++ {
		client.clientWeight[i] = client.clientMaxWeight
	}

	if client.config.WeightMaxOverload < MinOverload {
		client.clientMaxOverload = MinOverload
	} else if client.config.WeightMaxOverload > MaxOverload {
		client.clientMaxOverload = MaxOverload
	} else {
		client.clientMaxOverload = client.config.WeightMaxOverload
	}

	if client.config.WeightAdjustPeriod < MinAdjustPeriod {
		client.adjustPeriod = MinAdjustPeriod
	} else if client.config.WeightAdjustPeriod > MaxAdjustPeriod {
		client.adjustPeriod = MaxAdjustPeriod
	} else {
		client.adjustPeriod = client.config.WeightAdjustPeriod
	}

	if client.config.WeightAuto != 0 {
		go client.AutoAdjustWeight()
	}

	fmt.Printf("client config:%v\n", client)

	return client
}

func (c *HttpClient) AutoAdjustWeight() {
	nullCount := make([]ClientReqCount, c.addrLen)

	for {
		time.Sleep(time.Duration(c.adjustPeriod) * time.Second)

		oldCount := make([]ClientReqCount, c.addrLen)
		perAddrErrRate := make([]float64, c.addrLen)
		tmpWeight := make([]int, c.addrLen)

		copy(oldCount, c.clientReqCount)

		allReq := uint32(0)
		allFailReq := uint32(0)
		for index, val := range oldCount {
			allReq += val.totalReq
			allFailReq += val.failReq
			if val.totalReq != 0 {
				perAddrErrRate[index] = float64(val.failReq) / float64(val.totalReq)
			} else {
				perAddrErrRate[index] = -1.0
			}
		}

		/* 统计样本不够, 增加等待时间 */
		if allReq < AutoAdjustLowestReqCount {
			continue
		}

		copy(c.clientReqCount, nullCount)

		allErrRate := float64(allFailReq) / float64(allReq)

		log.Infof("client:%s AutoAdjustWeight: allErrRate: %f node err: %v", c.config.Name, allErrRate, perAddrErrRate)

		allWeightCount := 0
		/* 出错率大于平均，则减少权重，否则增加权重 */
		for index, errRate := range perAddrErrRate {
			adjustWeigth := c.clientWeight[index]
			if errRate < 0 {
				adjustWeigth = c.clientWeight[index] + WeightGrowStep
			} else {
				err_diff := errRate - allErrRate
				if err_diff > AllowReqErrDistance {
					adjustWeigth = c.clientWeight[index] - int(float64(c.clientWeight[index])*err_diff)

				} else if c.clientWeight[index] < c.clientMaxWeight {
					adjustWeigth = c.clientWeight[index] + int(float64(c.clientMaxWeight-c.clientWeight[index])*WeightGrowPercent) + WeightGrowStep
				}
			}

			if adjustWeigth < 1 {
				adjustWeigth = 1
			} else if adjustWeigth > c.clientMaxWeight {
				adjustWeigth = c.clientMaxWeight
			}

			tmpWeight[index] = adjustWeigth
			allWeightCount += adjustWeigth
		}

		avgWeight := allWeightCount / c.addrLen
		overloadWeight := int(float32(avgWeight) * c.clientMaxOverload)
		moveWeight := 0

		/* 避免某个节点过载，调整流量不超过最大过载限制 */
		for index, w := range tmpWeight {
			if w > overloadWeight {
				moveWeight += w - overloadWeight
				tmpWeight[index] = overloadWeight
			}
		}

		if moveWeight > 0 {
			perAddrMoveWeight := moveWeight / c.addrLen / 10
			if perAddrMoveWeight < 1 {
				perAddrMoveWeight = 1
			}
			for {
				beforeWeight := moveWeight
				for index, w := range tmpWeight {
					if w < overloadWeight {
						if w+perAddrMoveWeight <= overloadWeight {
							tmpWeight[index] = w + perAddrMoveWeight
							moveWeight -= perAddrMoveWeight
						} else {
							tmpWeight[index] = overloadWeight
							moveWeight -= overloadWeight - w
						}
					}

					if moveWeight <= 0 {
						break
					}

					if moveWeight < perAddrMoveWeight {
						perAddrMoveWeight = moveWeight
					}
				}

				if moveWeight <= 0 {
					break
				}

				if moveWeight == beforeWeight {
					log.Warningf("adjust overload maybe err, last weight: %d", moveWeight)
					break
				}
			}
		}

		log.Infof("client:%s AutoAdjustWeight: old weight: %v new weight: %v", c.config.Name, c.clientWeight, tmpWeight)

		copy(c.clientWeight, tmpWeight)
	}
}

func (c *HttpClient) getAddr(already_used []int) int {
	allWeight := 0
	for i, w := range c.clientWeight {
		if already_used[i] == 0 {
			allWeight += w
		}
	}
	if allWeight < 1 {
		log.Errorf("HttpClient:getAddr something err, allWeight < 1")
		rd := rand.Intn(c.addrLen)
		already_used[rd] = 1
		return rd
	}

	wrd := rand.Intn(allWeight)

	for i := 0; i < c.addrLen; i++ {
		if already_used[i] == 0 {
			if wrd < c.clientWeight[i] {
				already_used[i] = 1
				return i
			} else {
				wrd -= c.clientWeight[i]
			}
		}
	}

	log.Errorf("HttpClient:getAddr something err, not match any")
	rd := rand.Intn(c.addrLen)
	already_used[rd] = 1
	return rd

}

func (c *HttpClient) GetClientConfigInfo() *config.ClientInfo {
	return c.config
}

func (c *HttpClient) GetBodyFormat() string {
	if c.config.MessageFormat == "" {
		return HttpFormatJson
	} else {
		return c.config.MessageFormat
	}
}

func (c *HttpClient) doGet(url string) ([]byte, error) {

	var err error
	var resp *http.Response
	resp, err = c.httpclient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("send GET request error - %s", err)
	}

	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read GET response data error - %s", err)
	}

	if len(b) == 0 {
		return nil, ErrHttpResponseNull
	}

	//log.Debugf("httpclient Get:%s result:%s", url, string(b))

	return b, nil
}

func (c *HttpClient) doPost(url string, postdata []byte, contentType string) ([]byte, error) {

	var err error
	var resp *http.Response
	resp, err = c.httpclient.Post(url, contentType, bytes.NewReader(postdata))
	if err != nil {
		return nil, fmt.Errorf("send POST request error - %s", err)
	}

	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read POST response data error - %s", err)
	}

	if len(b) == 0 {
		return nil, ErrHttpResponseNull
	}

	//log.Debugf("httpclient Post:%s:%s result:%s", url,
	//	string(postdata), string(b))
	return b, nil
}

func (c *HttpClient) doPostWithHeader(url string, postdata []byte, contentType string, h http.Header) (http.Header, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(postdata))
	if err != nil {
		return nil, nil, fmt.Errorf("new req error - %s", err)
	}
	// 添加header
	if "" != contentType {
		req.Header.Set("Content-Type", contentType)
	}
	for k, vs := range h {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("send POST request error - %s", err)
	}

	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read POST response data error - %s", err)
	}

	if len(b) == 0 {
		return nil, nil, ErrHttpResponseNull
	}

	return resp.Header, b, nil
}

func (c *HttpClient) HttpGet(urlpath string) ([]byte, error) {

	countScope := countstat.GetCountGlobal().NewCountScope("Httpclient:" + c.config.Name + ":HttpGet")
	defer countScope.End()

	trynum := 1
	if c.config.TryNum > 0 {
		trynum = c.config.TryNum
	}

	if c.addrLen <= 0 {
		return nil, fmt.Errorf("no server addr")
	}

	startPos := rand.Intn(c.addrLen)

	var result []byte
	var err error
	var addr string
	var already_used_addr []int
	if c.config.WeightAuto != 0 {
		already_used_addr = make([]int, c.addrLen)
	}
	for t := 0; t < trynum && t < len(c.config.Addrs); t++ {

		if c.config.WeightAuto != 0 {
			startPos = c.getAddr(already_used_addr)
			atomic.AddUint32(&c.clientReqCount[startPos].totalReq, 1)
		} else {

			if c.addrLen <= 1 {
				startPos = (startPos + 1) % c.addrLen
			} else {
				startPos = (startPos + 1 + rand.Intn(c.addrLen-1)) % c.addrLen
			}
		}

		addr = c.config.Addrs[startPos]

		countScopeAddr := countstat.GetCountGlobal().NewCountScope("Httpclient:" + c.config.Name + ":HttpPost:" + addr)

		url := addr + urlpath
		result, err = c.doGet(url)
		if err == nil {
			countScopeAddr.End()
			return result, nil
		} else {
			countScopeAddr.SetErr()
			countScopeAddr.End()

			if c.config.WeightAuto != 0 {
				atomic.AddUint32(&c.clientReqCount[startPos].failReq, 1)
			}

			if strings.HasSuffix(err.Error(), "getsockopt: connection refused") {
				log.Warningf("connection refused retry: %s", err)
				trynum += 1
				continue
			}

			log.Warningf("Try HttpGet:%s fail:%s", url, err.Error())
		}
	}

	countScope.SetErr()
	log.Errorf("Client:%s HttpGet:%s fail - %s", c.config.Name, urlpath, err.Error())
	return nil, ErrHttpClient
}

func (c *HttpClient) HttpParamGet(path string, param map[string][]string) ([]byte, error) {
	urlValues := url.Values(param)
	url := fmt.Sprintf("%s?%s", path, urlValues.Encode())
	return c.HttpGet(url)
}

func (c *HttpClient) HttpPost(urlpath string, postdata []byte) ([]byte, error) {
	return c.httpPost(urlpath, postdata, ContentTypeJson)
}

func (c *HttpClient) HttpPostProtoBuf(urlpath string, postdata []byte) ([]byte, error) {
	return c.httpPost(urlpath, postdata, ContentTypeProtoBuf)
}

func (c *HttpClient) HttpPostRawData(urlpath string, postdata []byte) ([]byte, error) {
	return c.httpPost(urlpath, postdata, ContentTypeRawData)
}

func (c *HttpClient) HttpPostForm(urlpath string, data url.Values) ([]byte, error) {
	return c.httpPost(urlpath, []byte(data.Encode()), ContentTypeForm)
}

func (c *HttpClient) httpPost(urlpath string, postdata []byte, contentType string) ([]byte, error) {
	_, b, e := c.httpPostBase(urlpath, "", postdata, contentType, nil)
	return b, e
}

func (c *HttpClient) httpPostHeader(urlpath string, postdata []byte, contentType string, h http.Header) (http.Header, []byte, error) {
	return c.httpPostBase(urlpath, "", postdata, contentType, h)
}

func (c *HttpClient) httpPostWithAddr(urlpath string, inAddr string, postdata []byte, contentType string) ([]byte, error) {
	_, b, e := c.httpPostBase(urlpath, inAddr, postdata, contentType, nil)
	return b, e
}

func (c *HttpClient) httpPostBase(urlpath string, inAddr string, postdata []byte, contentType string, h http.Header) (http.Header, []byte, error) {

	countScope := countstat.GetCountGlobal().NewCountScope("Httpclient:" + c.config.Name + ":HttpPost")
	defer countScope.End()
	trynum := 1
	if c.config.TryNum > 0 {
		trynum = c.config.TryNum
	}

	var retHeader http.Header
	var result []byte
	var err error
	var addr string
	var already_used_addr []int
	if c.config.WeightAuto != 0 && inAddr == "" {
		already_used_addr = make([]int, c.addrLen)
	}

	confAddrs := c.config.Addrs
	if inAddr != "" {
		confAddrs = []string{inAddr}
	}
	addrsNum := len(confAddrs)
	if addrsNum <= 0 {
		return nil, nil, fmt.Errorf("no server addr")
	}

	startPos := rand.Intn(addrsNum)

	for t := 0; t < trynum && t < len(confAddrs); t++ {

		if c.config.WeightAuto != 0 && inAddr == "" {
			startPos = c.getAddr(already_used_addr)
			atomic.AddUint32(&c.clientReqCount[startPos].totalReq, 1)
		} else {
			if addrsNum <= 1 {
				startPos = (startPos + 1) % addrsNum
			} else {
				startPos = (startPos + 1 + rand.Intn(c.addrLen-1)) % addrsNum
			}
		}

		addr = confAddrs[startPos]

		countScopeAddr := countstat.GetCountGlobal().NewCountScope("Httpclient:" + c.config.Name + ":HttpPost:" + addr)

		url := addr + urlpath

		if nil == h {
			result, err = c.doPost(url, postdata, contentType)
		} else {
			retHeader, result, err = c.doPostWithHeader(url, postdata, contentType, h)
		}
		if err == nil || err == ErrHttpResponseNull {
			countScopeAddr.End()
			return retHeader, result, nil
		} else {

			countScopeAddr.SetErr()
			countScopeAddr.End()

			if c.config.WeightAuto != 0 && inAddr == "" {
				atomic.AddUint32(&c.clientReqCount[startPos].failReq, 1)
			}

			if strings.HasSuffix(err.Error(), "getsockopt: connection refused") {
				log.Warningf("connection refused retry: %s", err)
				trynum += 1
				continue
			}
			log.Warningf("Try HttpPost:%s fail:%s", url, err.Error())
		}
	}

	countScope.SetErr()
	log.Errorf("Client:%s HttpPost:%s fail - %s", c.config.Name, urlpath, err.Error())
	return retHeader, nil, ErrHttpClient
}

//根据seed参数作为参考选择请求哪个服务
func (c *HttpClient) HttpPostWithSeed(urlpath string, postdata []byte, contentType string, seed int) ([]byte, error) {

	countScope := countstat.GetCountGlobal().NewCountScope("Httpclient:" + c.config.Name + ":HttpPost")
	defer countScope.End()
	trynum := 1
	if c.config.TryNum > 0 {
		trynum = c.config.TryNum
	}
	addrnum := len(c.config.Addrs)

	if addrnum <= 0 {
		return nil, fmt.Errorf("no server addr")
	}

	startPos := seed
	if seed < 0 {
		startPos = rand.Intn(addrnum)
	}

	var result []byte
	var err error
	for t := 0; t < trynum; t++ {
		addr := c.config.Addrs[startPos%addrnum]
		startPos++
		url := addr + urlpath
		result, err = c.doPost(url, postdata, contentType)
		if err == nil {
			return result, nil
		} else {
			log.Warningf("Try HttpPost:%s dataLen:%d fail:%s ", url, len(postdata), err.Error())
		}
	}

	countScope.SetErr()
	log.Errorf("Client:%s HttpPost:%s fail - %s", c.config.Name, urlpath, err.Error())
	return nil, ErrHttpClient
}

//相当于导出httpPost，支持各种自定义的contentType
func (c *HttpClient) HttpPostCustom(urlpath string, postdata []byte, contentType string) ([]byte, error) {
	return c.httpPost(urlpath, postdata, contentType)
}

func (c *HttpClient) HttpPostWithHeader(urlpath string, h http.Header, postdata []byte) (http.Header, []byte, error) {
	return c.httpPostHeader(urlpath, postdata, ContentTypeJson, h)
}

func (c *HttpClient) HttpPostPBWithHeader(urlpath string, h http.Header, postdata []byte) (http.Header, []byte, error) {
	return c.httpPostHeader(urlpath, postdata, ContentTypeProtoBuf, h)
}

func (c *HttpClient) HttpPostPBWithAddr(urlpath string, inAddr string, postdata []byte) ([]byte, error) {
	return c.httpPostWithAddr(urlpath, inAddr, postdata, ContentTypeProtoBuf)
}
