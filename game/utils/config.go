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
	Name             string `yaml:"name"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	MaxConnection    uint32    `yaml:"MaxConnection"`
	MaxPackageSize   uint32 `yaml:"maxPackageSize"`
	MaxWorkPoolSize  uint32 `yaml:"maxWorkPoolSize"`
	MaxTaskQueueSize uint32 `yaml:"maxTaskQueueSize"`
}

func (c *config) Reload()  {
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