package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	srunModels "github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"os"
	"os/exec"
	"time"
)

// Guardian 守护模式逻辑
func Guardian(output bool) {
	util.Log.OutPut = output
	GuardianDuration := time.Duration(global.Config.Settings.Guardian.Duration) * time.Second

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
			if global.Config.Settings.Basic.Interfaces == "" { //单网卡
				if !util.Checker.NetOk(global.Transports(nil)) {
					util.Log.Println("Network down, trying to login")
					e := Login(output, true, nil)
					if e != nil {
						util.Log.Println("Error: ", e)
					}
				} else {
					if global.Config.Settings.Debug.Enable {
						util.Log.Println("Network ok")
					}
				}
			} else { //多网卡
				interfaces, e := util.GetInterfaceAddr()
				if e == nil {
					var down []srunModels.Eth
					for _, eth := range interfaces {
						if !util.Checker.NetOk(global.Transports(eth.Addr)) {
							util.Log.Println(eth.Name + " network down")
							down = append(down, eth)
						}
					}

					for _, eth := range down {
						util.Log.Println(eth.Name)
						e := Login(output, true, eth.Addr)
						if e != nil {
							util.Log.Println(eth.Name+" login error: ", e)
						}
					}
				}
			}

			c <- false
		}()
		<-c
		time.Sleep(GuardianDuration)
	}
}

// EnterGuardian 守护模式入口，控制是否进入daemon
func EnterGuardian() {
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
