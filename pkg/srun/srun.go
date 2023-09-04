package srun

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func (c Srun) LoginStatus() (online bool, ip string, err error) {
	res, err := c.api.GetUserInfo()
	if err != nil {
		return false, "", err
	}

	errRes, ok := res["error"]
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

	ip = ipInterface.(string)
	online = errRes.(string) == "ok"
	return
}

func (c Srun) DoLogin(clientIP string) error {
	log.Debugln("正在获取 Token")

	if c.LoginInfo.Form.UserType != "" {
		c.LoginInfo.Form.Username += "@" + c.LoginInfo.Form.UserType
	}

	res, err := c.api.GetChallenge(c.LoginInfo.Form.Username, clientIP)
	if err != nil {
		return err
	}
	token, ok := res["challenge"]
	if !ok {
		return ErrResultCannotFound
	}
	tokenStr := token.(string)
	log.Debugln("token: ", tokenStr)

	log.Debugln("发送登录请求")

	info, err := json.Marshal(map[string]string{
		"username": c.LoginInfo.Form.Username,
		"password": c.LoginInfo.Form.Password,
		"ip":       clientIP,
		"acid":     c.LoginInfo.Meta.Acid,
		"enc_ver":  c.LoginInfo.Meta.Enc,
	})
	if err != nil {
		return err
	}
	EncryptedInfo := "{SRBX1}" + Base64(XEncode(string(info), tokenStr))
	Md5Str := Md5(tokenStr)
	EncryptedMd5 := "{MD5}" + Md5Str
	EncryptedChkstr := Sha1(
		tokenStr + c.LoginInfo.Form.Username + tokenStr + Md5Str +
			tokenStr + c.LoginInfo.Meta.Acid + tokenStr + clientIP +
			tokenStr + c.LoginInfo.Meta.N + tokenStr + c.LoginInfo.Meta.Type +
			tokenStr + EncryptedInfo,
	)

	res, err = c.api.Login(
		c.LoginInfo.Form.Username,
		EncryptedMd5,
		c.LoginInfo.Meta.Acid,
		clientIP,
		EncryptedInfo,
		EncryptedChkstr,
		c.LoginInfo.Meta.N,
		c.LoginInfo.Meta.Type,
	)
	if err != nil {
		return err
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
