package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `yaml:"is_debug" env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `yaml:"is_development" env:"IS_DEVELOPMENT" env-default:"false"`
	Server        struct {
		ServerHost string `yaml:"server_host" env:"SERVER_HOST"`
		ServerPort string `yaml:"server_port" env:"SERVER_PORT"`
	} `yaml:"server"`
	Client struct {
		ClientHost string `yaml:"client_host" env:"CLIENT_IP"`
		ClientPort string `yaml:"client_port" env:"CLIENT_PORT"`
	} `yaml:"client"`
	Redis struct {
		Host string `yaml:"host" env:"REDIS_HOST"`
		Port string `yaml:"port" env:"REDIS_PORT"`
	} `yaml:"redis"`
	HashCash struct {
		RequiredZerosCount int   `yaml:"required_zero_count" env:"REQUIRED_ZEROS_COUNT"` // The number of leading zeros required in the hash for proof of work (server-side only)
		ChallengeLifetime  int64 `yaml:"challenge_life_time" env:"CHALLENGE_LIFE_TIME"`  // Lifetime of the challenge (server-side only)
		MaxIterations      int   `yaml:"max_iterations" env:"MAX_ITERATION"`             // Maximum iterations to prevent getting stuck on hard hashes (client-side only)
	}
}

const (
	EnvConfigPathName  = "CONFIG-PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(&configPath, FlagConfigPathName, "/Users/romanzaitsev/Desktop/faraway/configs/config.yaml", "this is app config file")
		flag.Parse()

		log.Print("config init")

		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			log.Fatal("config path is required")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			helpText := "faraway - test task client/server(PoW - Hashcash def DDoS)"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
