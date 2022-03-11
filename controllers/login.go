package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	BitSrun "github.com/Mmx233/BitSrunLoginGo/v1"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
	"net"
)

// Login 登录逻辑
func Login(output bool, localAddr net.Addr) error {
	return BitSrun.Login(&srunTransfer.Login{
		Https:    global.Config.Settings.Basic.Https,
		Debug:    global.Config.Settings.Debug.Enable,
		WriteLog: global.Config.Settings.Debug.WriteLog,
		OutPut:   output,
		LoginInfo: srunTransfer.LoginInfo{
			Form: &global.Config.Form,
			Meta: &global.Config.Meta,
		},
		Transport: global.Transports(localAddr),
	})
}
