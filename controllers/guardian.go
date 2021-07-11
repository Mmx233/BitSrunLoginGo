package controllers

import (
	"Mmx/global"
	"Mmx/util"
	"time"
)

func Guardian() {
	for {
		time.Sleep(time.Duration(global.Config.Settings.Guardian) * time.Second)
		go func() {
			defer func() {
				_ = recover()
			}()
			if !util.Checker.NetOk() {
				_ = Login(false)
			}
		}()
	}
}
