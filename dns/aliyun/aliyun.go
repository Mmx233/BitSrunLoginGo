package aliyun

import (
	"github.com/Mmx233/BitSrunLoginGo/dns/models"
	"github.com/Mmx233/BitSrunLoginGo/dns/util"
)

type DnsProvider struct {
}

func New(conf map[string]interface{}) (models.DnsProvider, error) {
	var p DnsProvider
	return &p, util.DecodeConfig(conf, &p)
}
