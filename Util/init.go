package Util

import (
	"Mmx/Global"
	"Mmx/Modles"
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	//配置文件初始化
	Path := "Config.json"
	if t, err := os.Executable(); err != nil {
		ErrHandler(err)
	} else {
		Path = filepath.Dir(t) + "/" + Path
	}
	if !File.Exists(Path) {
		if err := File.Write(Path, &Modles.Config{ //默认值
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
	if err := File.Read(Path, &c); err != nil {
		fmt.Println("读取配置文件失败:\n", err.Error())
		os.Exit(3)
	}

	Global.Config = &c
}
