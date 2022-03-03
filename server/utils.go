package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Server `yaml:"server"`
	Web    `yaml:"web"`
}
type Server struct {
	Token         string `yaml:"token"`
	Port          int    `yaml:"port"`
	WebsocketPort int    `yaml:"websocketPort"`
}
type Web struct {
	Enable bool   `yaml:"enable"`
	Title  string `yaml:"title"`
}

func (c *Conf) LoadConfig() (*Conf, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		yamlFile, err = ioutil.ReadFile("config.yml")
		if err != nil {
			return nil, err
		}
	}
	return c, yaml.Unmarshal(yamlFile, c)
}

func (c *Conf) SetOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if value != "" && value != 0 {
		return value
	}
	return defaultValue
}
