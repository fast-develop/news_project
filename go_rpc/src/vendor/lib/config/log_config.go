package config

import "encoding/xml"

type LogConfig struct {
	XMLName xml.Name `xml:"logconfig"`
	Path    string   `xml:"path"`
	Level   string   `xml:"level"`
}
