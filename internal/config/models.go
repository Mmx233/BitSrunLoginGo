package config

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/aliyun"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/cloudflare"
	"github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns/dnspod"
)

type (
	GuardianConf struct {
		Enable   bool `json:"enable" yaml:"enable"`
		Duration uint `json:"duration" yaml:"duration"`
	}

	BackoffConf struct {
		Enable          bool `json:"enable" yaml:"enable"`
		MaxRetries      uint `json:"max_retries" yaml:"max_retries"`
		InitialDuration uint `json:"initial_duration" yaml:"initial_duration"`
		MaxDuration     uint `json:"max_duration" yaml:"max_duration"`

		ExponentFactor   uint `json:"exponent_factor" yaml:"exponent_factor"`
		InterConstFactor uint `json:"inter_const_factor" yaml:"inter_const_factor"`
		OuterConstFactor uint `json:"outer_const_factor" yaml:"outer_const_factor"`
	}

	BasicConf struct {
		Https          bool   `json:"https" yaml:"https"`
		SkipCertVerify bool   `json:"skip_cert_verify" yaml:"skip_cert_verify"`
		Timeout        uint   `json:"timeout" yaml:"timeout"`
		Interfaces     string `json:"interfaces" yaml:"interfaces"`
	}

	LogConf struct {
		DebugLevel bool   `json:"debug_level" yaml:"debug_level"`
		WriteFile  bool   `json:"write_file" yaml:"write_file"`
		FilePath   string `json:"log_path" yaml:"log_path"`
		FileName   string `json:"log_name" yaml:"log_name"`
	}

	DdnsProviderConfigSum struct {
		dnspod.DnsPod         `yaml:",inline"`
		cloudflare.Cloudflare `yaml:",inline"`
		aliyun.Aliyun         `yaml:",inline"`
	}

	DdnsConf struct {
		Enable   bool                  `json:"enable" yaml:"enable"`
		TTL      uint                  `json:"ttl" yaml:"ttl"`
		Domain   string                `json:"domain" yaml:"domain"`
		Provider string                `json:"provider" yaml:"provider"`
		Config   DdnsProviderConfigSum `json:"config" yaml:"config"`
	}

	RealityConf struct {
		Enable bool   `json:"enable" yaml:"enable"`
		Addr   string `json:"addr" yaml:"addr"`
	}
)

type SettingsConf struct {
	Basic        BasicConf              `json:"basic" yaml:"basic"`
	Guardian     GuardianConf           `json:"guardian" yaml:"guardian"`
	Backoff      BackoffConf            `json:"backoff" yaml:"backoff"`
	Log          LogConf                `json:"log" yaml:"log"`
	DDNS         DdnsConf               `json:"ddns" yaml:"ddns"`
	Reality      RealityConf            `json:"reality" yaml:"reality"`
	CustomHeader map[string]interface{} `json:"custom_header" yaml:"custom_header"`
}
