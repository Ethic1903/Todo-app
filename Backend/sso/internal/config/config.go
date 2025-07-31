package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	if err := godotenv.Load("local.env"); err != nil {
		log.Fatal("could not find .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}

//func MustLoad() *Config {
//	path := fetchConfigPath()
//	if path == "" {
//		panic("config path is empty")
//	}
//
//	return MustLoadByPath(path)
//}
//
//func MustLoadByPath(configPath string) *Config {
//	if _, err := os.Stat(configPath); os.IsNotExist(err) {
//		panic("config file does not exist" + configPath)
//	}
//	var cfg Config
//
//	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
//		panic("can't read config: " + err.Error())
//	}
//	return &cfg
//}
//
//func fetchConfigPath() string {
//	var res string
//
//	flag.StringVar(&res, "config", "", "path to config file")
//	flag.Parse()
//
//	if res == "" {
//		res = os.Getenv("CONFIG_PATH")
//	}
//	return res
//}
