package main

import (
	"github.com/Mmx233/BitSrunLoginGo/controllers"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	if global.Flags.RunningDaemon {
		//后台挂起模式中
		controllers.Guardian()
	} else if global.Config.Settings.Guardian.Enable {
		//进入守护模式流程
		controllers.EnterGuardian()
	} else {
		if global.Config.Settings.Basic.Interfaces == "" { //单网卡
			if err := controllers.Login(nil); err != nil {
				log.Fatalln("运行出错，状态异常: ", err)
			}
		} else { //多网卡
			log.Debugln("多网卡模式")
			interfaces, _ := util.GetInterfaceAddr()
			for _, eth := range interfaces {
				log.Infoln("网卡: ", eth.Name)
				if err := controllers.Login(eth.Addr); err != nil {
					log.Errorln("运行出错，状态异常: ", err)
				}
			}
		}
	}
}
