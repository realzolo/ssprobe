package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Token string `yaml:"token"`
	Port  struct {
		Server    int `yaml:"server"`
		Web       int `yaml:"web"`
		Websocket int `yaml:"websocket"`
	} `yaml:"port"`
}

func (c *Conf) GetConf() (*Conf, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
