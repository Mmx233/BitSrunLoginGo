package controllers

import (
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"github.com/howeyc/fsnotify"
	"os"
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
		return false
	} else {
		return string(data) == a.Mark
	}
}

func (a *daemon) DaemonChan() bool {
	f, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	err = f.Watch(Daemon.Path)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-f.Event:
			if event.IsModify() && !a.CheckDaemon() {
				continue
			}
			os.Exit(0)
		case e := <-f.Error:
			panic(e)
		}
	}
}
