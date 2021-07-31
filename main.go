package main

import (
	"autoLogin/controllers"
	"autoLogin/global"
	"autoLogin/util"
	"os"
)

func main() {
	defer util.Log.CatchRecover()

	if err := controllers.Login(true); err != nil {
		util.Log.Println("运行出错，状态异常")
		if global.Config.Settings.ForceGuardianEvenFailed {
			util.Log.Println(err)
		} else if global.Config.Settings.DemoMode {
			util.Log.Fatalln(err)
		} else {
			util.Log.Println(err)
			os.Exit(3)
		}
	}

	controllers.EnterGuardian()
}
