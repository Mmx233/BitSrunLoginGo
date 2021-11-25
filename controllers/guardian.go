package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"os"
	"os/exec"
	"time"
)

func Guardian(output bool) {
	global.Status.Output = output

	go Daemon.DaemonChan()

	if e := Daemon.MarkDaemon(); e != nil {
		util.Log.Fatalln(e)
	}

	var c = make(chan bool)
	for {
		global.Status.Output = output
		go func() {
			defer func() {
				_ = recover()
			}()
			if !util.Checker.NetOk() {
				util.Log.Println("Network down, trying to login")
				_ = Login(output, true)
			} else {
				if global.Config.Settings.DemoMode {
					util.Log.Println("Network ok")
				}
			}
			c <- false
		}()
		<-c
		time.Sleep(time.Duration(global.Config.Settings.Guardian) * time.Second)
	}
}

func EnterGuardian() {
	global.Status.Output = true
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
