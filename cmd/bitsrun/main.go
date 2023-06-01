package main

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/controllers"
	"github.com/Mmx233/BitSrunLoginGo/internal/global"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	if global.Config.Settings.Guardian.Enable {
		//进入守护模式
		controllers.Guardian()
	} else {
		//登录流程
		var err error
		if global.Config.Settings.Basic.Interfaces == "" { //单网卡
			var eth *tools.Eth
			if global.Flags.Interface != "" {
				netEth, e := net.InterfaceByName(global.Flags.Interface)
				if e != nil {
					log.Warnf("获取指定网卡 %s 失败，使用默认网卡: %v", global.Flags.Interface, e)
				} else {
					eth, e = tools.ConvertInterface(*netEth)
					if e != nil {
						log.Warnf("获取指定网卡 %s ip 地址失败，使用默认网卡: %v", global.Flags.Interface, e)
					} else if eth == nil {
						log.Warnf("指定网卡 %s 无可用 ip 地址，使用默认网卡", global.Flags.Interface)
					}
				}
			}
			if err = controllers.Login(eth, false); err != nil {
				log.Errorln("登录出错: ", err)
				if !global.Config.Settings.Log.DebugLevel {
					log.Infoln("开启调试日志（debug_level）获取详细信息")
				}
				return
			}
		} else { //多网卡
			log.Infoln("多网卡模式")
			interfaces, _ := tools.GetInterfaceAddr(global.Config.Settings.Basic.Interfaces)
			for _, eth := range interfaces {
				log.Infoln("使用网卡: ", eth.Name)
				if err = controllers.Login(&eth, false); err != nil {
					log.Errorf("网卡 %s 登录出错: %v", eth.Name, err)
				}
			}
		}
	}
}
