package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"io/ioutil"

	"log"
	"net"
	"net/http"

	"new_user_predict/proto/newUserPredict"
)

type payloadFormat uint8

const (
	compressionNone payloadFormat = 0 // no compression
	compressionMade payloadFormat = 1 // compressed
)

const (
	payloadLen = 1
	sizeLen    = 4
	headerLen  = payloadLen + sizeLen
)

func msgHeader(data, compData []byte) (hdr []byte, payload []byte) {
	hdr = make([]byte, headerLen)
	if compData != nil {
		hdr[0] = byte(compressionMade)
		data = compData
	} else {
		hdr[0] = byte(compressionNone)
	}

	// Write length of payload into buf
	binary.BigEndian.PutUint32(hdr[payloadLen:], uint32(len(data)))
	return hdr, data
}

func main() {
	client := http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	reqPb := newUserPredict.PredictRequest{}
    userId := "127731698"
    reqPb.UserId = &userId
    docIdList := []int32{154985670,156398773,160201829,160285254,169245930,169347134,169521532,169638435,169732312,169975680}
    fromList := []int32{1,2,3,4,5,6,7,8,9,10}
    reqPb.DocIdList = docIdList
    reqPb.FromList = fromList

	byteReq, err := proto.Marshal(&reqPb)
	if err != nil {
		log.Fatal("proto marshal error")
		return
	}

	hdr, payload := msgHeader(byteReq, nil)
	req_data := append(hdr, payload...)

	req, err := http.NewRequest("POST", "http://localhost:30000/newUserPredict.Predict/Handler", bytes.NewReader(req_data))
	if err != nil {
		log.Fatal("new request error")
	}

	req.Header.Set("content-type", "application/grpc+proto")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/newUserPredict.Predict/Handler")
	req.Header.Set("scheme", "http")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("error request: %v", err))
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Proto)

	defer resp.Body.Close()

	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error")
		return
	}

	if len(b) == 0 {
		fmt.Println("body is nil")
		return
	}

	resp_data := b[5:]

	var rspPb newUserPredict.PredictResponse
	if err := proto.Unmarshal(resp_data, &rspPb); err != nil {
		fmt.Println("proto unmarshal error")
		return
	}

	fmt.Printf("%v", &rspPb)

}

