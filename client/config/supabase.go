package config

import "github.com/spf13/viper"

type Supabase struct {
	ID      string `mapstructure:"id"`
	AnonKey string `mapstructure:"anonkey"`
}

var _ configLoader = (*Supabase)(nil)

func (Supabase) loadDefault() {
	viper.SetDefault("supabase", map[string]interface{}{})
}
