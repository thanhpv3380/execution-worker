package main

import (
	"execution-worker/internal/configs"
	"execution-worker/internal/infra/redis"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/thanhpv3380/go-common/logger"
)

func main() {
	var wg sync.WaitGroup
	stop := make(chan struct{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg := configs.LoadConfig()
	logger.NewLogger(nil)

	initRedis(cfg)

	<-sigs
	close(stop)

	wg.Wait()
	logger.Info("Service is stopped, exiting.")
}

func initRedis(cfg *configs.Config) {
	redisAddress := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	logger.Infof("Initing Redis client ..., address: %s", redisAddress)
	err := redis.NewClient(redisAddress, cfg.Redis.Password)
	if err != nil {
		logger.Fatal("Failed to connect redis", err)
	}

	logger.Info("Init Redis client successfully")
}
