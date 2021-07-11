package main

import (
	"Mmx/Util"
	"Mmx/controllers"
	"Mmx/global"
	"flag"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	goDaemon := flag.Bool("daemon", false, "")
	flag.Parse()
	if *goDaemon {
		controllers.Guardian()
	}

	defer func() {
		if e := recover(); e != nil {
			util.Log.Println(e)
			var buf [4096]byte
			util.Log.Println(string(buf[:runtime.Stack(buf[:], false)]))
			os.Exit(1)
		}
	}()
	if err := controllers.Login(true); err != nil {
		util.ErrHandler(err)
		return
	}

	if global.Config.Settings.Guardian != 0 {
		global.Status.Daemon = true
		util.Log.Println("[Guardian mode]")
		if global.Config.Settings.Daemon {
			if err := exec.Command(os.Args[0], "-daemon").Start(); err != nil {
				util.ErrHandler(err)
				return
			}
			util.Log.Println("[Daemon mode entered]")
			return
		}
		controllers.Guardian()
	}
}
