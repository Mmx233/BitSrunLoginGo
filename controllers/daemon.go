package controllers

import (
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"time"
)

type daemon struct {
	Mark string
	Path string
}

var Daemon = daemon{
	Mark: fmt.Sprint(time.Now().UnixNano()),
	Path: ".autoLoginDaemon",
}

func init() {
	if global.Flags.Daemon {
		Guardian(false)
	}
}

func (a *daemon) MarkDaemon() error {
	return tool.File.Write(a.Path, []byte(a.Mark))
}

func (a *daemon) CheckDaemon() bool {
	if data, err := tool.File.Read(a.Path); err != nil {
		return true
	} else {
		return string(data) == a.Mark
	}
}
