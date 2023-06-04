package srun

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Conf struct {
	//调用 API 时直接访问 https URL
	Https bool
	//登录参数，不可缺省
	LoginInfo LoginInfo
	Client    *http.Client
}

func New(conf *Conf) *Srun {
	srun := &Srun{
		LoginInfo: conf.LoginInfo,
	}
	srun.api.Init(conf.Https, conf.LoginInfo.Form.Domain, conf.Client)
	return srun
}

type Srun struct {
	//登录参数，不可缺省
	LoginInfo LoginInfo
	api       Api
}

func (c Srun) LoginStatus() (online bool, ip string, e error) {
	res, e := c.api.GetUserInfo()
	if e != nil {
		return false, "", e
	}

	err, ok := res["error"]
	if !ok {
		return false, "", ErrResultCannotFound
	}

	ipInterface, ok := res["client_ip"]
	if !ok {
		ipInterface, ok = res["online_ip"]
		if !ok {
			return false, "", ErrResultCannotFound
		}
	}

	// 如果深澜分配的 ip 不是内网 ip，说明已经在线且拥有固定 ip
	ip = ipInterface.(string)

	inet := strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.")

	online = err.(string) == "ok" || !inet

	return
}

func (c Srun) DoLogin(clientIP string) error {
	log.Debugln("正在获取 Token")

	if c.LoginInfo.Form.UserType != "" {
		c.LoginInfo.Form.UserName += "@" + c.LoginInfo.Form.UserType
	}

	res, e := c.api.GetChallenge(c.LoginInfo.Form.UserName, clientIP)
	if e != nil {
		return e
	}
	token, ok := res["challenge"]
	if !ok {
		return ErrResultCannotFound
	}
	tokenStr := token.(string)
	log.Debugln("token: ", tokenStr)

	log.Debugln("发送登录请求")

	info, e := json.Marshal(map[string]string{
		"username": c.LoginInfo.Form.UserName,
		"password": c.LoginInfo.Form.PassWord,
		"ip":       clientIP,
		"acid":     c.LoginInfo.Meta.Acid,
		"enc_ver":  c.LoginInfo.Meta.Enc,
	})
	if e != nil {
		return e
	}
	EncryptedInfo := "{SRBX1}" + Base64(XEncode(string(info), tokenStr))
	Md5Str := Md5(tokenStr)
	EncryptedMd5 := "{MD5}" + Md5Str
	EncryptedChkstr := Sha1(
		tokenStr + c.LoginInfo.Form.UserName + tokenStr + Md5Str +
			tokenStr + c.LoginInfo.Meta.Acid + tokenStr + clientIP +
			tokenStr + c.LoginInfo.Meta.N + tokenStr + c.LoginInfo.Meta.Type +
			tokenStr + EncryptedInfo,
	)

	res, e = c.api.Login(
		c.LoginInfo.Form.UserName,
		EncryptedMd5,
		c.LoginInfo.Meta.Acid,
		clientIP,
		EncryptedInfo,
		EncryptedChkstr,
		c.LoginInfo.Meta.N,
		c.LoginInfo.Meta.Type,
	)
	if e != nil {
		return e
	}
	var result interface{}
	result, ok = res["error"]
	if !ok {
		return ErrResultCannotFound
	}
	LoginResult := result.(string)

	if LoginResult != "ok" {
		return errors.New(LoginResult)
	}

	return nil
}

func (c Srun) DetectAcid() (string, error) {
	return c.api.DetectAcid()
}
