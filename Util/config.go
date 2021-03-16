package Util

import (
	"Mmx/Modles"
	"fmt"
	"os"
)

type config struct {
	Path string
}

var Config = config{
	Path: "Config.json",
}

func (*config) Generate(Form *Modles.LoginForm, Meta *Modles.LoginMeta) *Modles.LoginInfo {
	return &Modles.LoginInfo{
		UrlLoginPage:       "http://" + Form.Domain + "/srun_portal_success",
		UrlGetChallengeApi: "http://" + Form.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + Form.Domain + "/cgi-bin/srun_portal",
		Meta:               Meta,
		Form: &Modles.LoginForm{
			UserName: Form.UserName + "@cmcc",
			PassWord: Form.PassWord,
		},
	}
}

func (a *config) Init() *Modles.LoginInfo {
	if !File.Exists(a.Path) {
		if err := File.Write(a.Path, &Modles.Config{ //默认值
			From: Modles.LoginForm{
				Domain:   "www.msftconnecttest.com",
				UserName: "",
				PassWord: "",
			},
			Meta: Modles.LoginMeta{
				N:     "200",
				VType: "1",
				Acid:  "5",
				Enc:   "srun_bx1",
			},
		}); err != nil {
			fmt.Println("Create 'Config.json' error:\n", err.Error())
			os.Exit(3)
		}
		fmt.Println("Please edit 'Config.json' and try again.")
		os.Exit(1)
	}

	var c Modles.Config
	if err := File.Read(a.Path, &c); err != nil {
		fmt.Println("Read config failed:\n", err.Error())
		os.Exit(3)
	}

	return a.Generate(
		&c.From,
		&c.Meta,
	)
}
