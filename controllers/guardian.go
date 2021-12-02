package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"os"
	"os/exec"
	"time"
)

func Guardian(output bool) {
	util.Log.OutPut = output

	if global.Config.Settings.Daemon.Enable {
		go Daemon.DaemonChan()

		if e := Daemon.MarkDaemon(); e != nil {
			util.Log.Fatalln(e)
		}
	}

	var c = make(chan bool)
	for {
		util.Log.OutPut = output
		go func() {
			defer func() {
				_ = recover()
			}()
			if !util.Checker.NetOk(global.Config.Settings.Timeout) {
				util.Log.Println("Network down, trying to login")
				e := Login(output, true)
				if e != nil {
					util.Log.Println("Error: ", e)
				}
			} else {
				if global.Config.Settings.DemoMode {
					util.Log.Println("Network ok")
				}
			}
			c <- false
		}()
		<-c
		time.Sleep(time.Duration(global.Config.Settings.Guardian.Duration) * time.Second)
	}
}

func EnterGuardian() {
	util.Log.OutPut = true
	util.Log.Println("[Guardian mode]")
	if global.Config.Settings.Daemon.Enable || global.Flags.Daemon {
		if err := exec.Command(os.Args[0], append(os.Args[1:], "--running-daemon")...).Start(); err != nil {
			util.Log.Fatalln(err)
		}
		util.Log.Println("[Daemon mode entered]")
		return
	}
	Guardian(true)
}
