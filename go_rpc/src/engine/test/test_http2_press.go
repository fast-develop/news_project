package main

import (
    "os"
    "bufio"
    "io"
    "strings"
    "time"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/http2"
	"io/ioutil"

	"encoding/json"
	"log"
	"net"
	"net/http"

	"w2vv2/proto/recallServer"
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
    f, err := os.Open("user_list")
    if err != nil {
        panic(err)
    }   
    defer f.Close()

    rd := bufio.NewReader(f)
    for {
        line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
        
        if err != nil || io.EOF == err {
            break
        }   
        fmt.Println(line)
        DoPost(strings.Replace(line, "\n", "", -1))
        time.Sleep(10*time.Millisecond)
    }    
}

func DoPost(userid string) {
	client := http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	reqPb := recallServer.RecallRequest{}
	param := make(map[string]string)
	param["top_n"] = "10"
	param["matrixCF_userRecall_client_name"] = "matrixCF_userRecall_client"
	param["Uid"] = userid

	data, err := json.Marshal(param)
	if err != nil {
		fmt.Println("json marshal failed:", err)
		return
	}
	fmt.Printf(string(data))

	reqPb.Data = data

	byteReq, err := proto.Marshal(&reqPb)
	if err != nil {
		log.Fatal("proto marshal error")
		return
	}

	hdr, payload := msgHeader(byteReq, nil)
	req_data := append(hdr, payload...)

	req, err := http.NewRequest("POST", "http://localhost:30000/recallServer.Recall/Handler", bytes.NewReader(req_data))
	if err != nil {
		log.Fatal("new request error")
	}

	req.Header.Set("content-type", "application/grpc+proto")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/serverpb.Recall/Handler")
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

	var rspPb recallServer.RecallResponse
	if err := proto.Unmarshal(resp_data, &rspPb); err != nil {
		fmt.Println("proto unmarshal error")
		return
	}

	fmt.Printf("%v", &rspPb)

}

