package global

import (
	"flag"
)

var Flags struct {
	// 配置文件路径
	Path string

	// settings overwrite
	Interface string
	Debug     bool
	AutoAcid  bool
	Acid      string
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.yaml", "config path")

	flag.StringVar(&Flags.Interface, "interface", "", "specify the eth name")
	flag.BoolVar(&Flags.Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&Flags.AutoAcid, "auto-acid", false, "auto detect acid")
	flag.StringVar(&Flags.Acid, "acid", "", "specify acid value")

	flag.Parse()
}
