package config

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/tool"
	"os"
	"time"
)

type ConfFromFile struct {
	Form     srun.LoginForm `json:"form" yaml:"form"`
	Meta     srun.LoginMeta `json:"meta" yaml:"meta"`
	Settings SettingsConf   `json:"settings" yaml:"settings"`
}

var (
	Form     *srun.LoginForm
	Meta     *srun.LoginMeta
	Settings *SettingsConf

	Timeout time.Duration
)

func init() {
	initLogPre()

	logger := Logger.WithField(keys.LogComponent, "init")
	reader := newReaderFromPath(flags.Path)

	// 生成配置文件
	exist, err := tool.File.Exists(flags.Path)
	if err != nil {
		logger.Fatalln("读取配置文件失败：", err)
	} else if !exist {
		var data []byte
		data, err = reader.Marshal(&defaultConfig)
		if err != nil {
			logger.Fatalln("生成配置文件失败：", err)
		}
		if err = os.WriteFile(flags.Path, data, 0600); err != nil {
			logger.Fatalln("写入配置文件失败：", err)
		}
		logger.Infoln("已生成配置文件，请编辑 '" + flags.Path + "' 然后重试")
		os.Exit(0)
	}

	// 读取配置文件
	data, err := os.ReadFile(flags.Path)
	if err != nil {
		logger.Fatalln("读取配置失败：", err)
	}
	var fileConf ConfFromFile
	if err = reader.Unmarshal(data, &fileConf); err != nil {
		logger.Fatalln("解析配置失败：", err)
	}
	Form = &fileConf.Form
	Meta = &fileConf.Meta
	Settings = &fileConf.Settings
	Timeout = time.Duration(Settings.Basic.Timeout) * time.Second

	// flag 配置覆写
	if flags.Debug {
		Settings.Log.DebugLevel = true
	}
	if flags.Acid != "" {
		Meta.Acid = flags.Acid
	}

	initLogFinal()
	initBackoff()
}
