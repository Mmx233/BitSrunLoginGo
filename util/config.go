package util

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
)

func GenerateLoginInfo(https bool, Form *srunTransfer.LoginForm, Meta *srunTransfer.LoginMeta) *srunModels.LoginInfo {
	portal := "http"
	if https {
		portal += "s"
	}
	portal += "://"
	return &srunModels.LoginInfo{
		UrlLoginPage:       portal + Form.Domain + "/srun_portal_success",
		UrlGetChallengeApi: portal + Form.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        portal + Form.Domain + "/cgi-bin/srun_portal",
		UrlCheckApi:        portal + Form.Domain + "/cgi-bin/rad_user_info",
		Meta:               Meta,
		Form: &srunTransfer.LoginForm{
			UserName: Form.UserName + "@" + Form.UserType,
			PassWord: Form.PassWord,
		},
	}
}
