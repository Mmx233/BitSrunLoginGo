package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	"time"
)

// Guardian 守护模式逻辑
func Guardian() {
	logger := config.Logger

	logger.Infoln("[以守护模式启动]")

	GuardianDuration := time.Duration(config.Settings.Guardian.Duration) * time.Second

	for {
		if config.Settings.Basic.Interfaces == "" { //单网卡
			err := Login(nil, true)
			if err != nil {
				logger.Errorln("登录出错: ", err)
			}
		} else { //多网卡
			interfaces, err := tools.GetInterfaceAddr(logger, config.Settings.Basic.Interfaces)
			if err == nil {
				for _, eth := range interfaces {
					logger.Debugf("使用 %s 网口登录 ", eth.Name)
					err = Login(&eth, true)
					if err != nil {
						logger.Errorln("网口 ", eth.Name+" 登录出错: ", err)
					}
				}
			}
		}
		time.Sleep(GuardianDuration)
	}
}
