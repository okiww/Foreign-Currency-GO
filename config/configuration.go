package config

import (
	"log"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// New create new configuration object
func New() (*Configuration, error) {
	viper.SetConfigName("default")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.SetConfigName(".env")
	if err := viper.MergeInConfig(); err != nil {
		glog.Warningf("Failed to load custom configuration from .env file: %s", err)
	}

	cfg := new(Configuration)
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}
