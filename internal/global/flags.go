package global

import (
	"flag"
)

var Flags struct {
	//配置文件路径
	Path string

	Interface string
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.yaml", "config path")
	flag.StringVar(&Flags.Interface, "interface", "", "specify the eth name")
	flag.Parse()
}
