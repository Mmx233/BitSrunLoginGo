package srunModels

type Daemon struct {
	Enable bool   `json:"enable"`
	Path   string `json:"path"`
}

type Guardian struct {
	Enable   bool `json:"enable"`
	Duration uint `json:"duration"`
}

type Settings struct {
	Timeout  uint `json:"timeout"`
	DemoMode bool `json:"demo_mode"`
	Guardian Guardian
	Daemon   Daemon
}

type Config struct {
	Form     LoginForm `json:"form"`
	Meta     LoginMeta `json:"meta"`
	Settings Settings  `json:"settings"`
}
