package global

import (
	"os"
	"time"
)

func init() {
	initFlags()

	//配置文件初始化
	if readConfig() != nil {
		os.Exit(1)
	}

	//初始化常变量
	Timeout = time.Duration(Config.Settings.Basic.Timeout) * time.Second

	//初始化日志配置
	initLog()
}
