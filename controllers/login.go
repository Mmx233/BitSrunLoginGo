package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	BitSrun "github.com/Mmx233/BitSrunLoginGo/v1"
	"net"
)

// Login 登录逻辑
func Login(localAddr net.Addr) error {
	return BitSrun.Login(&BitSrun.LoginConf{
		Https: global.Config.Settings.Basic.Https,
		LoginInfo: BitSrun.LoginInfo{
			Form: &global.Config.Form,
			Meta: &global.Config.Meta,
		},
		Client: util.HttpTools(localAddr).Client,
	})
}
