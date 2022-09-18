package util

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	srunModels "github.com/Mmx233/BitSrunLoginGo/models"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strings"
)

func GetInterfaceAddr() ([]srunModels.Eth, error) {
	var result []srunModels.Eth

	interfaces, e := net.Interfaces()
	if e != nil {
		return nil, e
	}
	reg, e := regexp.Compile(global.Config.Settings.Basic.Interfaces)
	if e != nil {
		log.Fatalln("interfaces设置异常，无法解析: ", e)
	}
	for _, eth := range interfaces {
		if reg.Match([]byte(eth.Name)) {
			addrs, e := eth.Addrs()
			if e != nil {
				log.Warnln(eth.Name+" 网卡地址获取失败: ", e)
				continue
			}
			for _, addr := range addrs {
				if strings.Contains(addr.String(), ".") {
					var ip *net.TCPAddr
					ip, e = net.ResolveTCPAddr("tcp", strings.Split(addr.String(), "/")[0]+":0")
					if e != nil {
						log.Warnln(eth.Name+" ip解析失败：", e)
						continue
					}
					result = append(result, srunModels.Eth{
						Name: eth.Name,
						Addr: ip,
					})
					break
				}
			}
		} else {
			log.Debugf("网卡 %s 不匹配", eth.Name)
		}
	}

	log.Debugln("有效匹配网卡：", result)

	if len(result) == 0 {
		log.Warnln("没有扫描到有效匹配网卡")
	}

	return result, nil
}
