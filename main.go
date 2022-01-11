package main

import (
	"github.com/Mmx233/BitSrunLoginGo/controllers"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
)

func main() {
	util.Log.Demo = global.Config.Settings.DemoMode
	util.Log.OutPut = true
	defer util.Log.CatchRecover()

	if global.Flags.RunningDaemon {
		//后台挂起模式中
		controllers.Guardian(false)
	} else if global.Config.Settings.Guardian.Enable {
		//进入守护模式流程
		controllers.EnterGuardian()
	} else {
		//单次登录模式
		if global.Config.Settings.Interfaces == "" { //单网卡
			if err := controllers.Login(true, false, nil); err != nil {
				util.Log.Println("运行出错，状态异常")
				if global.Config.Settings.DemoMode {
					util.Log.Fatalln(err)
				} else {
					util.Log.Println(err)
					return
				}
			}
		} else { //多网卡
			interfaces, e := util.GetInterfaceAddr()
			if e != nil {
				return
			}
			for _, eth := range interfaces {
				util.Log.Println(eth.Name)
				if err := controllers.Login(true, false, eth.Addr); err != nil {
					util.Log.Println(eth.Name + "运行出错，状态异常")
					util.Log.Println(err)
				}
			}
		}
	}
}
