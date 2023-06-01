package global

import (
	"flag"
)

var Flags struct {
	//配置文件路径
	Path string

	Interface string
	Debug     bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.yaml", "config path")
	flag.StringVar(&Flags.Interface, "interface", "", "specify the eth name")
	flag.BoolVar(&Flags.Debug, "debug", false, "enable debug mode")
	flag.Parse()
}
