package config

import (
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
)

var defaultConfig = ConfFromFile{
	Form: srun.LoginForm{
		Domain:   "www.msftconnecttest.com",
		UserType: "cmcc",
	},
	Meta: srun.LoginMeta{
		N:    "200",
		Type: "1",
		Acid: "5",
		Enc:  "srun_bx1",
	},
	Settings: SettingsConf{
		Basic: BasicConf{
			Timeout: 5,
		},
		Guardian: GuardianConf{
			Duration: 300,
		},
		Log: LogConf{
			FilePath: "./",
		},
		DDNS: DdnsConf{
			Enable:   false,
			TTL:      600,
			Domain:   "www.example.com",
			Provider: "cloudflare",
			Config: map[string]interface{}{
				"zone":  "",
				"token": "",
			},
		},
	},
}
