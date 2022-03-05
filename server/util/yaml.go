package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"ssprobe-server/model"
)

type Conf struct {
	model.Server   `yaml:"server"`
	model.Web      `yaml:"web"`
	model.Notifier `yaml:"notifier"`
}

func (c *Conf) LoadConfig() error {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		yamlFile, err = ioutil.ReadFile("config.yml")
		if err != nil {
			return err
		}
	}
	return yaml.Unmarshal(yamlFile, c)
}

func (c *Conf) SetOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if reflect.TypeOf(value).Name() == "bool" {
		return defaultValue
	}
	if value != "" && value != 0 {
		return value
	}
	return defaultValue
}
