package util

import (
	"Mmx/Global"
	"fmt"
	"log"
	"reflect"
	"time"
)

type loG struct {
	timeStamp string
}

var Log loG

func (*loG) WriteLog(name string, a ...interface{}) {
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
		err := File.Add(name, fmt.Sprintf(time.Now().Format("2006/01/02 15:04:05 "))+t)
		if err != nil {
			log.Println("Log error: ", err)
		}
	}
}

func (c *loG) genTimeStamp() {
	if c.timeStamp == "" {
		c.timeStamp = time.Now().Format("2006.01.02-15.04.05")
	}
}

func (c *loG) Println(a ...interface{}) {
	if !global.Status.Output {
		return
	}
	c.genTimeStamp()
	if global.Config.Settings.DemoMode {
		c.WriteLog("Login-"+c.timeStamp+".log", a...)
	}
	log.Println(a...)
}

func (c *loG) Fatalln(a ...interface{}) {
	if !global.Status.Output {
		return
	}
	c.genTimeStamp()
	c.WriteLog("LoginError-"+c.timeStamp+".log", a...)
	log.Fatalln(a...)
}
