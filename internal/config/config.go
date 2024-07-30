package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	StoragePath string         `yaml:"storage_path" env-required:"true"`
	PostgreSQL  PostgresConfig `yaml:"postgresql"`
	Kafka       KafkaConfig    `yaml:"kafka"`
	HTTP        HTTPConfig     `yaml:"http"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env-required:"true"`
	Topic   string   `yaml:"topic" env-required:"true"`
}

type HTTPConfig struct {
	Port    int    `yaml:"port" env-required:"true"`
	Timeout string `yaml:"timeout" env-default:"4s"`
	Iddle   string `yaml:"iddle" env-default:"60s"`
}

func MustLoad() *Config {
	//configPath := fetchConfigPath()
	configPath := "./local.yaml"
	//if configPath == "" {
	//	panic("config path is empty")
	//}

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

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
