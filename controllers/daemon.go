package controllers

import (
	"flag"
)

func init() {
	goDaemon := flag.Bool("daemon", false, "")
	flag.Parse()
	if *goDaemon {
		Guardian()
	}
}
