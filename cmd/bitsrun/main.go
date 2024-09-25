package main

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/internal/controllers"
)

func main() {
	logger := config.Logger
	if config.Settings.Basic.Interfaces != "" {
		logger.Infoln("[多网卡模式]")
	}

	if config.Settings.Guardian.Enable {
		//进入守护模式
		controllers.Guardian(logger.WithField(keys.LogComponent, "guard"))
	} else {
		//执行单次流程
		_ = controllers.Login(controllers.LoginConf{
			Logger:                      logger.WithField(keys.LogComponent, "login"),
			IsOnlineDetectLogDebugLevel: false,
		})
	}
}
