package config

import "github.com/spf13/viper"

const (
	DevEnv  = "development"
	ProdEnv = "production"
)

type App struct {
	Debug bool   `mapstructure:"debug"`
	Env   string `mapstructure:"env"`
	Port  uint16 `mapstructure:"port"`
}

var _ configLoader = (*App)(nil)

func (App) loadDefault() {
	viper.SetDefault("app", map[string]interface{}{
		"debug": true,
		"env":   DevEnv,
		"port":  8080,
	})
}
