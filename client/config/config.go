package config

import (
	"log"

	"github.com/spf13/viper"
)

type configLoader interface {
	loadDefault()
}

type config struct {
	App        `mapstructure:"app"`
	Secret     `mapstructure:"secret"`
	PostgreSQL `mapstructure:"postgresql"`
}

var cfg *config

// loadDefaultConfig loads the default config values.
func loadDefaultConfig(cfg *config) {
	cfg.App.loadDefault()
	cfg.Secret.loadDefault()
	cfg.PostgreSQL.loadDefault()
}

// Load loads config values from the given path and
// uses the default values if not set.
func Load(path string) *config {
	if cfg == nil {
		cfg = &config{}
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("yml")

		if err := viper.ReadInConfig(); err != nil {
			log.Println("Unable to load config:", err)
			log.Println("Use default config!")
		}

		loadDefaultConfig(cfg)
		if err := viper.Unmarshal(cfg); err != nil {
			log.Panic(err)
		}

		if len(cfg.Secret.Auth) < 32 {
			log.Panicf("The secret key is insecure: %s", cfg.Secret.Auth)
		}
	}

	return cfg
}
