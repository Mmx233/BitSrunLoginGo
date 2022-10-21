package util

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/v1"
)

func GenerateLoginInfo(Form *BitSrun.LoginForm, Meta *BitSrun.LoginMeta) *srunModels.LoginInfo {
	return &srunModels.LoginInfo{
		Meta: Meta,
		Form: &BitSrun.LoginForm{
			UserName: func() string {
				if Form.UserType == "" {
					return Form.UserName
				} else {
					return Form.UserName + "@" + Form.UserType
				}
			}(),
			PassWord: Form.PassWord,
		},
	}
}
