package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	BitSrun "github.com/Mmx233/BitSrunLoginGo/v1"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
	"net"
)

// Login 登录逻辑
func Login(localAddr net.Addr) error {
	return BitSrun.Login(&srunTransfer.Login{
		Https: global.Config.Settings.Basic.Https,
		LoginInfo: srunTransfer.LoginInfo{
			Form: &global.Config.Form,
			Meta: &global.Config.Meta,
		},
		Client: util.HttpTools(localAddr).Client,
	})
}
