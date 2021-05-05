package Util

import (
	"Mmx/Global"
	"fmt"
	"log"
	"reflect"
)

type loG struct{}

var Log loG

func (loG) WriteLog(name string, a ...interface{}) {
	for _, v := range a {
		var t string
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			t = v.(string)
		case reflect.Interface:
			t = v.(error).Error()
		default:
			t = fmt.Sprint(v)
		}
		_ = File.Add(name, t)
	}
}

func (c loG) Println(a ...interface{}) {
	if Global.Config.Settings.DemoMode {
		c.WriteLog("Login.loG", a...)
	}
	log.Println(a...)
}

func (c loG) Fatalln(a ...interface{}) {
	c.WriteLog("LoginError.loG", a...)
	log.Fatalln(a...)
}
