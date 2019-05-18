package config

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "os"
    log "github.com/golang/glog"
    "lib/config"
)

type Config struct {
    XMLName         xml.Name                    `xml:"config"`
    ServerConf      *config.ServerConfig        `xml:"serverconfig"`
    ClientConf      *config.ClientConfig        `xml:"clientconfig"`
    ShmCacheConf    *config.ShmCacheConfig      `xml:"shmcacheconfig"`
    PrometheConf    *config.PrometheConfig      `xml:"prometheusconfig"`
    LogConf         *config.LogConfig           `xml:"logconfig"`
}

func (conf *Config)Init(filePath *string) error {
	file, err := os.Open(*filePath) // For read access.
	if err != nil {
		log.Errorf(err.Error())
		fmt.Println(err.Error())
		return err
	}

	data, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		log.Errorf(err.Error())
		fmt.Println(err.Error())
		return err
	}

	err = xml.Unmarshal(data, conf)
	if err != nil {
		log.Errorf("xml unmarshal error")
		fmt.Println(err.Error())
		return err
	}

	return nil
}

var Conf = &Config{}

