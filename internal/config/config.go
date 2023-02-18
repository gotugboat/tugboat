package config

import (
	"io/fs"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) error {
	validFileNames := []string{"tugboat", ".tugboat"}

	for _, file := range validFileNames {
		viper.SetConfigName(file)
		viper.SetConfigType("yaml")

		validPaths := []string{".", "./ci", "./.ci"}
		for _, path := range validPaths {
			viper.AddConfigPath(path)
		}

		if configFile != "" {
			viper.SetConfigFile(configFile)
		}

		// Process the configuration
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error as we will also pull from flags and env variables
				log.Debugf("No configuration located: %v", err.(viper.ConfigFileNotFoundError))
			} else if _, ok := err.(*fs.PathError); ok {
				log.Warnf("No configuration file located for the provided config: %v", configFile)
			} else {
				// Config file was found but another error was produced
				return errors.Errorf("loading the config file failed: %v", err)
			}
		}
	}

	log.Debugf("Loaded '%v'", configFile)

	return nil
}
