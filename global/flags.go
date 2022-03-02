package global

import (
	"flag"
)

var Flags struct {
	//配置文件路径
	Path string
	//daemon模式内置标记
	RunningDaemon bool
	//强制daemon
	Daemon bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.yaml", "config path")
	flag.BoolVar(&Flags.RunningDaemon, "running-daemon", false, "")
	flag.BoolVar(&Flags.Daemon, "daemon", false, "")
	flag.Parse()
}
