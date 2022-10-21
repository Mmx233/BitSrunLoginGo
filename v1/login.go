package BitSrun

import (
	"encoding/json"

	"github.com/Mmx233/BitSrunLoginGo/util"
	srunTransfer "github.com/Mmx233/BitSrunLoginGo/v1/transfer"
	log "github.com/sirupsen/logrus"
)

func Login(c *srunTransfer.Login) error {
	G := util.GenerateLoginInfo(c.LoginInfo.Form, c.LoginInfo.Meta)
	api := SrunApi{
		BaseUrl: func() string {
			url := "http"
			if c.Https {
				url += "s"
			}
			return url + "://" + c.LoginInfo.Form.Domain + "/"
		}(),
		Client: c.Client,
	}

	var ok bool

	{
		log.Debugln("正在检查登录状态")

		res, e := api.GetUserInfo()
		if e != nil {
			return e
		}
		err := res["error"].(string)
		if err == "ok" {
			log.Debugln("已登录~")
			return nil
		}
		log.Infoln("检测到用户未登录，开始尝试登录...")

		{
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
		}
	}

	{
		log.Debugln("正在获取 Token")

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
		log.Debugln("token: ", G.Token)
	}

	{
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

		if G.LoginResult == "ok" {
			log.Infoln("已成功登录~")
		} else {
			log.Errorf("登陆失败: %s\n请开启日志 debug_level 获取更多信息", G.LoginResult)
			return nil
		}
	}

	return nil
}
