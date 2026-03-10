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

	ServerConfig   `yaml:"http"`
	WishlistConfig `yaml:"wishlist"`
	LogBullConfig  `yaml:"logbull"`
}

type ServerConfig struct {
	Addres  string        `yaml:"addres" env-default:"0.0.0.0"` // не WEB_CLIENT_HOST
	Port    string        `yaml:"port" env-default:"8095" env:"WEB_CLIENT_PORT"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"WEB_CLIENT_HTTP_TIMEOUT"`
	Secret  string        `yaml:"secret" env-required:"true" env:"WEB_CLIENT_SECRET"`
}

// конфиг для связи с приложением WishlistApp
type WishlistConfig struct {
	GatewayAddres string `yaml:"gateway_addres" env-default:"localhost" env:"API_GATEWAY_HTTP_ADRESS"`
	GatewayPort   string `yaml:"gateway_port" env-default:"8095" env:"API_GATEWAY_HTTP_PORT"`
	ApiURL        string `api_url:"gateway_port" env-default:"/api/v1" env:"WEB_CLIENT_WISHLIST_API_URL"`
}

type LogBullConfig struct {
	URL       string `yaml:"url" env:"LOGBULL_URL"`
	Port      string `yaml:"port" env-default:"4006" env:"LOGBULL_PORT"`
	ProjectID string `yaml:"api_key" env-default:"" env:"LOGBULL_PROJECT_ID"`
}

func LoadConfig() *Config {
	var path string
	// как запуск с параметром --config="./some/path/cfg.yaml"
	flag.StringVar(&path, "config", "", "Path to local config. Get from env WEB_CLIENT_CONFIG_PATH, if not set")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("WEB_CLIENT_CONFIG_PATH")
		if len(path) == 0 {
			log.Fatal("path is empty WEB_CLIENT_CONFIG_PATH not set")
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
