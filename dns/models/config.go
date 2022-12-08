package models

type BasicConfig struct {
	Domain string                 `mapstructure:"domain"`
	Other  map[string]interface{} `mapstructure:",remain"`
}
