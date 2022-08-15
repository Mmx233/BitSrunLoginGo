package global

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

func initLog() {
	if Config.Settings.Debug.Enable {
		log.SetLevel(log.DebugLevel)

		if Config.Settings.Debug.WriteLog {
			//日志路径初始化与处理
			if !strings.HasSuffix(Config.Settings.Debug.LogPath, "/") {
				Config.Settings.Debug.LogPath += "/"
			}
			e := os.MkdirAll(Config.Settings.Debug.LogPath, os.ModePerm)
			if e != nil {
				log.Fatalln(e)
			}

			f, e := os.OpenFile(Config.Settings.Debug.LogPath+time.Now().Format("2006.01.02-15.04.05")+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 700)
			if e != nil {
				log.Fatalln(e)
			}

			//设置双重输出
			mw := io.MultiWriter(os.Stdout, f)
			log.SetOutput(mw)
		}
	}
}
