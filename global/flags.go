package global

import (
	"flag"
)

var Flags struct {
	Path   string
	Daemon bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.json", "config path")
	flag.BoolVar(&Flags.Daemon, "daemon", false, "")
	flag.Parse()
}
