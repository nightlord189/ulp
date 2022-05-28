package config

import (
	"encoding/json"
	"os"
)

//Config - конфиг
type Config struct {
	HttpPort        int
	HttpDebug       bool
	TemplatesPath   string
	AttemptsPath    string
	RunTestsTimeout int
	TestPortStart   int
	DB              struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
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
	err2 := decoder.Decode(&config)
	if err2 != nil {
		panic("parse config error: " + err2.Error())
	}
	return &config
}
