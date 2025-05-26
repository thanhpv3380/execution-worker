package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	runnerService "github.com/thanhpv3380/execution-worker/internal/services/runner"

	"github.com/thanhpv3380/go-common/logger"

	"github.com/thanhpv3380/execution-producer/pkg/redis"
	"github.com/thanhpv3380/execution-producer/pkg/types/enums"

	redisTypes "github.com/thanhpv3380/execution-producer/pkg/types/redis"
)

func InitWorker(ctx context.Context, queue string, runner runnerService.RunnerService) {
	loggerCtx := logger.FromContext(ctx)
	loggerCtx.Info("Worker started")

	defer loggerCtx.Info("Worker finished")

	for {
		executionId, err := redis.BLPop(queue)
		if err != nil {
			loggerCtx.Error("Worker listen queue error", err)
			time.Sleep(time.Second)
			continue
		}

		subCtx := logger.AppendLogFields(ctx, map[string]interface{}{"executionId": executionId})
		subLoggerCtx := logger.FromContext(subCtx)

		subLoggerCtx.Infof("Received execution")

		executionInfoRedisKey := fmt.Sprintf("%s%s", enums.RedisKeyExecutionInfo, executionId)

		executionRaw, err := redis.Get(executionInfoRedisKey)
		if err != nil {
			subLoggerCtx.Errorw("Error get execution in redis", err)
			continue
		}

		var execution redisTypes.Execution

		err = json.Unmarshal([]byte(executionRaw), &execution)
		if err != nil {
			subLoggerCtx.Errorw("Error marshal execution", err)
			continue
		}

		result, err := runner.Run(subCtx, execution.Code)

		now := time.Now()
		execution.FinishedAt = &now

		if err != nil {
			subLoggerCtx.Errorw("Error run execution", err)
			execution.Result = err.Error()
			execution.Status = enums.ExecuteStatusFailed
		} else {
			subLoggerCtx.Infof("Execution result: %s", result)
			execution.Status = enums.ExecuteStatusCompleted
			execution.Result = result
		}

		executionByte, err := json.Marshal(execution)
		if err != nil {
			subLoggerCtx.Errorw("Error marshal execution", err)
			continue
		}

		ttl, err := redis.TTL(executionInfoRedisKey)
		if err != nil {
			subLoggerCtx.Errorw("Error get TTL", err)
			continue
		}

		if ttl == -2 || ttl == 0 {
			subLoggerCtx.Warnf("No TTL or key not exist: %s", ttl)
			continue
		}

		redis.Set(executionInfoRedisKey, executionByte, ttl)
		subLoggerCtx.Infof("Save execution result to redis successfully")
	}

}
