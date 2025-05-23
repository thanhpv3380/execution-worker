package services

import (
	"context"
	"encoding/json"
	"execution-worker/internal/infra/redis"
	runnerService "execution-worker/internal/services/runner"
	"fmt"
	"time"

	"github.com/thanhpv3380/go-common/logger"

	"github.com/thanhpv3380/execution-producer/internal/types/enums"

	redisTypes "github.com/thanhpv3380/execution-producer/internal/types/redis"
)

func getRunnerService(language string) runnerService.RunnerService {
	switch language {
	case "golang":
		return runnerService.GetGoRunnerService()
	default:
		logger.Warnf("Unsupported language: %s", language)
		return nil
	}
}

func initWorker(workerCount int, language string) {
	queueName := "code_jobs"

	runnerService := getRunnerService(language)
	if runnerService == nil {
		logger.Fatalf("No runner service found for language: %s", language)
	}

	for i := 1; i <= workerCount; i++ {
		go worker(i, runnerService)
	}
}

func worker(workerId int, runner runnerService.RunnerService) {
	logger.Infof("Worker %d started", workerId)

	executionId, err := listenQueue(queue)
	if err != nil {
		logger.Error(fmt.Sprintf("Worker %d listen queue error: %v", workerId), err)
		time.Sleep(time.Second)
		return
	}

	executionRaw, err := redis.Get(fmt.Sprintf("%s%s", enums.RedisKeyExecutionInfo, executionId))
	if err != nil {
		logger.Error("Error get execution in redis", err)
		return
	}

	var execution redisTypes.Execution

	err = json.Unmarshal([]byte(executionRaw), &execution)
	if err != nil {
		logger.Error("Error marshal execution", err)
		return
	}

	runner.Run(execution.C)
}

func listenQueue(ctx context.Context, queue string) (string, error) {
	result, err := redis.Client.BLPop(ctx, 0*time.Second, queue).Result()
	if err != nil {
		return "", err
	}

	if len(result) < 2 {
		return "", fmt.Errorf("unexpected BLPOP result: %v", result)
	}

	return result[1], nil
}
