package configs

import (
	"log"

	"github.com/thanhpv3380/go-common/env"
)

type Redis struct {
	Host     string
	Port     int
	Password string
}

type Config struct {
	Redis               Redis
	ExecutionExpireTime int
	WorkerCount         int
	WorkerLanguage      string
}

var Cfg *Config

func LoadConfig() *Config {
	if err := env.LoadEnv(); err != nil {
		log.Println("No .env file found")
	}

	Cfg = &Config{

		Redis: Redis{
			Host:     env.GetString("REDIS_HOST", ""),
			Port:     env.GetInt("REDIS_PORT", 6379),
			Password: env.GetString("REDIS_PASSWORD", ""),
		},
		ExecutionExpireTime: env.GetInt("EXECUTION_EXPIRE_TIME", 300), // seconds
		WorkerCount:         env.GetInt("WORKER_COUNT", 1),
		WorkerLanguage:      env.GetString("WORKER_LANGUAGE", "golang"),
	}

	log.Println("Load config successfully")

	return Cfg
}
