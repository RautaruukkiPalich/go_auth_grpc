package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env"`
	Database DatabaseConfig `yaml:"database" env_required:"true"`
	Server   ServerConfig   `yaml:"server" env_required:"true"`
	Token    TokenConfig    `yaml:"token" env_required:"true"`
}

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	ConnTimeout time.Duration `yaml:"conn_timeout"`
}

type TokenConfig struct {
	TTL    time.Duration `yaml:"ttl"`
	Secret string        `yaml:"secret"`
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("can not parse config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
