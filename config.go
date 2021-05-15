package main

import (
	"github.com/spf13/viper"
	"log"
)

func SetupConfig() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}
