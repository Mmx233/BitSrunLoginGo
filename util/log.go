package util

import (
	"fmt"
	"github.com/Mmx233/tool"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

type loG struct {
	timeStamp string
	Demo      bool
	OutPut    bool
}

var Log loG

func (c *loG) WriteLog(name string, a ...interface{}) {
	if !c.Demo {
		return
	}
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
		err := tool.File.Add(name, fmt.Sprintf(time.Now().Format("2006/01/02 15:04:05 "))+t, 700)
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
	c.genTimeStamp()
	c.WriteLog("Login-"+c.timeStamp+".log", a...)
	if !c.OutPut {
		return
	}
	log.Println(a...)
}

func (c *loG) Fatalln(a ...interface{}) {
	c.genTimeStamp()
	c.WriteLog("LoginError-"+c.timeStamp+".log", a...)
	if !c.OutPut {
		return
	}
	log.Fatalln(a...)
}

func (c *loG) CatchRecover() {
	if e := recover(); e != nil {
		c.Println(e)
		var buf [4096]byte
		c.Println(string(buf[:runtime.Stack(buf[:], false)]))
		os.Exit(1)
	}
}
