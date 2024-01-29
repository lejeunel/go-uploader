package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	dbpath string `mapstructure:"dbpath"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("dbpath", "~/.cache/uploads.db")

	viper.SetEnvPrefix("uploader")
	// viper.BindEnv("dbpath")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// fmt.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
