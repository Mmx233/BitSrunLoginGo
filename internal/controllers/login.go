package controllers

import (
	global2 "github.com/Mmx233/BitSrunLoginGo/internal/global"
	dns2 "github.com/Mmx233/BitSrunLoginGo/internal/pkg/dns"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

// Login 登录逻辑
func Login(localAddr net.Addr, debugOutput bool) error {
	// 登录状态检查

	httpClient := tools.HttpPackSelect(localAddr).Client
	conf := &srun.Conf{
		Https: global2.Config.Settings.Basic.Https,
		LoginInfo: srun.LoginInfo{
			Form: &global2.Config.Form,
			Meta: &global2.Config.Meta,
		},
		Client: httpClient,
	}

	var output func(args ...interface{})
	if debugOutput {
		output = log.Debugln
	} else {
		output = log.Infoln
	}

	output("正在获取登录状态")

	online, ip, e := srun.LoginStatus(conf)
	if e != nil {
		return e
	}

	if localAddr != nil && global2.Config.Settings.Basic.UseDhcpIP {
		ip = localAddr.(*net.TCPAddr).IP.String()
	} else if global2.Flags.ClientIP != "" {
		ip = global2.Flags.ClientIP
	}

	log.Debugln("认证客户端 ip: ", ip)

	// 登录执行

	if online {
		output("已登录~")

		if global2.Config.Settings.DDNS.Enable && global2.Config.Settings.Guardian.Enable && ipLast != ip {
			if ddns(ip, httpClient) == nil {
				ipLast = ip
			}
		}

		return nil
	} else {
		log.Infoln("检测到用户未登录，开始尝试登录...")

		if e = srun.DoLogin(ip, conf); e != nil {
			return e
		}

		log.Infoln("登录成功~")

		if global2.Config.Settings.DDNS.Enable {
			_ = ddns(ip, httpClient)
		}
	}

	return nil
}

var ipLast string

func ddns(ip string, httpClient *http.Client) error {
	return dns2.Run(&dns2.Config{
		Provider: global2.Config.Settings.DDNS.Provider,
		IP:       ip,
		Domain:   global2.Config.Settings.DDNS.Domain,
		TTL:      global2.Config.Settings.DDNS.TTL,
		Conf:     global2.Config.Settings.DDNS.Config,
		Http:     httpClient,
	})
}
