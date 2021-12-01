package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port            string `yaml:"port"`
	BaseURL         string `yaml:"baseurl"`
	UserEndpoint    string `yaml:"user_endpoint"`
	AddressEndpoint string `yaml:"address_endpoint"`
}

var Configuration Config

func init() {
	configBytes, err := ioutil.ReadFile("pkg/config/conf.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configBytes, &Configuration)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully parsed config: %+v\n", Configuration)
}
