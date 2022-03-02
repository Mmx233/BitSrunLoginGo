package global

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
	"github.com/Mmx233/tool"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Config srunModels.Config

var Timeout time.Duration

func readConfig() error {
	//配置文件默认值
	viper.SetDefault("form", srunTransfer.LoginForm{
		Domain:   "www.msftconnecttest.com",
		UserType: "cmcc",
	})
	viper.SetDefault("meta", srunTransfer.LoginMeta{
		N:    "200",
		Type: "1",
		Acid: "5",
		Enc:  "srun_bx1",
	})
	viper.SetDefault("settings", srunModels.Settings{
		Basic: srunModels.Basic{
			Timeout: 5,
		},
		Daemon: srunModels.Daemon{
			Path: ".autoLogin",
		},
		Guardian: srunModels.Guardian{
			Duration: 300,
		},
		Debug: srunModels.Debug{
			LogPath: "./",
		},
	})

	//生成配置文件
	if !tool.File.Exists(Flags.Path) {
		e := viper.WriteConfigAs(Flags.Path)
		if e != nil {
			log.Println("生成配置文件失败：", e)
			return e
		}
		log.Println("已生成配置文件，请编辑 '" + Flags.Path + "' 然后重试")
		os.Exit(0)
	}

	//读取配置文件
	viper.SetConfigFile(Flags.Path)
	if e := viper.ReadInConfig(); e != nil {
		log.Println("读取配置文件失败：", e)
		return e
	}
	if e := viper.Unmarshal(&Config); e != nil {
		log.Println("解析配置文件失败：", e)
		return e
	}

	return nil
}

func init() {
	initFlags()

	//配置文件初始化
	if readConfig() != nil {
		os.Exit(1)
	}

	Timeout = time.Duration(Config.Settings.Basic.Timeout) * time.Second
	initTransport()
}
