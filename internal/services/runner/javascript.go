package services

import (
	"context"
	"fmt"
	"sync"
)

type jsRunnerService struct {
}

var (
	jsRunnerServiceInstance RunnerService
	jsOnce                  sync.Once
)

func GetJSRunnerService() RunnerService {
	jsOnce.Do(func() {
		jsRunnerServiceInstance = &jsRunnerService{}
	})
	return jsRunnerServiceInstance
}

func (g *jsRunnerService) Run(ctx context.Context, code string) (string, error) {
	return "", fmt.Errorf("not supported yet")
}
