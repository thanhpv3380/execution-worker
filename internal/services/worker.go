package services

import (
	"encoding/json"
	runnerService "execution-worker/internal/services/runner"
	"fmt"
	"time"

	"github.com/thanhpv3380/go-common/logger"

	"github.com/thanhpv3380/execution-producer/pkg/redis"
	"github.com/thanhpv3380/execution-producer/pkg/types/enums"

	redisTypes "github.com/thanhpv3380/execution-producer/pkg/types/redis"
)

func InitWorker(workerId int, queue string, runner runnerService.RunnerService) {
	logger.Infof("Worker %d started", workerId)

	for {
		executionId, err := redis.BLPop(queue)
		if err != nil {
			logger.Error(fmt.Sprintf("Worker %d listen queue error", workerId), err)
			time.Sleep(time.Second)
			continue
		}

		logger.Infof("Worker %d received execution ID: %s", workerId, executionId)

		executionInfoRedisKey := fmt.Sprintf("%s%s", enums.RedisKeyExecutionInfo, executionId)
		executionRaw, err := redis.Get(executionInfoRedisKey)
		if err != nil {
			logger.Error("Error get execution in redis", err)
			continue
		}

		var execution redisTypes.Execution

		err = json.Unmarshal([]byte(executionRaw), &execution)
		if err != nil {
			logger.Error("Error marshal execution", err)
			continue
		}

		result, err := runner.Run(execution.Code)

		now := time.Now()
		execution.FinishedAt = &now

		if err != nil {
			logger.Error(fmt.Sprintf("Worker %d run execution error", workerId), err)
			execution.Result = err.Error()
			execution.Status = enums.ExecuteStatusFailed
		} else {
			logger.Infof("Worker %d execution result: %s", workerId, result)
			execution.Status = enums.ExecuteStatusCompleted
			execution.Result = result
		}

		executionByte, err := json.Marshal(execution)
		if err != nil {
			logger.Error("Error marshal execution", err)
			continue
		}

		redis.Set(executionInfoRedisKey, executionByte, -1)
	}
}
