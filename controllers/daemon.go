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

// Daemon 后台模式控制包
var Daemon = daemon{
	Mark: fmt.Sprint(time.Now().UnixNano()),
	Path: global.Config.Settings.Daemon.Path,
}

// MarkDaemon 写入后台标记文件
func (a *daemon) MarkDaemon() error {
	return tool.File.Write(a.Path, []byte(a.Mark))
}

// CheckDaemon 检查后台标记文件
func (a *daemon) CheckDaemon() bool {
	if data, err := tool.File.Read(a.Path); err != nil {
		return false
	} else {
		return string(data) == a.Mark
	}
}

// DaemonChan 后台标记文件监听
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
			if event.IsModify() && a.CheckDaemon() {
				continue
			}
			os.Exit(0)
		case e := <-f.Error:
			panic(e)
		}
	}
}
