package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Port    string `yaml:"port"`
	BaseURL string `yaml:"baseurl"`
	Mysql   struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"dbname"`
	}
	SNS     struct {
		Region         string `yaml:"region"`
		TopicArn       string `yaml:"topicArn"`
	}
	AWS     struct{
		PublicKey string `yaml:"publicKey"`
		SecretKey string `yaml:"secretKey"`
	}
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
