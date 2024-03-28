package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultConfigPath = "./config/local.yaml"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server" env-required:"true"`
	JwtSecret   string `yaml:"jwt_secret" env-requires:"true"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"3s"`
	TokenTL      time.Duration `yaml:"token_tl" env-default:"1h"`
}

// loading config from configPath
func MustLoad() *Config {

	// getting path of configuration file
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if config file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file doesn't exist:" + configPath)
	}

	// reading config
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("unable to read config file")
	}

	return &cfg
}

// fetch config path
// priority: command line flags (--config="pathtoconfig") > environmental variables > default
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG")
	}

	if path == "" {
		path = defaultConfigPath
	}

	return path
}
