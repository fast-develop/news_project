package shmcache

import (
	//"fmt"
	log "google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/grpclog/glogger"

	"lib/config"
)

type ShmCacheManager struct {
	shmCacheMap map[string]*ShmCache
}

func (shmCacheMgr *ShmCacheManager) GetCache(name string) (*ShmCache, bool) {
	shmCache, ok := shmCacheMgr.shmCacheMap[name]
	return shmCache, ok
}

func (shmCacheMgr *ShmCacheManager) Init(shmCacheConf *config.ShmCacheConfig) error {
	shmCacheMgr.shmCacheMap = make(map[string]*ShmCache)

	ShmLogInit(&shmCacheConf.Path, &shmCacheConf.Prefix)
	for _, shmCacheInfo := range shmCacheConf.ShmCacheInfos {
		shmCache, err := NewShmCache(shmCacheInfo)
		if err != nil {
			log.Errorf("cache init error")
			return err
		}

		shmCacheMgr.shmCacheMap[shmCacheInfo.Name] = shmCache
	}

	return nil
}

/*
func NewCacheManager(conf *config.CacheConfig) (*CacheManager, error) {
	var cache_mgr = &CacheManager{}
	err := cache_mgr.Init(conf)
	if err != nil {
		log.Errorf("cache manager init error")
		return nil, err
	}

	return cache_mgr, nil
}
*/

var ShmCacheMgr = &ShmCacheManager{}
