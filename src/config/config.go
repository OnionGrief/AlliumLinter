package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	SnakeCase  bool `yaml:"SnakeCase"`
	CamelCase  bool `yaml:"CamelCase"`
	ConstCount uint `yaml:"ConstCount"`
	ConstLen   uint `yaml:"ConstLen"`
	BlockLen   uint `yaml:"BlockLen"`
}

var Cfg = Config{}

const FuncComment = "OPT:"

var UsingFunctions []string = []string{"Go"}

func ReadConfigFromFile(filename string) *Config {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	var config = Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil
	}
	fmt.Println("Получен конфиг из файла", config)
	return &config
}
