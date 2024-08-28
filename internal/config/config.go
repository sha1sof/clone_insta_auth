package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"prod"`
	Storage  Storage       `yaml:"storage"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
}

type Storage struct {
	Type     string         `yaml:"type"`
	Sqlite   SqliteConfig   `yaml:"sqlite"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type SqliteConfig struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port_pg" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port_grpc" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
	Secret  string        `yaml:"secret" env-required:"true"`
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not exist" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("load config file fail" + err.Error())
	}

	return &cfg
}
