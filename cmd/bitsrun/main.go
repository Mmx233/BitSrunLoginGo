package main

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/controllers"
	"github.com/Mmx233/BitSrunLoginGo/tools"
)

func main() {
	if config.Settings.Guardian.Enable {
		//进入守护模式
		controllers.Guardian()
	} else {
		//登录流程
		var err error
		logger := config.Logger
		if config.Settings.Basic.Interfaces == "" { //单网卡
			if err = controllers.Login(nil, false); err != nil {
				logger.Errorln("登录出错: ", err)
				if !config.Settings.Log.DebugLevel {
					logger.Infoln("开启调试日志 (debug_level) 获取详细信息")
				}
				return
			}
		} else { //多网卡
			logger.Infoln("多网卡模式")
			interfaces, _ := tools.GetInterfaceAddr(logger, config.Settings.Basic.Interfaces)
			for _, eth := range interfaces {
				logger.Infoln("使用网卡: ", eth.Name)
				if err = controllers.Login(&eth, false); err != nil {
					config.Logger.Errorf("网卡 %s 登录出错: %v", eth.Name, err)
				}
			}
		}
	}
}
