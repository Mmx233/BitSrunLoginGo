package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	"github.com/Mmx233/BitSrunLoginGo/util"
	BitSrun "github.com/Mmx233/BitSrunLoginGo/v1"
	log "github.com/sirupsen/logrus"
	"net"
)

// Login 登录逻辑
func Login(localAddr net.Addr, debugOutput bool) error {
	conf := &BitSrun.Conf{
		Https: global.Config.Settings.Basic.Https,
		LoginInfo: BitSrun.LoginInfo{
			Form: &global.Config.Form,
			Meta: &global.Config.Meta,
		},
		Client: util.HttpTools(localAddr).Client,
	}

	var output func(args ...interface{})
	if debugOutput {
		output = log.Debugln
	} else {
		output = log.Infoln
	}

	output("正在获取登录状态")

	online, ip, e := BitSrun.LoginStatus(conf)
	if e != nil {
		return e
	}

	log.Debugln("认证客户端 ip: ", ip)

	if online {
		output("已登录~")

		return nil
	} else {
		output("检测到用户未登录，开始尝试登录...")

		return BitSrun.DoLogin(ip, conf)
	}
}
