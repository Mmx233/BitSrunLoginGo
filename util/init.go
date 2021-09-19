package util

import (
	"autoLogin/global"
	"autoLogin/models"
	"github.com/Mmx233/config"
	"github.com/Mmx233/tool"
	"log"
	"os"
	"time"
)

func init() {
	//配置文件初始化
	if e := config.Load(config.Options{
		Config: &global.Config,
		Default: &models.Config{
			From: models.LoginForm{
				Domain:   "www.msftconnecttest.com",
				UserType: "cmcc",
			},
			Meta: models.LoginMeta{
				N:    "200",
				Type: "1",
				Acid: "5",
				Enc:  "srun_bx1",
			},
			Settings: models.Settings{
				Timeout: 1,
			},
		},
		Path:        "Config.json",
		FillDefault: true,
		Overwrite:   true,
	}); e != nil {
		if config.IsNew(e) {
			log.Println("已生成配置文件，请编辑 'Config.json' 然后重试")
			os.Exit(0)
		}
		log.Println("读取配置文件失败:\n", e.Error())
		os.Exit(1)
	}

	//http工具设定
	tool.HTTP.Options.Timeout = 3 * time.Second
}
