package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/thanhpv3380/execution-worker/internal/configs"

	"github.com/google/uuid"
	"github.com/thanhpv3380/go-common/logger"
)

type goRunnerService struct {
}

var (
	goRunnerServiceInstance RunnerService
	goOnce                  sync.Once
)

func GetGoRunnerService() RunnerService {
	goOnce.Do(func() {
		goRunnerServiceInstance = &goRunnerService{}
	})
	return goRunnerServiceInstance
}

func (g *goRunnerService) Run(ctx context.Context, code string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	loggerCtx := logger.FromContext(ctx)

	err = os.MkdirAll(configs.Cfg.ExecuteTempDir, 0755)
	if err != nil {
		return "", err
	}

	mainFile := filepath.Join(configs.Cfg.ExecuteTempDir, fmt.Sprintf("%s.go", uuid.New().String()))

	err = os.WriteFile(mainFile, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	defer func() {
		if removeErr := os.Remove(mainFile); removeErr != nil {
			loggerCtx.Errorw("Error removing temp file", removeErr)
		}
	}()

	executeCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(executeCtx, "go", "run", mainFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s\n%s\n%s", err.Error(), stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}
