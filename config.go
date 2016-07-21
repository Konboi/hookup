package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port  int    `yaml:"port"`
	Hooks []Hook `yaml:"hooks"`
}

type Hook struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	reg  regexp.Regexp
}

func NewConfig(path string) (config Config, err error) {
	config = Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	if config.Port == 0 {
		return config, fmt.Errorf("please set port")
	}

	if len(config.Hooks) == 0 {
		return config, fmt.Errorf("please set hooks")
	}

	return config, nil
}

func (h Hook) Match(s string) bool {
	// todo regexp
	if strings.Compare(s, h.From) == 0 {
		return true
	}

	return false
}
