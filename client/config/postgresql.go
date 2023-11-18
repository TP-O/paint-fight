package config

import "github.com/spf13/viper"

type PostgreSQL struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	PoolSize int32  `mapstructure:"poolSize"`
	RootCA   string `mapstructure:"rootCA"`
}

var _ configLoader = (*PostgreSQL)(nil)

func (PostgreSQL) loadDefault() {
	viper.SetDefault("postgres", map[string]interface{}{
		"host":     "postgres",
		"port":     5432,
		"username": "dgame",
		"password": "dgame",
		"database": "dgame",
		"poolSize": 5,
	})
}
