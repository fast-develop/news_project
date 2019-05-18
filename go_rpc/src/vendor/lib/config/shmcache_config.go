package config

import "encoding/xml"

type ShmCacheInfo struct {
	Name                           string `xml:"name,attr"`
	FileName                       string `xml:"file_name"`
	MaxMemory                      uint64 `xml:"max_memory"`
	MinMemory                      uint64 `xml:"min_memory"`
	SegmentSize                    uint64 `xml:"segment_size"`
	MaxKeyCount                    uint64 `xml:"max_key_count"`
	MaxValueSize                   uint64 `xml:"max_value_size"`
	HashFunction                   string `xml:"hash_function"`
	RecycleKeyOnce                 uint64 `xml:"recycle_key_once"`
	RecycleValidEntries            bool   `xml:"recycle_valid_entries"`
	AvgKeyTTL                      uint64 `xml:"avg_key_ttl"`
	DiscardMemorySize              uint64 `xml:"discard_memory_size"`
	MaxFailTimes                   uint64 `xml:"max_fail_times"`
	SleepUsWhenRecycleValidEntries uint64 `xml:"sleep_us_when_recycle_valid_entries"`
	TrylockIntervalUs              uint64 `xml:"trylock_interval_us"`
	DetectDeadlockIntervalMs       uint64 `xml:"detect_deadlock_interval_ms"`
	Expire                         uint64 `xml:"expire"`
}

type ShmCacheConfig struct {
	XMLName       xml.Name        `xml:"shmcacheconfig"`
	Path          string          `xml:"logpath,attr"`
	Prefix        string          `xml:"logprefix,attr"`
	ShmCacheInfos []*ShmCacheInfo `xml:"shmcacheinfo"`
}
