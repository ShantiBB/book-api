package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTPServer struct {
	Address string        `yaml:"address" env-default:"localhost"`
	Port    string        `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		fmt.Println("failed to load .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		fmt.Println("config path variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("config path does not exist: %s\n", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		fmt.Printf("failed to read config: %s\n", err)
	}

	return &config
}
