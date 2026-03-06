package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local" env:"APP_ENV"`

	KafkaConfig   `yaml:"kafka"`
	LogBullConfig `yaml:"logbull"`
	SMTPConfig    `yaml:"smtp"`
}

type KafkaConfig struct {
	BrokerUrl  string `yaml:"broker_url" env-required:"true" env:"KAFKA_BROKER_URL"`
	BrokerPort string `yaml:"broker_port" env-required:"true" env:"KAFKA_BROKER_CONSUMER_PORT"`
	Topic      string `yaml:"topic" env-required:"true" env:"KAFKA_TOPIC_WISHLIST_EVENTS"`
	GroupId    string `yaml:"group_id" env-required:"true" env:"KAFKA_GROUP_ID"`
}

type LogBullConfig struct {
	URL       string `yaml:"url" env:"LOGBULL_URL"`
	Port      string `yaml:"port" env-default:"4006" env:"LOGBULL_PORT"`
	ProjectID string `yaml:"api_key" env-default:"" env:"LOGBULL_PROJECT_ID"`
}

type SMTPConfig struct {
	Host     string `yaml:"host" env-default:"smtp.gmail.com" env:"SMTP_HOST"`
	Port     string `yaml:"port" env-default:"587" env:"SMTP_PORT"`
	User     string `yaml:"user" env-required:"true" env:"SMTP_USER"`
	Password string `yaml:"password" env-required:"true" env:"SMTP_PASSWORD"`
	From     string `yaml:"from" env-required:"true" env:"SMTP_FROM"`
}

func LoadConfig() *Config {
	var path string
	// как запуск с параметром --config="./some/path/cfg.yaml"
	flag.StringVar(&path, "config", "", "Path to local config. Get from env NOTIFY_CONFIG_PATH, if not set")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("NOTIFY_CONFIG_PATH")
		if len(path) == 0 {
			log.Fatal("path is empty NOTIFY_CONFIG_PATH not set")
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
