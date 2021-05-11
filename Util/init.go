package Util

import (
	"Mmx/Global"
	"Mmx/Modles"
	"log"
	"os"
)

func init() {
	//配置文件初始化
	Path := "Config.json"
	if !File.Exists(Path) {
		if err := File.WriteJson(Path, &Modles.Config{ //默认值
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
			Settings: Modles.Settings{
				Dns: "1.2.4.8",
			},
		}); err != nil {
			log.Println("创建配置文件失败:\n", err.Error())
			os.Exit(1)
		}
		log.Println("已生成配置文件，请编辑 'Config.json' 然后重试")
		os.Exit(0)
	}

	var c Modles.Config
	if err := File.ReadJson(Path, &c); err != nil {
		log.Println("读取配置文件失败:\n", err.Error())
		os.Exit(1)
	}

	Global.Config = &c
}
