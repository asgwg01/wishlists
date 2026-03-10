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

	ServerConfig     `yaml:"http"`
	GRPCConfig       `yaml:"grpc"`
	JWTConfig        `yaml:"jwt"`
	CORSConfig       `yaml:"cors"`
	RateLimitsConfig `yaml:"rate_limits"`
	SwaggerConfig    `yaml:"swagger"`
	LogBullConfig    `yaml:"logbull"`
}

type ServerConfig struct {
	Addres  string        `yaml:"addres" env-default:"0.0.0.0"` //localhost, не API_GATEWAY_HTTP_ADRESS
	Port    string        `yaml:"port" env-default:"8095" env:"API_GATEWAY_HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"API_GATEWAY_HTTP_TIMEOUT"`
}

type GRPCConfig struct {
	AuthServiceAddr     string `yaml:"auth_service_addr" env-default:"localhost" env:"AUTH_SERVICE_HOST"`
	AuthServicePort     string `yaml:"auth_service_port" env-default:"50051" env:"AUTH_SERVICE_PORT"`
	WishlistServiceAddr string `yaml:"wishlist_service_addr" env-default:"localhost" env:"WISHLIST_SERVICE_HOST"`
	WishlistServicePort string `yaml:"wishlist_service_port" env-default:"50052" env:"WISHLIST_SERVICE_PORT"`
}

type JWTConfig struct {
	Secret string `yaml:"secret" env:"JWT_SECRET"`
}

type CORSConfig struct {
	Origins     []string `yaml:"origins" env:"ALLOWED_ORIGINS"`
	Methods     []string `yaml:"metods" env:"ALLOWED_METHODS"`
	Headers     []string `yaml:"headers" env:"ALLOWED_HEADERS"`
	Credendials bool     `yaml:"credentials" env:"CORS_ALLOW_CREDENTIALS"`
}

type RateLimitsConfig struct {
	Limit    int           `yaml:"limit" env:"RATE_LIMIT_REQUESTS"`
	Duration time.Duration `yaml:"duration" env:"RATE_LIMIT_DURATION"`
}

type SwaggerConfig struct {
	NeedRuning bool   `yaml:"run_swagger_ui" env-default:"false"`
	URL        string `yaml:"swagger_url" env-default:"/swagger/"`
}

type LogBullConfig struct {
	URL       string `yaml:"url" env:"LOGBULL_URL"`
	Port      string `yaml:"port" env-default:"4006" env:"LOGBULL_PORT"`
	ProjectID string `yaml:"api_key" env-default:"" env:"LOGBULL_PROJECT_ID"`
}

func LoadConfig() *Config {
	var path string
	// как запуск с параметром --config="./some/path/cfg.yaml"
	flag.StringVar(&path, "config", "", "Path to local config. Get from env API_GATEWAY_CONFIG_PATH, if not set")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("API_GATEWAY_CONFIG_PATH")
		if len(path) == 0 {
			log.Fatal("path is empty API_GATEWAY_CONFIG_PATH not set")
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
