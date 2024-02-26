package util

import (
	"testing"
	"time"
	"context"
	"os"
	"fmt"
)

func TestLog(t *testing.T){

	Context, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(500 * time.Millisecond)
	}()

	logPath := "/tmp/log"
	MODEL := "MODEL"
	NAME := "TEST"
	logLevel := "DEBUG"
 	ID := "1"

	logger := SetupLogging(Context,logPath, MODEL, NAME, logLevel, ID)
	
	logFilePath := "/tmp/log/TEST/TEST_1_DEBUG_go.log"
	go func() {
		for {
			select {
			case <-Context.Done():
				return
			case <-time.After(1 * time.Second):
				fileInfo, _ := os.Stat(logFilePath)

				fileSizeMB := float64(fileInfo.Size()) / 1024 / 1024
				fmt.Printf("File size: %.2f MB\n", fileSizeMB)
			}
		}
	}()


	frameDuration := time.Duration( 1 * time.Nanosecond)
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	for {
		select {
		case <- Context.Done():
			return
		case <- ticker.C:
			logger.Debug("LOG_TEST")
			logger.Debug("LOG_TESTING")
		}
	}
}
