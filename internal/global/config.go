package global

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/global/models"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/tool"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var Config models.Config

var Timeout time.Duration

func createDefaultConfig() error {
	configFile, err := os.OpenFile(Flags.Path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer configFile.Close()

	return yaml.NewEncoder(configFile).Encode(&models.Config{
		Form: srun.LoginForm{
			Domain:   "www.msftconnecttest.com",
			UserType: "cmcc",
		},
		Meta: srun.LoginMeta{
			N:    "200",
			Type: "1",
			Acid: "5",
			Enc:  "srun_bx1",
		},
		Settings: models.Settings{
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
		},
	})
}

func initConfig() {
	// 生成配置文件
	if exist, err := tool.File.Exists(Flags.Path); err != nil {
		log.Fatalln("[init] 读取配置文件失败：", err)
	} else if !exist {
		err = createDefaultConfig()
		if err != nil {
			log.Fatalln("[init] 生成配置文件失败：", err)
		}
		log.Infoln("[init] 已生成配置文件，请编辑 '" + Flags.Path + "' 然后重试")
		os.Exit(0)
	}

	// 读取配置文件
	configBytes, err := os.ReadFile(Flags.Path)
	if err != nil {
		log.Fatalln("[init] 读取配置失败：", err)
	}
	if err = yaml.Unmarshal(configBytes, &Config); err != nil {
		log.Fatalln("[init] 解析配置失败：", err)
	}

	// flag 配置覆写
	if Flags.Debug {
		Config.Settings.Log.DebugLevel = true
	}
	if Flags.Acid != "" {
		Config.Meta.Acid = Flags.Acid
	}
}
