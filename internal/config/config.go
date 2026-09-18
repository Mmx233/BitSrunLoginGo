package config

import (
	"os"
	"time"

	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/internal/dns/aliyun"
	"github.com/Mmx233/BitSrunLoginGo/internal/dns/cloudflare"
	"github.com/Mmx233/BitSrunLoginGo/internal/dns/dnspod"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/tool"
)

type ConfFromFile struct {
	Form     srun.LoginForm `json:"form" yaml:"form"`
	Meta     srun.LoginMeta `json:"meta" yaml:"meta"`
	Settings SettingsConf   `json:"settings" yaml:"settings"`
}

var (
	Form     *srun.LoginForm
	Meta     *srun.LoginMeta
	Settings *SettingsConf

	Timeout time.Duration
)

func init() {
	initLogPre()

	logger := Logger.WithField(keys.LogComponent, "init")
	reader := newReaderFromPath(flags.Path)

	// 生成配置文件
	exist, err := tool.File.Exists(flags.Path)
	if err != nil {
		logger.Fatalln("读取配置文件失败：", err)
	} else if !exist {
		var data []byte
		data, err = reader.Marshal(&defaultConfig)
		if err != nil {
			logger.Fatalln("生成配置文件失败：", err)
		}
		if err = os.WriteFile(flags.Path, data, 0600); err != nil {
			logger.Fatalln("写入配置文件失败：", err)
		}
		logger.Infoln("已生成配置文件，请编辑 '" + flags.Path + "' 然后重试")
		os.Exit(0)
	}

	// 读取配置文件
	data, err := os.ReadFile(flags.Path)
	if err != nil {
		logger.Fatalln("读取配置失败：", err)
	}
	var fileConf ConfFromFile
	if err = reader.Unmarshal(data, &fileConf); err != nil {
		logger.Fatalln("解析配置失败：", err)
	}
	Form = &fileConf.Form
	Meta = &fileConf.Meta
	Settings = &fileConf.Settings
	Timeout = time.Duration(Settings.Basic.Timeout) * time.Second

	// flag 配置覆写
	if flags.Debug {
		Settings.Log.DebugLevel = true
	}
	if flags.Acid != "" {
		Meta.Acid = flags.Acid
	}

	initLogFinal()
	initBackoff()
}

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
		Https              bool   `json:"https" yaml:"https"`
		SkipCertVerify     bool   `json:"skip_cert_verify" yaml:"skip_cert_verify"`
		Timeout            uint   `json:"timeout" yaml:"timeout"`
		Interfaces         string `json:"interfaces" yaml:"interfaces"`
		InterfacesInterval uint   `json:"interfaces_interval" yaml:"interfaces_interval"`
		DNSServer          string `json:"dns_server" yaml:"dns_server"`
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

	InterfaceForm struct {
		Domain   *string `json:"domain,omitempty" yaml:"domain,omitempty"`
		Username *string `json:"username,omitempty" yaml:"username,omitempty"`
		UserType *string `json:"user_type,omitempty" yaml:"user_type,omitempty"`
		Password *string `json:"password,omitempty" yaml:"password,omitempty"`
	}

	InterfaceMeta struct {
		N           *string `json:"n,omitempty" yaml:"n,omitempty"`
		Type        *string `json:"type,omitempty" yaml:"type,omitempty"`
		Acid        *string `json:"acid,omitempty" yaml:"acid,omitempty"`
		Enc         *string `json:"enc,omitempty" yaml:"enc,omitempty"`
		OS          *string `json:"os,omitempty" yaml:"os,omitempty"`
		Name        *string `json:"name,omitempty" yaml:"name,omitempty"`
		InfoPrefix  *string `json:"info_prefix,omitempty" yaml:"info_prefix,omitempty"`
		DoubleStack *bool   `json:"double_stack,omitempty" yaml:"double_stack,omitempty"`
	}

	InterfaceConf struct {
		Form *InterfaceForm `json:"form,omitempty" yaml:"form,omitempty"`
		Meta *InterfaceMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	}
)

type SettingsConf struct {
	Basic            BasicConf                `json:"basic" yaml:"basic"`
	Guardian         GuardianConf             `json:"guardian" yaml:"guardian"`
	Backoff          BackoffConf              `json:"backoff" yaml:"backoff"`
	Log              LogConf                  `json:"log" yaml:"log"`
	DDNS             DdnsConf                 `json:"ddns" yaml:"ddns"`
	Reality          RealityConf              `json:"reality" yaml:"reality"`
	CustomHeader     map[string]interface{}   `json:"custom_header" yaml:"custom_header"`
	InterfacesConfig map[string]InterfaceConf `json:"interfaces_config,omitempty" yaml:"interfaces_config,omitempty"`
}
