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
	Timeout      uint   `json:"timeout"`
	Interfaces   string `json:"interfaces"`
	DemoMode     bool   `json:"demo_mode" yaml:"demo_mode"`
	SkipNetCheck bool   `json:"skip_net_check" yaml:"skip_net_check"`
}

type Settings struct {
	Basic    Basic
	Guardian Guardian
	Daemon   Daemon
}

type Config struct {
	Form     srunTransfer.LoginForm `json:"form"`
	Meta     srunTransfer.LoginMeta `json:"meta"`
	Settings Settings               `json:"settings"`
}
