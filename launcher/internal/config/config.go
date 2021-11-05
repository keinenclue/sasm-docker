package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

// Setup reads the config file and parses it
func Setup(configPath string) error {
	return SetupWithName(configPath, "launcher")
}

// SetupWithName reads the config file and parses it and allows to set the name
func SetupWithName(configPath string, configName string) error {

	viper.Reset()
	viper.SetConfigType("yml")
	viper.SetConfigName(configName)
	viper.SetConfigPermissions(0666)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.SetDefault("dataPath", configPath)

	if err := viper.ReadInConfig(); err != nil {

		if _, e := err.(viper.ConfigFileNotFoundError); !e {
			return err
		}

		os.MkdirAll(configPath, os.ModePerm)
		viper.SetConfigFile(fmt.Sprintf("%s/%s", configPath, fmt.Sprintf("%s.yml", configName)))
		if err = viper.WriteConfig(); err != nil {
			return err
		}
	}

	viper.WatchConfig()

	return nil
}

// GetConfig returns the config
func Get(key string) interface{} {
	return viper.Get(key)
}

// PathUsed returns the used config path
func PathUsed() string {
	return path.Dir(viper.ConfigFileUsed())
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
	viper.WriteConfig()
}
