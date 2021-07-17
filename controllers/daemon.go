package controllers

import (
	"autoLogin/util"
	"flag"
	"fmt"
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
	goDaemon := flag.Bool("daemon", false, "")
	flag.Parse()
	if *goDaemon {
		Guardian(false)
	}
}

func (a *daemon) MarkDaemon() error {
	return util.File.Write(a.Path, []byte(a.Mark))
}

func (a *daemon) CheckDaemon() bool {
	if data, err := util.File.Read(a.Path); err != nil {
		return true
	} else {
		return string(data) == a.Mark
	}
}
