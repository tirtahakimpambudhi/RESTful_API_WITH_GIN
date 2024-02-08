package test

import (
	"github.com/spf13/viper"
	"go_gin/internal/config"
	"testing"
)

type server struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
}
type database struct {
	Host         string   `mapstructure:"host"`
	Port         string   `mapstructure:"port"`
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
	DatabaseName string   `mapstructure:"database_name"`
	TableName    []string `mapstructure:"table_name"`
}

type configt struct {
	Server   *server   `mapstructure:"server"`
	Database *database `mapstructure:"database"`
}

func TestConfigServer(t *testing.T) {
	var cfg configt
	viper.SetConfigName("conf") // name of config file (without extension)
	viper.SetConfigType("toml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		t.Error(err.Error())
	}
	viper.Unmarshal(&cfg)
	t.Log(cfg.Database)
}

func TestConfigs(t *testing.T) {
	t.Log(config.Other)
	t.Log(config.Database)
	t.Log(config.Server)
	t.Log([]byte(config.JWT.SecretKey))
}
