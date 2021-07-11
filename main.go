package main

import (
	"autoLogin/controllers"
	"autoLogin/util"
)

func main() {
	defer util.Log.CatchRecover()

	if err := controllers.Login(true); err != nil {
		util.ErrHandler(err)
		return
	}

	controllers.EnterGuardian()
}
