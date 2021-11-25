package global

import "flag"

var Flags struct {
	Path   string
	Daemon bool
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "", "config path")
	flag.BoolVar(&Flags.Daemon, "daemon", false, "")
	flag.Parse()
}
