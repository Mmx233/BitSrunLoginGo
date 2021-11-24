package main

import (
	"github.com/Mmx233/BitSrunLoginGo/controllers"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
)

func main() {
	defer util.Log.CatchRecover()

	if global.Config.Settings.Guardian != 0 {
		controllers.EnterGuardian()
	} else if err := controllers.Login(true, false); err != nil {
		util.Log.Println("运行出错，状态异常")
		if global.Config.Settings.DemoMode {
			util.Log.Fatalln(err)
		} else {
			util.Log.Println(err)
			return
		}
	}
}
