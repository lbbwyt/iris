package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	WebServer Web  `yaml:"web"`
	TcpServer Web  `yaml:"tcp"`
	Uart      Uart `yaml:"uart"`
}

type Uart struct {
	PortName string `yaml:"port_name"`
	BaudRate uint   `yaml:"baud_rate"`
}

type Web struct {
	Port string `yaml:"port"`
}

var GConfig *Config

func Init(cfgPath string) error {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	GConfig = c
	return nil
}
