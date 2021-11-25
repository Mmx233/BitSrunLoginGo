package global

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"github.com/Mmx233/config"
	"log"
	"os"
)

var Config srunModels.Config

func init() {
	initFlags()

	//配置文件初始化
	if e := config.Load(config.Options{
		Config: &Config,
		Default: &srunModels.Config{
			Form: srunModels.LoginForm{
				Domain:   "www.msftconnecttest.com",
				UserType: "cmcc",
			},
			Meta: srunModels.LoginMeta{
				N:    "200",
				Type: "1",
				Acid: "5",
				Enc:  "srun_bx1",
			},
			Settings: srunModels.Settings{
				Timeout: 1,
				Daemon: srunModels.Daemon{
					Path: ".autoLogin",
				},
				Guardian: srunModels.Guardian{
					Duration: 300,
				},
			},
		},
		Path:        Flags.Path,
		FillDefault: true,
		Overwrite:   true,
	}); e != nil {
		if config.IsNew(e) {
			log.Println("已生成配置文件，请编辑 '" + Flags.Path + "' 然后重试")
			os.Exit(0)
		}
		log.Println("读取配置文件失败:\n", e.Error())
		os.Exit(1)
	}

	util.Log.Demo = Config.Settings.DemoMode
}
