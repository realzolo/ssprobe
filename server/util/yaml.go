package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"ssprobe-server/model"
)

type Conf struct {
	model.Server   `yaml:"server"`
	model.Web      `yaml:"web"`
	model.Notifier `yaml:"notifier"`
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
