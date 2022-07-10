package utils

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
)

var Config *config
var ConfigPath *string

func InitConfig() {
	Config = &config{}
	Config.Reload()
}

type config struct {
	Name             string `yaml:"Name"`
	IP               string `yaml:"IP"`
	Port             int    `yaml:"Port"`
	MaxConnection    uint32 `yaml:"MaxConnection"`
	MaxPackageSize   uint32 `yaml:"MaxPackageSize"`
	MaxWorkPoolSize  uint32 `yaml:"MaxWorkPoolSize"`
	MaxTaskQueueSize uint32 `yaml:"MaxTaskQueueSize"`
}

func (c *config) Reload() {
	fmt.Printf("==> ConfigPath is [%s]\n", *ConfigPath)
	data, err := ioutil.ReadFile(*ConfigPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
}
