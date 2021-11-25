package util

import "github.com/Mmx233/BitSrunLoginGo/models"

func GenerateLoginInfo(Form *models.LoginForm, Meta *models.LoginMeta) *models.LoginInfo {
	return &models.LoginInfo{
		UrlLoginPage:       "http://" + Form.Domain + "/srun_portal_success",
		UrlGetChallengeApi: "http://" + Form.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + Form.Domain + "/cgi-bin/srun_portal",
		UrlCheckApi:        "http://" + Form.Domain + "/cgi-bin/rad_user_info",
		Meta:               Meta,
		Form: &models.LoginForm{
			UserName: Form.UserName + "@" + Form.UserType,
			PassWord: Form.PassWord,
		},
	}
}
