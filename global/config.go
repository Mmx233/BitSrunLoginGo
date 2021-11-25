package global

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/config"
	"log"
	"os"
)

var Config models.Config

func init() {
	initFlags()

	//配置文件初始化
	if e := config.Load(config.Options{
		Config: &Config,
		Default: &models.Config{
			Form: models.LoginForm{
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
				Daemon: models.Daemon{
					Path: ".autoLogin",
				},
				Guardian: models.Guardian{
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
}
