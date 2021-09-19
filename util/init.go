package util

import (
	"autoLogin/global"
	"github.com/Mmx233/tool"
	"time"
)

func init() {
	//tool设定
	tool.HTTP.Options.Timeout = time.Duration(global.Config.Settings.Timeout) * time.Second
	tool.File.Options.ForceRoot = true
}
