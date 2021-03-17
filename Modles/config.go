package Modles

type Settings struct {
	QuitIfNetOk bool `json:"quit_if_net_ok"`
	DemoMode    bool `json:"demo_mode"`
}

type Config struct {
	From     LoginForm `json:"from"`
	Meta     LoginMeta `json:"meta"`
	Settings Settings  `json:"settings"`
}
