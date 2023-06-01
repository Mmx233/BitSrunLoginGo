package tools

import (
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strings"
)

type Eth struct {
	Name string
	Addr net.Addr
}

// ConvertInterface 当没有 ipv4 地址时 eth 可能为 nil
func ConvertInterface(eth net.Interface) (*Eth, error) {
	addresses, e := eth.Addrs()
	if e != nil {
		return nil, e
	}
	for _, addr := range addresses {
		if strings.Contains(addr.String(), ".") {
			var ip *net.TCPAddr
			ip, e = net.ResolveTCPAddr("tcp", strings.Split(addr.String(), "/")[0]+":0")
			if e != nil {
				log.Warnln(eth.Name+" ip解析失败：", e)
				continue
			}
			return &Eth{
				Name: eth.Name,
				Addr: ip,
			}, nil
		}
	}
	return nil, nil
}

func GetInterfaceAddr(regexpStr string) ([]Eth, error) {
	var result []Eth

	interfaces, e := net.Interfaces()
	if e != nil {
		return nil, e
	}
	reg, e := regexp.Compile(regexpStr)
	if e != nil {
		log.Fatalln("interfaces设置异常，无法解析: ", e)
	}
	for _, eth := range interfaces {
		if reg.Match([]byte(eth.Name)) {
			cEth, e := ConvertInterface(eth)
			if e != nil {
				log.Warnln(eth.Name+" 网卡地址获取失败: ", e)
				continue
			}

			if cEth != nil {
				result = append(result, *cEth)
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
