package ctxlog

import (
	"strings"

	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"time"
)

type loggerSettings struct {
	logDir string
}

type loggerOption func(*loggerSettings)

func WithLogDir(logDir string) loggerOption {
	return func(s *loggerSettings) {
		s.logDir = logDir
	}
}

func NewCmdLogger(opts ...loggerOption) *zap.Logger {
	cfg := zap.NewProductionConfig()
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		cfg.Level.SetLevel(zap.DebugLevel)
	}

	settings := loggerSettings{
		logDir: "./logs",
	}

	for _, opt := range opts {
		opt(&settings)
	}

	if _, err := os.Stat(settings.logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(settings.logDir, 0755); err != nil {
			panic(fmt.Errorf("error creating log directory: %v", err))
		}
	}

	fileName := os.Getenv("LOG_FILE_NAME")
	if fileName == "" {
		fileName = "run"
	}
	fileName = strings.TrimSuffix(fileName, ".log")

	filePath := filepath.Join(settings.logDir, fmt.Sprintf("%s-%s.log", fileName, time.Now().Format("2006-01-02_15-04-05")))
	latestFilePath := filepath.Join(settings.logDir, fmt.Sprintf("latest-%s.log", fileName))
	cfg.OutputPaths = []string{filePath, "stdout"}

	// f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	panic(fmt.Errorf("error opening file: %s: %w", filePath, err))
	// }
	// if err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd())); err != nil {
	// 	panic(fmt.Errorf("error redirecting stderr to file: %w", err))
	// }

	os.Remove(latestFilePath)

	err := os.Symlink(filePath, latestFilePath)
	if err != nil {
		panic(fmt.Errorf("error symlinking file: %s: %w", latestFilePath, err))
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("error creating logger: %v", err))
	}

	return logger
}
