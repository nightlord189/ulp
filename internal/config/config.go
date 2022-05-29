package config

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"os"
)

//Config - конфиг
type Config struct {
	HttpPort              int `envconfig:"HTTP_PORT"`
	HttpDebug             bool
	TemplatesPath         string
	AttemptsPath          string
	RunTestsTimeout       int
	TestPortStart         int
	TestHost              string `envconfig:"TEST_HOST"`
	ContainerArchitecture string `envconfig:"CONTAINER_ARCH"`
	DB                    struct {
		Host     string `envconfig:"DB_HOST"`
		Port     int    `envconfig:"DB_PORT"`
		Name     string `envconfig:"DB_NAME"`
		User     string `envconfig:"DB_USER"`
		Password string `envconfig:"DB_PASSWORD"`
		Migrate  bool
		Log      bool
	}
	Auth struct {
		Secret  string
		Issuer  string
		ExpTime int
	}
}

//Load - загрузить конфиг по указанному пути
func Load(path string) *Config {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		panic("load config error: " + err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic("parse config error: " + err.Error())
	}
	err = envconfig.Process("", &config)
	if err != nil {
		panic("get env error: " + err.Error())
	}
	return &config
}
