package srunModels

import "github.com/Mmx233/BitSrunLoginGo/v1/transfer"

type Daemon struct {
	Enable bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
	Path   string `json:"path" yaml:"path" mapstructure:"path"`
}

type Guardian struct {
	Enable   bool `json:"enable" yaml:"enable" mapstructure:"enable"`
	Duration uint `json:"duration" yaml:"duration" mapstructure:"duration"`
}

type Basic struct {
	Https          bool   `json:"https" yaml:"https" mapstructure:"https"`
	SkipCertVerify bool   `json:"skip_cert_verify" yaml:"skip_cert_verify" mapstructure:"skip_cert_verify"`
	Timeout        uint   `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Interfaces     string `json:"interfaces" yaml:"interfaces" mapstructure:"interfaces"`
	SkipNetCheck   bool   `json:"skip_net_check" yaml:"skip_net_check" mapstructure:"skip_net_check"`
}

type Settings struct {
	Basic    Basic    `json:"basic" yaml:"basic" mapstructure:"basic"`
	Guardian Guardian `json:"guardian" yaml:"guardian" mapstructure:"guardian"`
	Daemon   Daemon   `json:"daemon" yaml:"daemon" mapstructure:"daemon"`
	Debug    Debug    `json:"debug" yaml:"debug" mapstructure:"debug"`
}

type Debug struct {
	Enable   bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
	WriteLog bool   `json:"write_log" yaml:"write_log" mapstructure:"write_log"`
	LogPath  string `json:"log_path" yaml:"log_path" mapstructure:"log_path"`
}

type Config struct {
	Form     srunTransfer.LoginForm `json:"form" yaml:"form" mapstructure:"form"`
	Meta     srunTransfer.LoginMeta `json:"meta" yaml:"meta" mapstructure:"meta"`
	Settings Settings               `json:"settings" yaml:"settings" mapstructure:"settings"`
}
