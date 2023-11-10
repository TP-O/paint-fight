package config

import "github.com/spf13/viper"

type Service struct {
	Web string `mapstructure:"web"`
}

var _ configLoader = (*Service)(nil)

func (Service) loadDefault() {
	viper.SetDefault("service", map[string]interface{}{
		"web": "http://127.0.0.1:3000",
	})
}
