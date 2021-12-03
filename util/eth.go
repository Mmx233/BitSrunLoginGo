package util

import (
	"github.com/Mmx233/BitSrunLoginGo/global"
	srunModels "github.com/Mmx233/BitSrunLoginGo/models"
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
	reg, e := regexp.Compile(global.Config.Settings.Interfaces)
	if e != nil {
		Log.Println("interfaces设置异常，无法解析")
		return nil, e
	}
	for _, eth := range interfaces {
		if reg.Match([]byte(eth.Name)) {
			addrs, e := eth.Addrs()
			if e != nil {
				Log.Println(eth.Name + " 地址获取失败")
			}
			for _, ip := range addrs {
				if strings.Contains(ip.String(), ".") {
					result = append(result, srunModels.Eth{
						Name: eth.Name,
						Addr: func() *net.TCPAddr {
							n, _ := net.ResolveTCPAddr("tcp", ip.String())
							return n
						}(),
					})
					break
				}
			}
		}
	}

	return result, nil
}
