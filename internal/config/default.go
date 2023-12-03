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
		N:           "200",
		Type:        "1",
		Acid:        "5",
		Enc:         "srun_bx1",
		OS:          "Windows 10",
		Name:        "windows",
		InfoPrefix:  "SRBX1",
		DoubleStack: false,
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
		CustomHeader: map[string]interface{}{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		},
	},
}
