package global

import (
	"flag"
)

var Flags struct {
	Path          string
	RunningDaemon bool
	Daemon        bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "Config.json", "config path")
	flag.BoolVar(&Flags.RunningDaemon, "running-daemon", false, "")
	flag.BoolVar(&Flags.Daemon, "daemon", false, "")
	flag.Parse()
}
