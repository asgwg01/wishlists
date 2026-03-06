package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local" env:"APP_ENV"`

	HTTPPort string `yaml:"http_port" env-default:"8082"`

	GRPCConfig    `yaml:"grpc"`
	StorageConfig `yaml:"storage"`
	KafkaConfig   `yaml:"kafka"`
	LogBullConfig `yaml:"logbull"`
}

type GRPCConfig struct {
	Port            string `yaml:"port" env-default:"50051" env:"WISHLIST_SERVICE_PORT"`
	AuthServiceAddr string `yaml:"auth_service_addr" env-default:"localhost" env:"AUTH_SERVICE_HOST"`
	AuthServicePort string `yaml:"auth_service_port" env-default:"50051" env:"AUTH_SERVICE_PORT"`
}

type StorageConfig struct {
	Host     string `yaml:"host" nv-required:"true" env:"WISHLIST_DB_HOST"`
	Port     string `yaml:"port" env-default:"5432" env:"WISHLIST_DB_PORT"`
	DBName   string `yaml:"db_name" nv-required:"true" env:"WISHLIST_DB_NAME"`
	User     string `yaml:"user" nv-required:"true" env:"WISHLIST_DB_USER"`
	Password string `yaml:"password" nv-required:"true" env:"WISHLIST_DB_PASSWORD"`
}

type KafkaConfig struct {
	Broker string `yaml:"broker" nv-required:"true" env:"KAFKA_BROKER"`
	Topic  string `yaml:"topic" nv-required:"true" env:"KAFKA_TOPIC_WISHLIST_EVENTS"`
}

type LogBullConfig struct {
	URL       string `yaml:"url" env:"LOGBULL_URL"`
	Port      string `yaml:"port" env-default:"4006" env:"LOGBULL_PORT"`
	ProjectID string `yaml:"api_key" env-default:"" env:"LOGBULL_PROJECT_ID"`
}

func LoadConfig() *Config {
	var path string
	// как запуск с параметром --config="./some/path/cfg.yaml"
	flag.StringVar(&path, "config", "", "Path to local config. Get from env WISHLIST_CONFIG_PATH, if not set")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("WISHLIST_CONFIG_PATH")
		if len(path) == 0 {
			log.Fatal("path is empty WISHLIST_CONFIG_PATH not set")
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
