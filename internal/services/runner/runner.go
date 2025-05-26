package services

import (
	"context"

	"github.com/thanhpv3380/execution-producer/pkg/types/enums"
	"github.com/thanhpv3380/go-common/logger"
)

type RunnerService interface {
	Run(ctx context.Context, code string) (string, error)
}

func GetRunnerService(language enums.ProgrammingLanguage) RunnerService {
	switch language {
	case enums.Golang:
		return GetGoRunnerService()
	case enums.Javascript:
		return GetJSRunnerService()
	default:
		logger.Warnf("Unsupported programming language: %s", language)
		return nil
	}
}
