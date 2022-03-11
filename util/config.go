package util

import (
	"github.com/Mmx233/BitSrunLoginGo/models"
	"github.com/Mmx233/BitSrunLoginGo/v1/transfer"
)

func GenerateLoginInfo(Form *srunTransfer.LoginForm, Meta *srunTransfer.LoginMeta) *srunModels.LoginInfo {
	return &srunModels.LoginInfo{
		Meta: Meta,
		Form: &srunTransfer.LoginForm{
			UserName: Form.UserName + "@" + Form.UserType,
			PassWord: Form.PassWord,
		},
	}
}
