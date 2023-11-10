package config

import "github.com/spf13/viper"

type Secret struct {
	Auth string `mapstructure:"auth"`
}

var _ configLoader = (*Secret)(nil)

func (Secret) loadDefault() {
	viper.SetDefault("secret", map[string]interface{}{})
}
