package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

	viper.ReadInConfig()

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
