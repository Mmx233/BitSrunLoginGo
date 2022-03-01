package util

import (
	"fmt"
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/tool"
	"log"
	"os"
	"runtime"
	"time"
)

type loG struct {
	timeStamp string
	WriteFile bool
	Debug     bool
	OutPut    bool
}

var Log loG

func (c *loG) WriteLog(name string, a ...interface{}) {
	if !(c.Debug && c.WriteFile) {
		return
	}
	var t string
	for i, v := range a {
		t += fmt.Sprint(v)
		if i != len(a)-1 {
			t += " "
		}
	}
	err := tool.File.Add(global.Config.Settings.Debug.Path+name, fmt.Sprintf(time.Now().Format("2006/01/02 15:04:05 "))+t, 700)
	if err != nil {
		log.Println("Write log error: ", err)
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
	if c.OutPut {
		log.Println(a...)
	}
}

func (c *loG) Fatalln(a ...interface{}) {
	c.genTimeStamp()
	c.WriteLog("LoginError-"+c.timeStamp+".log", a...)
	if c.OutPut {
		log.Fatalln(a...)
	}
}

func (c *loG) CatchRecover() {
	if e := recover(); e != nil {
		c.Println(e)
		var buf [4096]byte
		c.Println(string(buf[:runtime.Stack(buf[:], false)]))
		os.Exit(1)
	}
}
