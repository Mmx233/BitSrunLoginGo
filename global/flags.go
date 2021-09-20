package global

import "flag"

var Flags struct {
	Path string
}

func initFlags() {
	flag.StringVar(&Flags.Path, "config", "", "config path")
	flag.Parse()
}
