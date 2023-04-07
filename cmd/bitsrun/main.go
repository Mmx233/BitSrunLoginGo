package main

import (
	"fmt"
	controllers2 "github.com/Mmx233/BitSrunLoginGo/internal/controllers"
	global2 "github.com/Mmx233/BitSrunLoginGo/internal/global"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
)

func main() {
	if global2.Flags.RunningDaemon {
		//后台挂起模式中
		controllers2.Guardian()
	} else if global2.Config.Settings.Guardian.Enable {
		//进入守护模式流程
		controllers2.EnterGuardian()
	} else {
		//登录流程
		var err error
		if global2.Config.Settings.Basic.Interfaces == "" { //单网卡
			if err = controllers2.Login(nil, false); err != nil {
				log.Errorln("登录出错: ", err)
				if !global2.Config.Settings.Log.DebugLevel {
					fmt.Printf("开启调试日志（debug_level）获取详细信息")
				}
				return
			}
		} else { //多网卡
			log.Infoln("多网卡模式")
			interfaces, _ := tools.GetInterfaceAddr()
			for _, eth := range interfaces {
				log.Infoln("使用网卡: ", eth.Name)
				if err = controllers2.Login(eth.Addr, false); err != nil {
					log.Errorf("网卡 %s 登录出错: %v", eth.Name, err)
				}
			}
		}
	}
}
