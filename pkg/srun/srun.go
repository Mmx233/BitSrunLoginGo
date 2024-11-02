package srun

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Conf struct {
	//调用 API 时直接访问 https URL
	Https bool
	//登录参数，不可缺省
	LoginInfo    LoginInfo
	Client       *http.Client
	CustomHeader map[string]interface{}

	Logger log.FieldLogger
}

func New(conf *Conf) *Srun {
	if conf.Logger == nil {
		conf.Logger = log.New()
	}

	srun := &Srun{
		LoginInfo: conf.LoginInfo,
		Logger:    conf.Logger,
	}
	srun.Api.Init(&ApiConfig{
		Https:        conf.Https,
		Domain:       conf.LoginInfo.Form.Domain,
		Client:       conf.Client,
		CustomHeader: conf.CustomHeader,
		Logger:       conf.Logger,
	})
	return srun
}

type Srun struct {
	//登录参数，不可缺省
	LoginInfo LoginInfo
	Api       Api
	Logger    log.FieldLogger
}

// LoginStatus err will not be nil only when both online and ip is found.
func (c Srun) LoginStatus() (online *bool, ip *string, err error) {
	res, err := c.Api.GetUserInfo()
	if err != nil {
		return nil, nil, err
	}

	errRes, ok := res["error"]
	if ok {
		isOnlineStr, ok := errRes.(string)
		if ok {
			isOnline := isOnlineStr == "ok"
			online = &isOnline
		}
	}

	ipInterface, ok := res["client_ip"]
	if !ok {
		ipInterface, ok = res["online_ip"]
	}
	if ipInterface != nil {
		ipStr, ok := ipInterface.(string)
		if ok {
			ip = &ipStr
		}
	}

	if online == nil || ip == nil {
		err = ErrResultCannotFound
	}
	return
}

func (c Srun) DoLogin(clientIP string) error {
	c.Logger.Debugln("正在获取 Token")

	if c.LoginInfo.Form.UserType != "" {
		c.LoginInfo.Form.Username += "@" + c.LoginInfo.Form.UserType
	}

	res, err := c.Api.GetChallenge(c.LoginInfo.Form.Username, clientIP)
	if err != nil {
		return err
	}
	token, ok := res["challenge"]
	if !ok {
		return ErrResultCannotFound
	}
	tokenStr := token.(string)
	c.Logger.Debugln("token: ", tokenStr)

	c.Logger.Debugln("发送登录请求")

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

	var infoPrefix string
	if c.LoginInfo.Meta.InfoPrefix != "" {
		infoPrefix = fmt.Sprintf("{%s}", c.LoginInfo.Meta.InfoPrefix)
	}
	EncryptedInfo := infoPrefix + Base64(XEncode(string(info), tokenStr))
	Md5Str, err := Md5(tokenStr, c.LoginInfo.Form.Password)
	if err != nil {
		return err
	}
	EncryptedMd5 := "{MD5}" + Md5Str
	EncryptedChkstr := Sha1(
		tokenStr + c.LoginInfo.Form.Username + tokenStr + Md5Str +
			tokenStr + c.LoginInfo.Meta.Acid + tokenStr + clientIP +
			tokenStr + c.LoginInfo.Meta.N + tokenStr + c.LoginInfo.Meta.Type +
			tokenStr + EncryptedInfo,
	)

	var doubleStack string
	if c.LoginInfo.Meta.DoubleStack {
		doubleStack = "1"
	} else {
		doubleStack = "0"
	}

	res, err = c.Api.Login(&LoginRequest{
		Username:    c.LoginInfo.Form.Username,
		Password:    EncryptedMd5,
		AcID:        c.LoginInfo.Meta.Acid,
		IP:          clientIP,
		Info:        EncryptedInfo,
		ChkSum:      EncryptedChkstr,
		N:           c.LoginInfo.Meta.N,
		Type:        c.LoginInfo.Meta.Type,
		OS:          c.LoginInfo.Meta.OS,
		Name:        c.LoginInfo.Meta.Name,
		DoubleStack: doubleStack,
	})
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
