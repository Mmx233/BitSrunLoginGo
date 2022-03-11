package BitSrun

import (
	"encoding/json"
	"errors"
	"github.com/Mmx233/BitSrunLoginGo/util"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
)

func Login(c *srunTransfer.Login) error {
	util.Log.DebugMode = c.Debug
	util.Log.WriteFile = c.WriteLog
	util.Log.OutPut = c.OutPut

	G := util.GenerateLoginInfo(c.LoginInfo.Form, c.LoginInfo.Meta)
	api := SrunApi{
		BaseUrl: func() string {
			url := "http"
			if c.Https {
				url += "s"
			}
			return url + "://" + c.LoginInfo.Form.Domain + "/"
		}(),
		Transport: c.Transport,
	}

	var ok bool

	{
		util.Log.Info("Step.0: 正在检查状态")
		res, e := api.GetUserInfo()
		if e != nil {
			return e
		}
		err := res["error"].(string)
		if err == "ok" {
			util.Log.Info("--已登录--")
			return nil
		}

		util.Log.Info("Step.1: 正在获取客户端ip")
		var ip interface{}
		ip, ok = res["client_ip"]
		if !ok {
			ip, ok = res["online_ip"]
			if !ok {
				return ErrResultCannotFound
			}
		}
		G.Ip = ip.(string)
		util.Log.Debug("ip: ", G.Ip)
	}

	util.Log.Info("Step.2: 正在获取token")
	{
		res, e := api.GetChallenge(G.Form.UserName, G.Ip)
		if e != nil {
			return e
		}
		var token interface{}
		token, ok = res["challenge"]
		if !ok {
			return ErrResultCannotFound
		}
		G.Token = token.(string)
		util.Log.Debug("token: ", G.Token)
	}

	util.Log.Info("Step.3: 执行登录…")
	{
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
		G.EncryptedInfo = "{SRBX1}" + util.Base64(util.XEncode(string(info), G.Token))
		G.Md5 = util.Md5(G.Token)
		G.EncryptedMd5 = "{MD5}" + G.Md5

		var chkstr = G.Token + G.Form.UserName + G.Token + G.Md5
		chkstr += G.Token + G.Meta.Acid + G.Token + G.Ip
		chkstr += G.Token + G.Meta.N + G.Token + G.Meta.Type
		chkstr += G.Token + G.EncryptedInfo
		G.EncryptedChkstr = util.Sha1(chkstr)

		res, e := api.Login(
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

		util.Log.Info("登录结果: " + G.LoginResult)
		if G.LoginResult != "ok" {
			return errors.New(G.LoginResult)
		}
	}

	return nil
}
