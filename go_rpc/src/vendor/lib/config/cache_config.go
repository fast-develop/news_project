package config

import "encoding/xml"

type CacheInfo struct {
	Name   string `xml:"name,attr"`
	Expire uint32 `xml:"expire,attr"`
	Size   uint64 `xml:"size,attr"`
}

type CacheConfig struct {
	XMLName    xml.Name     `xml:"cacheconfig"`
	CacheInfos []*CacheInfo `xml:"cacheinfo"`
}
