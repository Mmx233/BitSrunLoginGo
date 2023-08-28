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
	addresses, err := eth.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addresses {
		if strings.Contains(addr.String(), ".") {
			var ip *net.TCPAddr
			ip, err = net.ResolveTCPAddr("tcp", strings.Split(addr.String(), "/")[0]+":0")
			if err != nil {
				log.Warnln(eth.Name+" ip解析失败：", err)
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

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	reg, err := regexp.Compile(regexpStr)
	if err != nil {
		log.Fatalln("interfaces设置异常，无法解析: ", err)
	}
	for _, eth := range interfaces {
		if reg.Match([]byte(eth.Name)) {
			cEth, err := ConvertInterface(eth)
			if err != nil {
				log.Warnln(eth.Name+" 网卡地址获取失败: ", err)
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
