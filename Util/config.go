package Util

import (
	"Mmx/Global"
	"Mmx/Modles"
	"fmt"
	"os"
	"path/filepath"
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
		UrlCheckApi:        "http://" + Form.Domain + "/cgi-bin/rad_user_info",
		Meta:               Meta,
		Form: &Modles.LoginForm{
			UserName: Form.UserName + "@cmcc",
			PassWord: Form.PassWord,
		},
	}
}

func (a *config) Init() *Modles.LoginInfo {
	if t, err := os.Executable(); err != nil {
		ErrHandler(err)
	} else {
		a.Path = filepath.Dir(t) + "/" + a.Path
	}
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
			fmt.Println("创建配置文件失败:\n", err.Error())
			os.Exit(3)
		}
		fmt.Println("已生成配置文件，请编辑 'Config.json' 然后重试")
		os.Exit(1)
	}

	var c Modles.Config
	if err := File.Read(a.Path, &c); err != nil {
		fmt.Println("读取配置文件失败:\n", err.Error())
		os.Exit(3)
	}

	Global.Config = &c

	return a.Generate(
		&c.From,
		&c.Meta,
	)
}
