package cache

import (
	log "google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/grpclog/glogger"

	"lib/config"
)

type CacheManager struct {
	cacheMap map[string]*MCache
}

func (cacheMgr *CacheManager) GetCache(name string) (*MCache, bool) {
	cache, ok := cacheMgr.cacheMap[name]
	return cache, ok
}

func (cacheMgr *CacheManager) Init(cacheconf *config.CacheConfig) error {
	cacheMgr.cacheMap = make(map[string]*MCache)

	for _, cacheinfo := range cacheconf.CacheInfos {
		cache, err := NewCache(cacheinfo)
		if err != nil {
			log.Errorf("cache init error")
			return err
		}

		cacheMgr.cacheMap[cacheinfo.Name] = cache
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

var CacheMgr = &CacheManager{}
