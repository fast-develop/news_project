package shmcache

// #cgo CFLAGS: -I"./include"
// #cgo LDFLAGS: ${SRCDIR}/lib/libshmcache.a ${SRCDIR}/lib/libfastcommon.a -ldl -lm
/*
#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>
#include <unistd.h>
#include <dirent.h>
#include <sys/types.h>
#include <sys/stat.h>
#include "fastcommon/logger.h"
#include "fastcommon/hash.h"
#include "shmcache/shm_hashtable.h"
#include "shmcache/shmcache.h"

typedef struct ShmCacheParameter {
	char        *file_name;
	uint64_t    max_memory;
	uint64_t    min_memory;
	uint64_t    segment_size;
	uint64_t    max_key_count;
	uint64_t    max_value_size;
	char        *hash_function;
	uint64_t    recycle_key_once;
	bool        recycle_valid_entries;
	uint64_t    avg_key_ttl;
	uint64_t    discard_memory_size;
	uint64_t    max_fail_times;
	uint64_t    sleep_us_when_recycle_valid_entries;
	uint64_t    trylock_interval_us;
	uint64_t    detect_deadlock_interval_ms;
	uint64_t    expire;

} ShmCacheParam;

int shm_log_init(char *path, char *prefix)
{
	if (log_init() != 0) {
		return -1;
	}

	log_set_prefix(path, prefix);
	return 0;
}

int shm_make_dir(char *path)
{
    char dir_name[1024] = {0};
    int i = 0,e=0;
    strncpy(dir_name, path, strlen(path));
    int len = strlen(path);
    dir_name[len] = '\0';
    if (dir_name[len-1] != '/') {
        dir_name[len] = '/';
        dir_name[len+1] = '\0';
        len = len + 1;
    }


    for (i = 1; i < len; i++) {
        if (dir_name[i] == '/') {
            dir_name[i] = '\0';
		    e = access(dir_name, F_OK);
			//printf("access %s %d\n", dir_name, e);
            if (e != 0) {
                if (mkdir(dir_name, 0755) == -1) {
                    return -1;
                }
            }
            dir_name[i] = '/';
        }
    }

    return 0;
}

int del_path_file(char *folder)
{
    struct dirent * file_name;    // return value for readdir()
    DIR * dir;                   // return value for opendir()
    dir = opendir(folder);
    if (NULL == dir) {
        return -1;
    }

    //LOG(INFO) << "Successfully opened the dir:" << folder;

    while ((file_name = readdir(dir)) != NULL) {
        // get rid of "." and ".."
        if(strcmp(file_name->d_name , ".") == 0 ||
            strcmp(file_name->d_name , "..") == 0) {
            continue;
        }
        char tmp_path[1024] = {0};
        sprintf(tmp_path, "%s/%s", folder, file_name->d_name);
        if(unlink(tmp_path) == -1)
        {
            return -1;
        }

    }
    return 0;
}

int shm_cache_init(void **context, ShmCacheParam *param)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)malloc(sizeof(struct shmcache_context));
	if (!shm_context) {
		printf("shm_cache_init error\n");
		return -1;
	}

	char path[1024] = {0};
	if (strlen(param->file_name) == 0) {
		printf("file_name len is 0\n");
        return -1;
    } else {
        char *last = strrchr(param->file_name, '/');
        if (last) {
            int path_len = last - param->file_name;
            memcpy(path, param->file_name, path_len);
            if (shm_make_dir(path) != 0) {
				printf("shm_make_dir error\n");
                return -1;
            }
        }
    }

    struct shmcache_config cache_config;
    cache_config.type = SHMCACHE_TYPE_MMAP;
   	memcpy(cache_config.filename, param->file_name, strlen(param->file_name));
   	cache_config.max_memory = param->max_memory;
    cache_config.min_memory = param->min_memory;
    cache_config.segment_size = param->segment_size;
    if (cache_config.max_memory / cache_config.segment_size > 255) {
        cache_config.segment_size = cache_config.max_memory / 255;
    }
    cache_config.max_key_count = param->max_key_count;
    if (cache_config.max_key_count <= 0) {
		printf("max_key_count <= 0\n");
        return -1;
    }
    cache_config.max_value_size = param->max_value_size;
    if (cache_config.max_value_size <= 0) {
        printf("max_value_size <= 0\n");
        return -1;
    }

    if (strlen(param->hash_function) == 0) {
        cache_config.hash_func = simple_hash;
    } else {
        void *handle;
        handle = dlopen(NULL, RTLD_LAZY);
        if (handle == NULL) {
		    printf("dlopen error\n");
            return -1;
        }
        cache_config.hash_func = (HashFunc)dlsym(handle, param->hash_function);
        if (cache_config.hash_func == NULL) {
		    printf("hash_func is nil\n");
            return -1;
        }
        dlclose(handle);
    }
    cache_config.va_policy.avg_key_ttl = param->avg_key_ttl;
    cache_config.va_policy.discard_memory_size = param->discard_memory_size;
    cache_config.va_policy.max_fail_times = param->max_fail_times;
    cache_config.lock_policy.trylock_interval_us = param->trylock_interval_us;
    if (cache_config.lock_policy.trylock_interval_us <= 0) {
		printf("trylock_interval_us <=0 \n");
        return -1;
    }
    cache_config.va_policy.sleep_us_when_recycle_valid_entries = param->sleep_us_when_recycle_valid_entries;
    cache_config.lock_policy.detect_deadlock_interval_ms = param->detect_deadlock_interval_ms;
    if (cache_config.lock_policy.detect_deadlock_interval_ms <= 0) {
		printf("detect_deadlock_interval_ms<=0\n");
        return -1;
    }
    cache_config.recycle_key_once = param->recycle_key_once;
    if (cache_config.recycle_key_once <= 0) {
        cache_config.recycle_key_once = -1;
    }
    cache_config.recycle_valid_entries = param->recycle_valid_entries;

	if (shmcache_init(shm_context, &cache_config, true, true) != 0) {
		if (del_path_file(path)) {
			if (shmcache_init(shm_context, &cache_config, true, true) != 0) {
		        printf("shmcache_init error again\n");
				return -1;
			}
		} else {
		    printf("del_path_file error\n");
			return -1;
		}
	}

	*context = shm_context;

	return 0;
}

void shm_cache_release(void *context)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)context;

	if (shm_context) {
		shmcache_destroy(context);
		free (context);
	}
}

int shm_cache_set(void *context, char *key, uint64_t key_size, char *val, uint64_t val_size, uint64_t time)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)context;

	struct shmcache_key_info cache_key;
	struct shmcache_value_info cache_value;

	cache_key.length = key_size;
	cache_key.data = strdup(key);

	cache_value.length = val_size;
	cache_value.data = val;
	cache_value.options = SHMCACHE_SERIALIZER_STRING;
	cache_value.expires = time;

	int err = shmcache_set_ex(shm_context, &cache_key, &cache_value);
    free(cache_key.data);
	if (err != 0) {
		return -1;
	}

	return 0;
}

int shm_cache_get(void *context, char *key, uint64_t key_size, char **val, int *val_size, uint64_t time)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)context;

	struct shmcache_key_info cache_key;
    struct shmcache_value_info cache_value;

    cache_key.length = key_size;
    cache_key.data = strdup(key);

    int err = shmcache_get(shm_context, &cache_key, &cache_value);
    free(cache_key.data);
    if (err != 0) {
        return -1;
    }

    *val = (char *)malloc(cache_value.length);
    memcpy(*val, cache_value.data, cache_value.length);
    *val_size = cache_value.length;

    return 0;
}

uint64_t shm_cache_get_size(void *context)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)context;
	struct shmcache_stats stat;
	shmcache_stats(shm_context, &stat);
	return stat.memory.used;
}

uint64_t shm_cache_get_count(void *context)
{
	struct shmcache_context *shm_context = (struct shmcache_context*)context;
    struct shmcache_stats stat;
    shmcache_stats(shm_context, &stat);
    return stat.hashtable.count;
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

type ShmCache struct {
	context unsafe.Pointer
	config  ShmCacheConfig
}

type ShmCacheConfig struct {
	file_name                           string
	max_memory                          uint64
	min_memory                          uint64
	segment_size                        uint64
	max_key_count                       uint64
	max_value_size                      uint64
	hash_function                       string
	recycle_key_once                    uint64
	recycle_valid_entries               bool
	avg_key_ttl                         uint64
	discard_memory_size                 uint64
	max_fail_times                      uint64
	sleep_us_when_recycle_valid_entries uint64
	trylock_interval_us                 uint64
	detect_deadlock_interval_ms         uint64
	expire                              uint64
}

func (cache *ShmCache) Init(config *ShmCacheConfig) error {
	shmCachePara := &C.ShmCacheParam{}
	shmCachePara.file_name = C.CString(config.file_name)
	shmCachePara.max_memory = C.uint64_t(config.max_memory)
	shmCachePara.min_memory = C.uint64_t(config.min_memory)
	shmCachePara.segment_size = C.uint64_t(config.segment_size)
	shmCachePara.max_key_count = C.uint64_t(config.max_key_count)
	shmCachePara.max_value_size = C.uint64_t(config.max_value_size)
	shmCachePara.hash_function = C.CString(config.hash_function)
	shmCachePara.recycle_key_once = C.uint64_t(config.recycle_key_once)
	shmCachePara.recycle_valid_entries = C.bool(config.recycle_valid_entries)
	shmCachePara.avg_key_ttl = C.uint64_t(config.avg_key_ttl)
	shmCachePara.discard_memory_size = C.uint64_t(config.discard_memory_size)
	shmCachePara.max_fail_times = C.uint64_t(config.max_fail_times)
	shmCachePara.sleep_us_when_recycle_valid_entries = C.uint64_t(config.sleep_us_when_recycle_valid_entries)
	shmCachePara.trylock_interval_us = C.uint64_t(config.trylock_interval_us)
	shmCachePara.detect_deadlock_interval_ms = C.uint64_t(config.detect_deadlock_interval_ms)
	shmCachePara.expire = C.uint64_t(config.expire)

	ret := C.shm_cache_init(&cache.context, shmCachePara)
	if ret != 0 {
		var err error = errors.New("cache init error")
		log.Infof(err.Error())
		return err
	}

	return nil
}

func (cache *ShmCache) Release() error {
	C.shm_cache_release(cache.context)
	return nil
}

func (cache *ShmCache) Set(key *string, value []byte, cur_time time.Time, expire uint32) error {
	c_key := C.CString(*key)
	c_key_len := C.uint64_t(len(*key))
	c_val := (*C.char)(unsafe.Pointer(&value[0]))
	c_val_len := C.uint64_t(len(value))
	c_expire := C.uint64_t(uint64(cur_time.Unix()) + uint64(expire))

	//fmt.Println("set", *key, c_val_len)
	if C.shm_cache_set(cache.context, c_key, c_key_len, c_val, c_val_len, c_expire) != 0 {
		var err error = errors.New("cache set error")
		log.Errorf(err.Error())
		return err

	}

	return nil

}

func (cache *ShmCache) Get(key *string, cur_time time.Time) ([]byte, bool) {
	c_key := C.CString(*key)
	c_key_len := C.uint64_t(len(*key))
	var get_val *C.char
	var get_val_size C.int
	c_time := C.uint64_t(uint64(cur_time.Unix()))

	if C.shm_cache_get(cache.context, c_key, c_key_len, &get_val, &get_val_size, c_time) != 0 {
		return nil, false
	}
	//fmt.Println("get", key, get_val_size)

	data := C.GoBytes(unsafe.Pointer(get_val), get_val_size)

	C.free(unsafe.Pointer(get_val))
	return data, true
}

func NewShmCache(cacheInfo *config.ShmCacheInfo) (*ShmCache, error) {
	var shmCache = &ShmCache{}
	var shmCacheConfig = &ShmCacheConfig{}

	shmCacheConfig.file_name = cacheInfo.FileName
	shmCacheConfig.max_memory = cacheInfo.MaxMemory
	shmCacheConfig.min_memory = cacheInfo.MinMemory
	shmCacheConfig.segment_size = cacheInfo.SegmentSize
	shmCacheConfig.max_key_count = cacheInfo.MaxKeyCount
	shmCacheConfig.max_value_size = cacheInfo.MaxValueSize
	shmCacheConfig.hash_function = cacheInfo.HashFunction
	shmCacheConfig.recycle_key_once = cacheInfo.RecycleKeyOnce
	shmCacheConfig.recycle_valid_entries = cacheInfo.RecycleValidEntries
	shmCacheConfig.avg_key_ttl = cacheInfo.AvgKeyTTL
	shmCacheConfig.discard_memory_size = cacheInfo.DiscardMemorySize
	shmCacheConfig.max_fail_times = cacheInfo.MaxFailTimes
	shmCacheConfig.sleep_us_when_recycle_valid_entries = cacheInfo.SleepUsWhenRecycleValidEntries
	shmCacheConfig.trylock_interval_us = cacheInfo.TrylockIntervalUs
	shmCacheConfig.detect_deadlock_interval_ms = cacheInfo.DetectDeadlockIntervalMs
	shmCacheConfig.expire = cacheInfo.Expire

	err := shmCache.Init(shmCacheConfig)
	if err != nil {
		log.Errorf("cache init error")
		fmt.Println(err)
		return nil, err
	}

	return shmCache, nil
}

func ShmLogInit(path *string, prefix *string) error {
	c_path := C.CString(*path)
	c_prefix := C.CString(*prefix)
    if C.shm_make_dir(c_path) != 0 {
        return errors.New("shm_make_dir logpath error")
    }
	if C.shm_log_init(c_path, c_prefix) != 0 {
		return errors.New("shm log init error")
	}
	return nil
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
