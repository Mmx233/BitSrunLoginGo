package global

import (
	"flag"
)

var Flags struct {
	Path          string
	RunningDaemon bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.json", "config path")
	flag.BoolVar(&Flags.RunningDaemon, "running-daemon", false, "")
	flag.Parse()
}
