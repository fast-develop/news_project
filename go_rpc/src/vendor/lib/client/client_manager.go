package client

import (

    _ "google.golang.org/grpc/grpclog/glogger"
    "lib/config"
    "fmt"
)

type KvClient interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
    SetWithttl(key string, value interface{}, ttl int) error
}

type ClientManager struct {
    httpClientMap  map[string]*HttpClient
    kvClientMap    map[string]KvClient
    redisClientMap map[string]*RedisClient
}

func (cm *ClientManager) GetHttpClient(name string) (*HttpClient, bool) {
    c, ok := cm.httpClientMap[name]
    return c, ok
}

func (cm *ClientManager) GetKvClient(name string) (KvClient, bool) {
    c, ok := cm.kvClientMap[name]
    return c, ok
}

func (cm *ClientManager) GetRedisClient(name string) (*RedisClient, bool) {
    c, ok := cm.redisClientMap[name]
    return c, ok
}

func (cm *ClientManager) Init(clientconf *config.ClientConfig) error {
    cm.httpClientMap = make(map[string]*HttpClient)
    cm.kvClientMap = make(map[string]KvClient)
    cm.redisClientMap = make(map[string]*RedisClient)

    for _, cinfo := range clientconf.HttpClients {
        client := NewHttpClient(cinfo)
        cm.httpClientMap[cinfo.Name] = client
    }

    for _, cinfo := range clientconf.RedisClients {
        client := NewRedisClient(cinfo)
        cm.kvClientMap[cinfo.Name] = client
        cm.redisClientMap[cinfo.Name] = client
    }

    return nil
}

var ClientMgr = &ClientManager{}
