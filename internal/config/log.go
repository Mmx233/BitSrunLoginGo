package config

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func initLog() {
	Logger = log.New()

	if Settings.Log.DebugLevel {
		Logger.SetLevel(log.DebugLevel)
	}

	if Settings.Log.WriteFile {
		//日志路径初始化与处理
		if !strings.HasSuffix(Settings.Log.FilePath, "/") {
			Settings.Log.FilePath += "/"
		}
		err := os.MkdirAll(Settings.Log.FilePath, os.ModePerm)
		if err != nil {
			Logger.Fatalln(err)
		}

		if Settings.Log.FileName == "" {
			Settings.Log.FileName = time.Now().Format("2006.01.02-15.04.05") + ".log"
		}

		f, err := os.OpenFile(Settings.Log.FilePath+Settings.Log.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			Logger.Fatalln(err)
		}

		//设置双重输出
		mw := io.MultiWriter(os.Stdout, f)
		Logger.SetOutput(mw)
	}

	//设置输出格式
	Logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		NoColors:        Settings.Log.WriteFile,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
