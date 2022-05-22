package config

import (
	"encoding/json"
	"os"
)

//Config - конфиг
type Config struct {
	HttpPort int
	GinMode  string //debug или release
	DB       struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
		Migrate  bool
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
		panic("load config error: " + err2.Error())
	}
	return &config
}
