package util

import (
	"fmt"
	"github.com/Mmx233/tool"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type loG struct {
	timeStamp string
	WriteFile bool
	Path      string
	OutPut    bool
	DebugMode bool
}

var Log loG

func (c *loG) Init(debug, logFile, outPut bool, path string) error {
	c.DebugMode = debug
	c.WriteFile = logFile
	c.OutPut = outPut
	c.timeStamp = time.Now().Format("2006.01.02-15.04.05")

	//日志路径初始化与处理
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	c.Path = path
	return os.MkdirAll(path, os.ModePerm)
}

func (c *loG) time() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func (c *loG) WriteLog(name string, a ...interface{}) {
	err := tool.File.Add(c.Path+name, c.time()+" "+fmt.Sprint(a...), 700)
	if err != nil && c.OutPut {
		log.Println(err)
	}
}

func (c *loG) print(fatal bool, a ...interface{}) {
	if c.DebugMode && c.WriteFile {
		c.WriteLog("Login-"+c.timeStamp+".log", a...)
	}
	if c.OutPut {
		if fatal {
			if c.DebugMode {
				log.Panicln(a...)
			} else {
				log.Fatalln(a...)
			}
		} else {
			log.Println(a...)
		}
	}
}

func (c *loG) Debug(a ...interface{}) {
	if c.DebugMode {
		c.print(false, append([]interface{}{"[DEBUG] "}, a...)...)
	}
}

func (c *loG) Info(a ...interface{}) {
	c.print(false, append([]interface{}{"[INFO] "}, a...)...)
}

func (c *loG) Warn(a ...interface{}) {
	c.print(false, append([]interface{}{"[WARN] "}, a...)...)
}

func (c *loG) Fatal(a ...interface{}) {
	c.print(true, append([]interface{}{"[FATAL] "}, a...)...)
}

func (c *loG) CatchRecover() {
	if e := recover(); e != nil {
		var buf [4096]byte
		c.Fatal(e, "\n", string(buf[:runtime.Stack(buf[:], false)]))
	}
}
