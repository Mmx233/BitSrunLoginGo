package main

import (
	"Mmx/Util"
	"Mmx/controllers"
)

func main() {
	defer util.Log.CatchRecover()

	if err := controllers.Login(true); err != nil {
		util.ErrHandler(err)
		return
	}

	controllers.EnterGuardian()
}
