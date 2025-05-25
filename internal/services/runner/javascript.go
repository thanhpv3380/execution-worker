package services

import (
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

func (g *jsRunnerService) Run(code string) (string, error) {
	return "", fmt.Errorf("Not supported yet")
}
