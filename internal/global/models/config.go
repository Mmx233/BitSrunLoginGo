package models

import (
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
)

type Guardian struct {
	Enable   bool `yaml:"enable"`
	Duration uint `yaml:"duration"`
}

type Basic struct {
	Https          bool   `yaml:"https"`
	SkipCertVerify bool   `yaml:"skip_cert_verify"`
	Timeout        uint   `yaml:"timeout"`
	Interfaces     string `yaml:"interfaces"`
}

type Log struct {
	DebugLevel bool   `yaml:"debug_level"`
	WriteFile  bool   `yaml:"write_file"`
	FilePath   string `yaml:"log_path"`
	FileName   string `yaml:"log_name"`
}

type DDNS struct {
	Enable   bool                   `yaml:"enable"`
	TTL      uint                   `yaml:"ttl"`
	Domain   string                 `yaml:"domain"`
	Provider string                 `yaml:"provider"`
	Config   map[string]interface{} `yaml:",inline"`
}

type Settings struct {
	Basic    Basic    `yaml:"basic"`
	Guardian Guardian `yaml:"guardian"`
	Log      Log      `yaml:"log"`
	DDNS     DDNS     `yaml:"ddns"`
}

type Config struct {
	Form     srun.LoginForm `yaml:"form"`
	Meta     srun.LoginMeta `yaml:"meta"`
	Settings Settings       `yaml:"settings"`
}
