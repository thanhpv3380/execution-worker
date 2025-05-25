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
	Redis          Redis
	WorkerCount    int
	WorkerLanguage string
	ExecuteTempDir string
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
		WorkerCount:    env.GetInt("WORKER_COUNT", 1),
		WorkerLanguage: env.GetString("WORKER_LANGUAGE", "golang"),
		ExecuteTempDir: env.GetString("EXECUTE_TEMP_DIR", "./tmp/sandbox"),
	}

	log.Println("Load config successfully")

	return Cfg
}
