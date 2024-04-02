package main

import (
	"github.com/CDN-Security/Echo/pkg/config"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := yaml.Marshal(config.DefaultConfig)
	if err != nil {
		panic(err)
	}
	println(string(data))
}
