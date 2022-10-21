package global

import (
	"io"
	"os"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	if Config.Settings.Log.DebugLevel {
		log.SetLevel(log.DebugLevel)
	}

	if Config.Settings.Log.WriteFile {
		//日志路径初始化与处理
		if !strings.HasSuffix(Config.Settings.Log.FilePath, "/") {
			Config.Settings.Log.FilePath += "/"
		}
		e := os.MkdirAll(Config.Settings.Log.FilePath, os.ModePerm)
		if e != nil {
			log.Fatalln(e)
		}

		if Config.Settings.Log.FileName == "" {
			Config.Settings.Log.FileName = time.Now().Format("2006.01.02-15.04.05") + ".log"
		}

		f, e := os.OpenFile(Config.Settings.Log.FilePath+Config.Settings.Log.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if e != nil {
			log.Fatalln(e)
		}

		//设置双重输出
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
		log.SetFormatter(&nested.Formatter{
			HideKeys:        true,
			NoColors:        Config.Settings.Log.WriteFile,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}
