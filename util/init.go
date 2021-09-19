package util

import (
	"github.com/Mmx233/tool"
	"time"
)

func init() {
	//http工具设定
	tool.HTTP.Options.Timeout = 3 * time.Second
}
