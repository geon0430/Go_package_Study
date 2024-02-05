package util

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestIncrementLogFileName(t *testing.T) {
	basePath := "/tmp/log_test_go.log"
	expected := "/tmp/log_test_go_1.log"
	actual := incrementLogFileName(basePath, 1)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestRotateLogFile(t *testing.T) {
	logPath := "/tmp/log"
	logLevel := "INFO"
	name := "testlog"

	logger := SetupLogging(logPath, logLevel, name)
	defer func() {
		if file, ok := logger.Out.(*os.File); ok {
			file.Close()
		}
	}()

	logFilePath := filepath.Join(logPath, fmt.Sprintf("%s_%s_go.log", name, logLevel))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			logger.Info("Test log entry")
			time.Sleep(1 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	file, _ := logger.Out.(*os.File)
	rotateLogFile(logger, file, logFilePath, 1)

	rotatedFilePath := incrementLogFileName(logFilePath, 1)
	if _, err := os.Stat(rotatedFilePath); os.IsNotExist(err) {
		t.Errorf("Rotated log file does not exist: %s", rotatedFilePath)
	}

	info, err := os.Stat(logFilePath)
	if err != nil {
		t.Fatalf("Failed to stat the log file: %v", err)
	}
	if info.Size() > 1*1024*1024 {
		t.Errorf("Log file size did not decrease after rotation")
	}

}
