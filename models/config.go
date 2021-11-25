package models

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

func (a *Config) Generate() *LoginInfo {
	return &LoginInfo{
		UrlLoginPage:       "http://" + a.Form.Domain + "/srun_portal_success",
		UrlGetChallengeApi: "http://" + a.Form.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + a.Form.Domain + "/cgi-bin/srun_portal",
		UrlCheckApi:        "http://" + a.Form.Domain + "/cgi-bin/rad_user_info",
		Meta:               &a.Meta,
		Form: &LoginForm{
			UserName: a.Form.UserName + "@" + a.Form.UserType,
			PassWord: a.Form.PassWord,
		},
	}
}
