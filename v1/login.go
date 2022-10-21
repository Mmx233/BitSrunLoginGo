package BitSrun

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
)

func Login(c *Conf) error {
	c.initApi()
	G := GenerateLoginInfo(c.LoginInfo.Form, c.LoginInfo.Meta)

	log.Debugln("正在检查登录状态")

	res, e := c.api.GetUserInfo()
	if e != nil {
		return e
	}
	err, ok := res["error"]
	if err == "ok" {
		log.Debugln("已登录~")
		return nil
	}
	log.Infoln("检测到用户未登录，开始尝试登录...")

	log.Debugln("正在获取客户端 IP")

	var ip interface{}
	ip, ok = res["client_ip"]
	if !ok {
		ip, ok = res["online_ip"]
		if !ok {
			return ErrResultCannotFound
		}
	}
	G.Ip = ip.(string)
	log.Debugln("ip: ", G.Ip)

	log.Debugln("正在获取 Token")

	res, e = c.api.GetChallenge(G.Form.UserName, G.Ip)
	if e != nil {
		return e
	}
	var token interface{}
	token, ok = res["challenge"]
	if !ok {
		return ErrResultCannotFound
	}
	G.Token = token.(string)
	log.Debugln("token: ", G.Token)

	log.Debugln("发送登录请求")

	info, e := json.Marshal(map[string]string{
		"username": G.Form.UserName,
		"password": G.Form.PassWord,
		"ip":       G.Ip,
		"acid":     G.Meta.Acid,
		"enc_ver":  G.Meta.Enc,
	})
	if e != nil {
		return e
	}
	G.EncryptedInfo = "{SRBX1}" + Base64(XEncode(string(info), G.Token))
	G.Md5 = Md5(G.Token)
	G.EncryptedMd5 = "{MD5}" + G.Md5

	var chkstr = G.Token + G.Form.UserName + G.Token + G.Md5
	chkstr += G.Token + G.Meta.Acid + G.Token + G.Ip
	chkstr += G.Token + G.Meta.N + G.Token + G.Meta.Type
	chkstr += G.Token + G.EncryptedInfo
	G.EncryptedChkstr = Sha1(chkstr)

	res, e = c.api.Login(
		G.Form.UserName,
		G.EncryptedMd5,
		G.Meta.Acid,
		G.Ip,
		G.EncryptedInfo,
		G.EncryptedChkstr,
		G.Meta.N,
		G.Meta.Type,
	)
	if e != nil {
		return e
	}
	var result interface{}
	result, ok = res["error"]
	if !ok {
		return ErrResultCannotFound
	}
	G.LoginResult = result.(string)

	if G.LoginResult == "ok" {
		return errors.New(G.LoginResult)
	}

	return nil
}
