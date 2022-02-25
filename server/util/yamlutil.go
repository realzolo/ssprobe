package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Conf struct {
	Token string `yaml:"token"`
	Port  struct {
		Server int `yaml:"server"`
		WebApi int `yaml:"web-api"`
	} `yaml:"port"`
}

func (c *Conf) GetConf() *Conf {
	path, _ := filepath.Abs("config.yaml")
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Yaml file load error. %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
