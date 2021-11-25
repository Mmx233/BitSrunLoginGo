package util

import "github.com/Mmx233/BitSrunLoginGo/models"

func GenerateLoginInfo(Form *srunModels.LoginForm, Meta *srunModels.LoginMeta) *srunModels.LoginInfo {
	return &srunModels.LoginInfo{
		UrlLoginPage:       "http://" + Form.Domain + "/srun_portal_success",
		UrlGetChallengeApi: "http://" + Form.Domain + "/cgi-bin/get_challenge",
		UrlLoginApi:        "http://" + Form.Domain + "/cgi-bin/srun_portal",
		UrlCheckApi:        "http://" + Form.Domain + "/cgi-bin/rad_user_info",
		Meta:               Meta,
		Form: &srunModels.LoginForm{
			UserName: Form.UserName + "@" + Form.UserType,
			PassWord: Form.PassWord,
		},
	}
}
