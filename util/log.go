package util

import (
	"fmt"
	"github.com/Mmx233/tool"
	"log"
	"os"
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
	if c.DebugMode && c.WriteFile {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		c.Path = path
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
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

func (c *loG) print(name string, fatal bool, a ...interface{}) {
	a = append([]interface{}{"[" + name + "] "}, a...)
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
		c.print("DEBUG", false, a...)
	}
}

func (c *loG) Info(a ...interface{}) {
	c.print("INFO", false, a...)
}

func (c *loG) Warn(a ...interface{}) {
	c.print("WARN", false, a...)
}

func (c *loG) Fatal(a ...interface{}) {
	c.print("FATAL", true, a...)
}
