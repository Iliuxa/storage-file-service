package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env                 string     `yaml:"env" env-default:"local"`
	StoragePath         string     `yaml:"storage_path" env-required:"true"`
	GRPC                GRPCConfig `yaml:"grpc"`
	ListLimit           int        `yaml:"listLimit"`
	DownloadUploadLimit int        `yaml:"downloadUploadLimit"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad(defaultPath string) *Config {
	configPath := fetchConfigPath(defaultPath)
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath(defaultPath string) string {
	if defaultPath == "" {
		defaultPath = "config/config.yaml"
	}
	return defaultPath
}
