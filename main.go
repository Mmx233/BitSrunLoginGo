package main

import (
	"autoLogin/controllers"
	"autoLogin/global"
	"autoLogin/util"
)

func main() {
	defer util.Log.CatchRecover()

	if global.Config.Settings.Guardian != 0 {
		controllers.EnterGuardian()
	} else if err := controllers.Login(true); err != nil {
		util.Log.Println("运行出错，状态异常")
		if global.Config.Settings.DemoMode {
			util.Log.Fatalln(err)
		} else {
			util.Log.Println(err)
			return
		}
	}
}
