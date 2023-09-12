package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"time"

	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
)

// Guardian 守护模式逻辑
func Guardian() {
	log.Infoln("[以守护模式启动]")

	GuardianDuration := time.Duration(config.Settings.Guardian.Duration) * time.Second

	var c = make(chan bool)
	for {
		go func() {
			if config.Settings.Basic.Interfaces == "" { //单网卡
				err := Login(nil, true)
				if err != nil {
					log.Errorln("登录出错: ", err)
				}
			} else { //多网卡
				interfaces, err := tools.GetInterfaceAddr(config.Settings.Basic.Interfaces)
				if err == nil {
					for _, eth := range interfaces {
						log.Debugf("使用 %s 网口登录 ", eth.Name)
						err = Login(&eth, true)
						if err != nil {
							log.Errorln("网口 ", eth.Name+" 登录出错: ", err)
						}
					}
				}
			}

			c <- false
		}()
		<-c
		time.Sleep(GuardianDuration)
	}
}
