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
			defer func() {
				_ = recover()
			}()
			if config.Settings.Basic.Interfaces == "" { //单网卡
				e := Login(nil, true)
				if e != nil {
					log.Errorln("登录出错: ", e)
				}
			} else { //多网卡
				interfaces, e := tools.GetInterfaceAddr(config.Settings.Basic.Interfaces)
				if e == nil {
					for _, eth := range interfaces {
						log.Debugf("使用 %s 网口登录 ", eth.Name)
						e = Login(&eth, true)
						if e != nil {
							log.Errorln("网口 ", eth.Name+" 登录出错: ", e)
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
