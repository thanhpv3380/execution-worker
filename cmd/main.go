package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/thanhpv3380/execution-worker/internal/configs"
	services "github.com/thanhpv3380/execution-worker/internal/services"
	runnerServices "github.com/thanhpv3380/execution-worker/internal/services/runner"

	"github.com/thanhpv3380/execution-producer/pkg/redis"
	"github.com/thanhpv3380/execution-producer/pkg/types/enums"
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
	initWorkers(cfg)

	<-sigs
	close(stop)

	wg.Wait()
	logger.Info("Service is stopped, exiting.")
}

func initWorkers(cfg *configs.Config) {
	language := enums.ProgrammingLanguage(cfg.WorkerLanguage)

	runnerService := runnerServices.GetRunnerService(language)
	if runnerService == nil {
		logger.Fatalf("No runner service found for language: %s", language)
	}

	rootCtx := context.Background()
	for workerId := 1; workerId <= cfg.WorkerCount; workerId++ {
		ctx := context.WithValue(rootCtx, logger.ContextLogFieldsKey, map[string]interface{}{
			"workerId": workerId,
		})
		go services.InitWorker(ctx, fmt.Sprintf("%s%s", enums.RedisKeyExecutionQueue, language), runnerService)
	}
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
