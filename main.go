package main

import (
	"github.com/Mmx233/BitSrunLoginGo/controllers"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
)

func main() {
	defer util.Log.CatchRecover()

	if global.Flags.RunningDaemon {
		//后台挂起模式中
		controllers.Guardian(false)
	} else if global.Config.Settings.Guardian.Enable {
		//进入守护模式流程
		controllers.EnterGuardian()
	} else if err := controllers.Login(true, false); err != nil { //单次登录模式
		util.Log.Println("运行出错，状态异常")
		if global.Config.Settings.DemoMode {
			util.Log.Fatalln(err)
		} else {
			util.Log.Println(err)
			return
		}
	}
}
