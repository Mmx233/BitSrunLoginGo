package main

import (
	"github.com/Mmx233/BitSrunLoginGo/controllers"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"log"
)

func main() {
	if e := util.Log.Init(
		global.Config.Settings.Debug.Enable,
		global.Config.Settings.Debug.WriteLog,
		true,
		global.Config.Settings.Debug.LogPath,
	); e != nil {
		log.Fatalln("初始化日志失败: ", e)
	}
	defer util.Log.CatchRecover()

	if global.Flags.RunningDaemon {
		//后台挂起模式中
		controllers.Guardian(false)
	} else if global.Config.Settings.Guardian.Enable {
		//进入守护模式流程
		controllers.EnterGuardian()
	} else {
		//单次登录模式
		if global.Config.Settings.Basic.Interfaces == "" { //单网卡
			if err := controllers.Login(true, global.Config.Settings.Basic.SkipNetCheck, nil); err != nil {
				util.Log.Fatal("运行出错，状态异常: ", err)
			}
		} else { //多网卡
			util.Log.Debug("多网卡模式")
			interfaces, _ := util.GetInterfaceAddr()
			for _, eth := range interfaces {
				util.Log.Info("网卡: ", eth.Name)
				if err := controllers.Login(true, global.Config.Settings.Basic.SkipNetCheck, eth.Addr); err != nil {
					util.Log.Warn("运行出错，状态异常: ", err)
				}
			}
		}
	}
}
