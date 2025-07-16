package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"./data"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoadConfig() *Config {
	path := fetchingConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	//проверим существует ли такой файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}
	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic(err)
	}
	return &config
}

// будет получать инфу о пути до файла конфига либо из переменной окрудения либо из флага
func fetchingConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "./config/local.yaml", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
