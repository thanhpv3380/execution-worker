package services

import (
	"bytes"
	"context"
	"execution-worker/internal/configs"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/thanhpv3380/go-common/logger"
)

type goRunnerService struct {
}

var (
	goRunnerServiceInstance RunnerService
	goOnce                  sync.Once
	ctx                     = context.Background()
)

func GetGoRunnerService() RunnerService {
	goOnce.Do(func() {
		goRunnerServiceInstance = &goRunnerService{}
	})
	return goRunnerServiceInstance
}

func (g *goRunnerService) Run(code string) (string, error) {
	err := os.MkdirAll(configs.Cfg.ExecuteTempDir, 0755)
	if err != nil {
		return "", err
	}

	mainFile := filepath.Join(configs.Cfg.ExecuteTempDir, fmt.Sprintf("%s.go", uuid.New().String()))

	err = os.WriteFile(mainFile, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	defer func() {
		err = os.Remove(mainFile)
		if err != nil {
			logger.Error("Error remove temp file", err)
		}
	}()

	cmd := exec.CommandContext(ctx,
		"go", "run", mainFile,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%s\n%s\n%s", err.Error(), stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}
