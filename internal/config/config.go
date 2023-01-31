package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type clickhouse struct {
	Address  string `json:"Address"`
	Database string `json:"Database"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type Config struct {
	Port       uint64     `json:"port"`
	Clickhouse clickhouse `json:"clickhouse"`
}

func NewConfig() *Config {
	var c Config
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.Unmarshal(&c)
	return &c
}
