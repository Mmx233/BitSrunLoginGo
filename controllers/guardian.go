package controllers

import (
	"autoLogin/global"
	"autoLogin/util"
	"os"
	"os/exec"
	"time"
)

func Guardian(output bool) {
	global.Status.Output = output

	if e := Daemon.MarkDaemon(); e != nil {
		util.Log.Fatalln(e)
	}

	var c = make(chan bool)
	for {
		if !Daemon.CheckDaemon() {
			os.Exit(1)
		}
		global.Status.Output = output
		time.Sleep(time.Duration(global.Config.Settings.Guardian) * time.Second)
		go func() {
			defer func() {
				_ = recover()
			}()
			if !util.Checker.NetOk() {
				util.Log.Println("Network down, trying to login")
				_ = Login(false)
			} else {
				if global.Config.Settings.DemoMode {
					util.Log.Println("Network ok")
				}
			}
			c <- false
		}()
		<-c
	}
}

func EnterGuardian() {
	if global.Config.Settings.Guardian != 0 {
		global.Status.Guardian = true
		util.Log.Println("[Guardian mode]")
		if global.Config.Settings.Daemon {
			if err := exec.Command(os.Args[0], "-daemon").Start(); err != nil {
				util.Log.Fatalln(err)
			}
			util.Log.Println("[Daemon mode entered]")
			return
		}
		Guardian(true)
	}
}
