package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type SRV struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
}

type DB struct {
	Host         string   `mapstructure:"host"`
	Port         string   `mapstructure:"port"`
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
	DatabaseName string   `mapstructure:"database_name"`
	TableName    []string `mapstructure:"table_name"`
}

type OTH struct {
	SecretKey   string `mapstructure:"secret_key"`
	SaltLevel   int    `mapstructure:"salt_level"`
	BatchSize   int    `mapstructure:"batch_size"`
	LimitInsert int    `mapstructure:"limit_insert"`
	Limit       int    `mapstructure:"limit"`
}

type CFG struct {
	Server   *SRV   `mapstructure:"server"`
	Database *DB    `mapstructure:"database"`
	Other    *OTH   `mapstructure:"other"`
	JWT      *TOKEN `mapstructure:"jwt"`
}

type TOKEN struct {
	AppName   string `mapstructure:"app_name"`
	Exp       int    `mapstructure:"exp"`
	SecretKey string `mapstructure:"secret_key"`
}

var (
	Database *DB
	Server   *SRV
	Other    *OTH
	JWT      *TOKEN
)

func init() {
	var config CFG
	workdir, _ := os.Getwd()
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Join(workdir, "config"))    //Untuk Direktori Utama
	viper.AddConfigPath(filepath.Join(workdir, "../config")) //Untuk Direktori Uji

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	viper.Unmarshal(&config)
	Database = config.Database
	Server = config.Server
	Other = config.Other
	JWT = config.JWT
}
