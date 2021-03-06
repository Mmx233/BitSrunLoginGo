package models

type Settings struct {
	Timeout  uint `json:"timeout"`
	DemoMode bool `json:"demo_mode"`
	Guardian uint `json:"guardian"`
	Daemon   bool `json:"daemon"`
}

type Config struct {
	From     LoginForm `json:"from"`
	Meta     LoginMeta `json:"meta"`
	Settings Settings  `json:"settings"`
}

func (a *Config) Generate() *LoginInfo {
	return &LoginInfo{
		UrlLoginPage:       "http://" + a.From.Domain + "/srun_portal_success",
		UrlGetChallengeApi: "http://" + a.From.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + a.From.Domain + "/cgi-bin/srun_portal",
		UrlCheckApi:        "http://" + a.From.Domain + "/cgi-bin/rad_user_info",
		Meta:               &a.Meta,
		Form: &LoginForm{
			UserName: a.From.UserName + "@" + a.From.UserType,
			PassWord: a.From.PassWord,
		},
	}
}
