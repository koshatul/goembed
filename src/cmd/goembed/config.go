package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func configDefaults() {
}

func configInit() {
	viper.SetConfigName("afero-static")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./artifacts")
	viper.AddConfigPath("./test")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/etc/afero-static")
	viper.AddConfigPath("/usr/local/afero-static/etc")
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath(".")

	configDefaults()

	viper.ReadInConfig()

	configFormatting()

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func configFormatting() {
}

func zapConfig() zap.Config {
	var cfg zap.Config
	if viper.GetBool("debug") {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	return cfg
}
