package models

import (
	"autoLogin/models/util"
	"reflect"
)

type Settings struct {
	QuitIfNetOk bool   `json:"quit_if_net_ok"`
	DemoMode    bool   `json:"demo_mode"`
	Dns         string `json:"dns"`
	Guardian    uint   `json:"guardian"`
	Daemon      bool   `json:"daemon"`
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

func (a *Config) FillDefault() *Config {
	var m = map[interface{}]map[string]interface{}{
		&a.From: {
			"Domain":   "www.msftconnecttest.com",
			"UserType": "cmcc",
		},
		&a.Meta: {
			"N":    "200",
			"Type": "1",
			"Acid": "5",
			"Enc":  "srun_bx1",
		},
		&a.Settings: {
			"Dns": "1.2.4.8",
		},
	}

	for q, w := range m {
		t := reflect.ValueOf(q).Elem()
		for k, v := range w {
			tt := t.FieldByName(k)
			if util.Reflect.IsEmpty(tt) {
				tt.Set(reflect.ValueOf(v))
			}
		}
	}

	return a
}
