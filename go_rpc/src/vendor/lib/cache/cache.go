package cache

// #cgo CFLAGS: -I"./include"
// #cgo LDFLAGS: -L"./lib" -lmcache
/*
#include <stdlib.h>
#include <stdio.h>
#include "./include/mcache.h"

void print(const char* msg) {
	printf("%s\n", msg);
}

int cache_init(void **p_cache, uint64_t size, uint32_t expire)
{
	char err_buf[1024];
	*p_cache = mcache_kv_init(size, expire, err_buf, sizeof(err_buf));
	if (!*p_cache) {
		return -1;
	}
	return 0;
}

void cache_release(void *p_cache)
{
	if (p_cache) {
		mcache_kv_free(p_cache);
	}
}

int cache_set(void *p_cache, char *key, char *val, uint64_t val_size, uint64_t time)
{
	int ret = mcache_kv_set(p_cache, (u_char*)key, (u_char*)val, val_size, time);
	if (ret != MC_SUCCESS) {
		return -1;
	}

	return 0;
}

int cache_get(void *p_cache, char *key, char **val, int *val_size, uint64_t time)
{
	if (mcache_kv_get(p_cache, (u_char *)key, (u_char **)val, (uint32_t*)val_size, time) != MC_SUCCESS) {
        return -1;
	}

	return 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
	log "google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/grpclog/glogger"
	_ "reflect"
	"time"
	"unsafe"

	"lib/config"
)

type MCache struct {
	p_cache unsafe.Pointer
	name    string
	size    uint64
	expire  uint32
}

func (cache *MCache) Init(size uint64, expire uint32) error {
	c_size := C.uint64_t(size)
	c_expire := C.uint32_t(expire)
	ret := C.cache_init(&cache.p_cache, c_size, c_expire)
	if ret != 0 {
		var err error = errors.New("cache init error")
		log.Infof(err.Error())
		return err
	}

	return nil
}

func (cache *MCache) Release() error {
	C.cache_release(cache.p_cache)
	return nil
}

/*
func (cache *MCache) Set(key *string, value *string, cur_time time.Time, expire uint32) error {
	c_key := C.CString(*key)
	c_val := C.CString(*value)
	c_val_len := C.uint64_t(len(*value))
	c_expire := C.uint64_t(uint64(cur_time.Unix()) + uint64(expire))

    fmt.Println("set", key, len(*value))
	if C.cache_set(cache.p_cache, c_key, c_val, c_val_len, c_expire) != 0 {
		var err error = errors.New("cache set error")
		log.Errorf(err.Error())
		return err

	}

	return nil

}

func (cache *MCache) Get(key *string, value *string, cur_time time.Time) bool {
	c_key := C.CString(*key)
	var get_val *C.char
	var get_val_size C.uint32_t
	c_time := C.uint64_t(uint64(cur_time.Unix()))

	if C.cache_get(cache.p_cache, c_key, &get_val, &get_val_size, c_time) != 0 {
		return false
	}
    fmt.Println("get", key, get_val_size)

	*value = C.GoString(get_val)

	C.free(unsafe.Pointer(get_val))
	return true
}
*/

func (cache *MCache) Set(key *string, value []byte, cur_time time.Time, expire uint32) error {
	c_key := C.CString(*key)
	c_val := (*C.char)(unsafe.Pointer(&value[0]))
	c_val_len := C.uint64_t(len(value))
	c_expire := C.uint64_t(uint64(cur_time.Unix()) + uint64(expire))

    fmt.Println("set", *key, c_val_len)
	if C.cache_set(cache.p_cache, c_key, c_val, c_val_len, c_expire) != 0 {
		var err error = errors.New("cache set error")
		log.Errorf(err.Error())
		return err

	}

	return nil

}

func (cache *MCache) Get(key *string, cur_time time.Time) ([]byte, bool) {
	c_key := C.CString(*key)
	var get_val *C.char
	var get_val_size C.int
	c_time := C.uint64_t(uint64(cur_time.Unix()))

	if C.cache_get(cache.p_cache, c_key, &get_val, &get_val_size, c_time) != 0 {
		return nil, false
	}
    fmt.Println("get", key, get_val_size)

    data := C.GoBytes(unsafe.Pointer(get_val), get_val_size)

	C.free(unsafe.Pointer(get_val))
	return data, true
}


func NewCache(cacheinfo *config.CacheInfo) (*MCache, error) {
	var cache = &MCache{}
	size := cacheinfo.Size
	expire := cacheinfo.Expire

	err := cache.Init(size, expire)
	if err != nil {
		log.Errorf("cache init error")
		return nil, err
	}

	return cache, nil
}

/*
func main() {
	cache := new(MCache)
	err := cache.Init(5000000, 5)
	if err != nil {
		fmt.Println("cache_init error")
		return
	}

	key := "hello"
	val := "11111111111111111111112e3123213213"
	cur_time := time.Now()
	if err = cache.Set(&key, &val, cur_time, 5); err != nil {
		fmt.Println("cache_set error")
	}

	var get_val string
	if err = cache.Get(&key, &get_val, cur_time); err != nil {
		fmt.Println("cache_get error")
	}

	fmt.Println(get_val)

	cache.Release()
}
*/
