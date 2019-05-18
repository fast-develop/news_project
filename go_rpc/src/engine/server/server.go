package server

import (
    "context"
    "errors"
    "fmt"
    "net"
    log "github.com/golang/glog"

    "lib/shmcache"
    "lib/client"

    "engine/config"
    "engine/handler"
    "engine/proto/engine"
)

type ServerImpl struct {
    Listen net.Listener
    Config *config.Config
}

func (s *ServerImpl) Handler(ctx context.Context, in *engine.EngineRequest) (*engine.EngineResponse, error) {
    rsp, err := handler.HandleRequest(ctx, in)
    return rsp, err
}

func NewServerImpl(conf *config.Config) (*ServerImpl, error) {
	if conf.ShmCacheConf == nil {
		log.Errorf("cache config is nil")
		err := errors.New("cache config is nil")
		return nil, err
	}
	err := shmcache.ShmCacheMgr.Init(conf.ShmCacheConf)
	if err != nil {
		log.Errorf("cache init error")
		return nil, err
	}
    fmt.Println("cache config sucess")

	if conf.ClientConf == nil {
		log.Errorf("client config is nil")
		err := errors.New("client config is nil")
		return nil, err
	}
	err = client.ClientMgr.Init(conf.ClientConf)
	if err != nil {
		log.Errorf("client init error")
		return nil, err
	}
    fmt.Println("client config sucess")

	listen, err := net.Listen("tcp", conf.ServerConf.Port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return nil, err
	}

	srv := &ServerImpl{
		listen,
		conf,
	}

	return srv, nil
}
