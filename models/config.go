package srunModels

import "github.com/Mmx233/BitSrunLoginGo/v1/transfer"

type Daemon struct {
	Enable bool   `json:"enable"`
	Path   string `json:"path"`
}

type Guardian struct {
	Enable   bool `json:"enable"`
	Duration uint `json:"duration"`
}

type Basic struct {
	Https        bool   `json:"https"`
	Timeout      uint   `json:"timeout"`
	Interfaces   string `json:"interfaces"`
	SkipNetCheck bool   `json:"skip_net_check" yaml:"skip_net_check"`
}

type Settings struct {
	Basic    Basic    `json:"basic"`
	Guardian Guardian `json:"guardian"`
	Daemon   Daemon   `json:"daemon"`
	Debug    Debug
}

type Debug struct {
	Enable   bool   `json:"enable"`
	WriteLog bool   `json:"write_log" yaml:"write_log"`
	LogPath  string `json:"log_path" yaml:"log_path"`
}

type Config struct {
	Form     srunTransfer.LoginForm `json:"form"`
	Meta     srunTransfer.LoginMeta `json:"meta"`
	Settings Settings               `json:"settings"`
}
