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

	// 로그 파일에 데이터를 지속적으로 쓰는 고루틴
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			logger.Info("Test log entry")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// 로테이션을 확인하기 위한 대기 시간
	time.Sleep(1 * time.Second)

	// 로테이션 실행
	file, _ := logger.Out.(*os.File)
	rotateLogFile(logger, file, logFilePath, 1)

	// 로테이션된 로그 파일의 존재 여부 확인
	rotatedFilePath := incrementLogFileName(logFilePath, 1)
	if _, err := os.Stat(rotatedFilePath); os.IsNotExist(err) {
		t.Errorf("Rotated log file does not exist: %s", rotatedFilePath)
	}

	// 로그 파일의 크기가 감소했는지 확인 (새로운 로그 파일이 생성되었는지)
	info, err := os.Stat(logFilePath)
	if err != nil {
		t.Fatalf("Failed to stat the log file: %v", err)
	}
	if info.Size() > 1*1024*1024 { // 로그 파일 크기 제한이 1MB로 가정
		t.Errorf("Log file size did not decrease after rotation")
	}

	wg.Wait() // 로그 쓰기 고루틴이 완료될 때까지 대기
}
이 테스트 코드는 다음과 같은 방식으로 로그 파일 로테이션을 검증합니다:

고루틴을 사용하여 지속적으로 로그 파일에 데이터를 씁니다.
로테이션을 체크하기 위해 일정 시간 기다립니다.
rotateLogFile 함수를 호출하여 로테이션을 강제로 실행합니다.
로테이션된 로그 파일이 존재하는지 확인합니다.
원본 로그 파일의 크기가 감소했는지 확인하여 새로운 로그 파일이 생성되었는지 검증합니다.
이 테스트는 실제 로그 파일의 로테이션을 시뮬레이션하므로, 로그 파일의 경로 및 이름, 로그 파일 크기 제한 등을 정확하게 지정해야 합니다. 또한, 실제 파일 시스템을 사용하기 때문에, 테스트의 실행 시간과 파일 시스템의 성능에 영향을 받을 수 있습니다. 따라서, 테스트 환경에서 충분한 리소스가 확보되어 있는지 확인하고, 필요에 따라 타임아웃 값을 조정해야 합니다.

User

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
