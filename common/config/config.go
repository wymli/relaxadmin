package config

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	configPath = "./config.yaml"
	c          *Config
	once       sync.Once
)

type Config struct {
	DB     *DBConfig
	Server *ServerConfig
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DB       string
}

type ServerConfig struct {
	ReloadUI  bool
	ConnectDB bool
	Host      string
	Port      string
	Env       string
}

func SetConfigPath(s string) {
	configPath = s
}

func Init() {
	once.Do(func() {
		var cc Config = Config{
			Server: &ServerConfig{
				ReloadUI:  true,
				ConnectDB: false,
				Host:      "0.0.0.0",
				Port:      "9999",
				Env:       "dev",
			},
		}

		bs, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(bs, &cc); err != nil {
			panic(err)
		}

		c = &cc
	})
}

func GetConfig() *Config {
	if c == nil {
		Init()
	}

	return c
}
