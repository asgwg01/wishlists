package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local" env:"APP_ENV"`

	GRPCConfig    `yaml:"grpc"`
	StorageConfig `yaml:"storage"`
	RedisConfig   `yaml:"redis"`
	JWTConfig     `yaml:"jwt"`
	LogBullConfig `yaml:"logbull"`
}

type GRPCConfig struct {
	Port string `yaml:"port" env-default:"50051" env:"AUTH_SERVICE_PORT"`
}

type StorageConfig struct {
	Host     string `yaml:"host" nv-required:"true" env:"AUTH_DB_HOST"`
	Port     string `yaml:"port" env-default:"5432" env:"AUTH_DB_PORT"`
	DBName   string `yaml:"db_name" nv-required:"true" env:"AUTH_DB_NAME"`
	User     string `yaml:"user" nv-required:"true" env:"AUTH_DB_USER"`
	Password string `yaml:"password" nv-required:"true" env:"AUTH_DB_PASSWORD"`
}

type RedisConfig struct {
	Host     string `yaml:"host" nv-required:"true" env:"REDIS_HOST"`
	Port     string `yaml:"port" env-default:"6379" env:"REDIS_PORT"`
	Password string `yaml:"password" nv-required:"true" env:"REDIS_PASSWORD"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret" env:"JWT_SECRET"`
	TTL        time.Duration `yaml:"ttl" env-default:"1h" env:"JWT_TTL"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"24h" env:"JWT_REFRESH_TTL" `
}

type LogBullConfig struct {
	URL       string `yaml:"url" env:"LOGBULL_URL"`
	Port      string `yaml:"port" env-default:"4006" env:"LOGBULL_PORT"`
	ProjectID string `yaml:"api_key" env-default:"" env:"LOGBULL_PROJECT_ID"`
}

func LoadConfig() *Config {
	var path string
	// как запуск с параметром --config="./some/path/cfg.yaml"
	flag.StringVar(&path, "config", "", "Path to local config. Get from env AUTH_CONFIG_PATH, if not set")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("AUTH_CONFIG_PATH")
		if len(path) == 0 {
			log.Fatal("path is empty AUTH_CONFIG_PATH not set")
		}
	}

	//check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file %s is not exist", path)
	}

	//load File
	var result Config

	if err := cleanenv.ReadConfig(path, &result); err != nil {
		log.Fatalf("read file error %s", err)
	}

	return &result

}
