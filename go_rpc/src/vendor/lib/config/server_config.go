package config

import "encoding/xml"

type ServerConfig struct {
	XMLName xml.Name `xml:"serverconfig"`
	Port    string   `xml:"serverport"`
    Name    string   `xml:"servername"`
}
