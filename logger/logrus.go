package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"
	"context"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000000000")
	logLevel := strings.ToUpper(entry.Level.String())

	var file string
	var line int

	if entry.HasCaller() {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}

	msg := fmt.Sprintf("%s | %s | [%s:%d] %s\n", timestamp, logLevel, file, line, entry.Message)
	return []byte(msg), nil
}

func SetupLogging(ctx context.Context, logPath, model, name, logLevel, pipelineID string) *logrus.Logger {
	logLevel = strings.ToUpper(logLevel)
	maxSizeMB := int64(100)
	maxFiles := 10


	if logPath == "" {
		logPath = "/tmp/log"
	}
	logDir := filepath.Join(logPath, name)

	if err := os.MkdirAll(logDir, 0750); err != nil {
		logrus.Fatalf("Failed to create log directory: %v", err)
	}

	logFileName := fmt.Sprintf("%s_%s_%s_go.log", name, pipelineID, logLevel)
	logFilePath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	logger := logrus.New()
	logger.SetOutput(file)
	logger.SetFormatter(new(CustomFormatter))
	logger.SetReportCaller(true)

	switch logLevel {
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
	case "CRITICAL":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	go rotateLogFile(ctx, logger, logFilePath,maxSizeMB, maxFiles)


	return logger
}

func findLatestLogFile(logDir string) string {
	files, err := os.ReadDir(logDir)
	if err != nil {
		logrus.Fatalf("Failed to read log directory: %v", err)
	}

	var latestFile os.DirEntry
	var latestIdx int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		idx, err := parseLogIndex(name)
		if err == nil && idx > latestIdx {
			latestIdx = idx
			latestFile = file
		}
	}
	if latestFile != nil {
		return latestFile.Name()
	}
	return ""
}

func parseLogIndex(fileName string) (int, error) {
	base := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	base = strings.TrimSuffix(base, ext)
	parts := strings.Split(base, "_")
	
	if len(parts) > 1 {
		return strconv.Atoi(parts[len(parts)-1])
	}
	return 0, fmt.Errorf("log file name does not contain index")
}

func rotateLogFile(ctx context.Context, logger *logrus.Logger, logFilePath string, maxSizeMB int64, maxFiles int) {
    checkLogDuration := 1 * time.Second
    checkSizeTicker := time.NewTicker(checkLogDuration)
    defer checkSizeTicker.Stop()

    for {
        select {
        case <-ctx.Done(): // 종료 신호를 받음
            logger.Info("Log rotation stopped due to context cancellation")
            return
        case <-checkSizeTicker.C:
            fileInfo, err := os.Stat(logFilePath)
            if err != nil {
                logger.Errorf("Failed to get log file info: %v", err)
                continue
            }

            if fileInfo.Size() < maxSizeMB*1024*1024 {
                continue
            }

            for i := maxFiles - 1; i > 0; i-- {
                oldPath := fmt.Sprintf("%s.%d", logFilePath, i-1)
                if i == 1 {
                    oldPath = logFilePath
                }
                newPath := fmt.Sprintf("%s.%d", logFilePath, i)

                if _, err := os.Stat(oldPath); os.IsNotExist(err) {
                    continue
                }
                os.Rename(oldPath, newPath)
            }

            os.Truncate(logFilePath, 0)
            file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
            if err != nil {
                logger.Errorf("Failed to open new log file: %v", err)
                return
            }
            logger.SetOutput(file)
        }
    }
}




func incrementLogFileName(basePath string, index int) string {
	dir, file := filepath.Split(basePath)
	ext := filepath.Ext(file)
	base := strings.TrimSuffix(file, ext)

	newFileName := fmt.Sprintf("%s_%d%s", base, index, ext)
	return filepath.Join(dir, newFileName)
}

func resetFileIndexIfNecessary(currentIndex, maxIndex int) int {
	if currentIndex >= maxIndex {
		return 0
	}
	return currentIndex + 1
}


