package main

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/internal/controllers"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
	"net"
)

func main() {
	if config.Settings.Guardian.Enable {
		//进入守护模式
		controllers.Guardian()
	} else {
		//登录流程
		var err error
		if config.Settings.Basic.Interfaces == "" { //单网卡
			var eth *tools.Eth
			if flags.Interface != "" {
				netEth, e := net.InterfaceByName(flags.Interface)
				if e != nil {
					log.Warnf("获取指定网卡 %s 失败，使用默认网卡: %v", flags.Interface, e)
				} else {
					eth, e = tools.ConvertInterface(*netEth)
					if e != nil {
						log.Warnf("获取指定网卡 %s ip 地址失败，使用默认网卡: %v", flags.Interface, e)
					} else if eth == nil {
						log.Warnf("指定网卡 %s 无可用 ip 地址，使用默认网卡", flags.Interface)
					} else {
						log.Debugf("使用指定网卡 %s ip: %s", eth.Name, eth.Addr.String())
					}
				}
			}
			if err = controllers.Login(eth, false); err != nil {
				log.Errorln("登录出错: ", err)
				if !config.Settings.Log.DebugLevel {
					log.Infoln("开启调试日志（debug_level）获取详细信息")
				}
				return
			}
		} else { //多网卡
			log.Infoln("多网卡模式")
			interfaces, _ := tools.GetInterfaceAddr(config.Settings.Basic.Interfaces)
			for _, eth := range interfaces {
				log.Infoln("使用网卡: ", eth.Name)
				if err = controllers.Login(&eth, false); err != nil {
					log.Errorf("网卡 %s 登录出错: %v", eth.Name, err)
				}
			}
		}
	}
}
