package global

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/global/models"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

var Config models.Config

var Timeout time.Duration

func readConfig() {
	//配置文件默认值
	viper.SetDefault("form", srun.LoginForm{
		Domain:   "www.msftconnecttest.com",
		UserType: "cmcc",
	})
	viper.SetDefault("meta", srun.LoginMeta{
		N:    "200",
		Type: "1",
		Acid: "5",
		Enc:  "srun_bx1",
	})
	viper.SetDefault("settings", models.Settings{
		Basic: models.Basic{
			Timeout: 5,
		},
		Guardian: models.Guardian{
			Duration: 300,
		},
		Log: models.Log{
			FilePath: "./",
		},
		DDNS: models.DDNS{
			Enable: false,
			TTL:    600,
			Domain: "www.example.com",
		},
	})

	//生成配置文件
	if exist, e := tool.File.Exists(Flags.Path); e != nil {
		log.Fatalln("[init] 读取配置文件失败：", e)
	} else if !exist {
		e = viper.WriteConfigAs(Flags.Path)
		if e != nil {
			log.Fatalln("[init] 生成配置文件失败：", e)
		}
		log.Infoln("[init] 已生成配置文件，请编辑 '" + Flags.Path + "' 然后重试")
		os.Exit(0)
	}

	//读取配置文件
	viper.SetConfigFile(Flags.Path)
	if e := viper.ReadInConfig(); e != nil {
		log.Fatalln("[init] 读取配置失败：", e)
	}
	if e := viper.Unmarshal(&Config); e != nil {
		log.Fatalln("[init] 解析配置失败：", e)
	}
}
