package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service struct {
		Name string `yaml:"name"`
	}
	Folders struct {
		SrcFolder string `yaml:"source_folder"`
		DstFolder string `yaml:"destination_folder"`
	}
	Options struct {
		TwoWay   bool   `yaml:"two_way"`
		Internal int    `yaml:"internal"`
		LogFile  string `yaml:"logFile"`
		Debug    int    `yaml:"debug"`
	}
}

func ReadConfig(configPath string) *Config {
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var AppConfig Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&AppConfig)

	if err != nil {
		fmt.Println(err)
	}
	val := AppConfig.Options.Internal
	AppConfig.Options.Internal = val * 60
	return &AppConfig
}
