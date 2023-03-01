package dnsUtil

import (
	"github.com/mitchellh/mapstructure"
)

func DecodeConfig(conf map[string]interface{}, output interface{}) error {
	return mapstructure.Decode(conf, output)
}
