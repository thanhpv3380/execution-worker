package services

import (
	"sync"
)

type goRunnerService struct {
}

var (
	goRunnerServiceInstance RunnerService
	once                    sync.Once
)

func GetGoRunnerService() RunnerService {
	once.Do(func() {
		goRunnerServiceInstance = &goRunnerService{}
	})
	return goRunnerServiceInstance
}

func (g *goRunnerService) Run(code string) error {
	return nil
}
