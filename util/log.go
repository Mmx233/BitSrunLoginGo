package util

import (
	"fmt"
	"github.com/Mmx233/tool"
	"log"
	"os"
	"runtime"
	"time"
)

type loG struct {
	timeStamp string
	WriteFile bool
	Path      string
	Debug     bool
	OutPut    bool
}

var Log loG

func (c *loG) Init(debug, logFile, outPut bool, path string) {
	c.Debug = debug
	c.WriteFile = logFile
	c.OutPut = outPut
	c.Path = path
}

func (c *loG) time() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func (c *loG) WriteLog(name string, a ...interface{}) {
	if !(c.Debug && c.WriteFile) {
		return
	}
	var t string
	for _, v := range a {
		t += fmt.Sprint(v)
	}
	err := tool.File.Add(c.Path+name, c.time()+" "+t, 700)
	if err != nil {
		log.Println("Write log error: ", err)
	}
}

func (c *loG) genTimeStamp() {
	if c.timeStamp == "" {
		c.timeStamp = c.time()
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
