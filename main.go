package main

import (
	"Mmx/controllers"
	"Mmx/util"
)

func main() {
	defer util.Log.CatchRecover()

	if err := controllers.Login(true); err != nil {
		util.ErrHandler(err)
		return
	}

	controllers.EnterGuardian()
}
