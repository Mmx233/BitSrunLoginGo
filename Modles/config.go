package Modles

type Settings struct {
	QuitIfNetOk bool   `json:"quit_if_net_ok"`
	DemoMode    bool   `json:"demo_mode"`
	Dns         string `json:"dns"`
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
			UserName: a.From.UserName + "@cmcc",
			PassWord: a.From.PassWord,
		},
	}
}
