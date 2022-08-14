package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

// Guardian 守护模式逻辑
func Guardian() {
	GuardianDuration := time.Duration(global.Config.Settings.Guardian.Duration) * time.Second

	if global.Config.Settings.Daemon.Enable {
		go Daemon.DaemonChan()

		if e := Daemon.MarkDaemon(); e != nil {
			log.Warnln("写入daemon标记文件失败: ", e)
		}
	}

	var c = make(chan bool)
	for {
		go func() {
			defer func() {
				_ = recover()
			}()
			if global.Config.Settings.Basic.Interfaces == "" { //单网卡
				e := Login(nil)
				if e != nil {
					log.Errorln("登陆失败: ", e)
				}
			} else { //多网卡
				interfaces, e := util.GetInterfaceAddr()
				if e == nil {
					for _, eth := range interfaces {
						log.Infoln(eth.Name)
						e = Login(eth.Addr)
						if e != nil {
							log.Errorln("网口 ", eth.Name+" 登录失败: ", e)
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
	log.Infoln("[Guardian mode]")
	if global.Config.Settings.Daemon.Enable || global.Flags.Daemon {
		if err := exec.Command(os.Args[0], append(os.Args[1:], "--running-daemon")...).Start(); err != nil {
			log.Fatalln("启动守护失败: ", err)
		}
		log.Infoln("[Daemon mode entered]")
		return
	}
	Guardian()
}
