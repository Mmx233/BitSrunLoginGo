package flags

import (
	"flag"
)

var (
	// Path 配置文件路径
	Path string

	Interface string
	Debug     bool
	AutoAcid  bool
	Acid      string
	AutoEnc   bool
)

func init() {
	flag.StringVar(&Path, "config", "Config.yaml", "config path")

	flag.StringVar(&Interface, "interface", "", "specify the eth name")
	flag.BoolVar(&Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&AutoAcid, "auto-acid", false, "auto detect acid")
	flag.StringVar(&Acid, "acid", "", "specify acid value")
	flag.BoolVar(&AutoEnc, "auto-enc", false, "auto detect enc")

	flag.Parse()
}
