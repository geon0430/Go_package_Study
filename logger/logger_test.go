func rotateLogFileOnce(logger *logrus.Logger, file *os.File, basePath string, maxSizeMB int64) {
	// 로그 파일 크기를 가져옵니다.
	fileInfo, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Fatalf("Failed to stat log file: %v", err)
	}

	// 파일 크기가 최대 크기를 초과하는 경우 로테이션을 수행합니다.
	if fileInfo.Size() > maxSizeMB*1024*1024 {
		file.Close()

		for i := maxIndex - 1; i >= 0; i-- {
			oldPath := incrementLogFileName(basePath, i)
			newPath := incrementLogFileName(basePath, i+1)
			if _, err := os.Stat(oldPath); err == nil {
				if i == (maxIndex - 1) {
					os.Remove(oldPath)
				} else {
					os.Rename(oldPath, newPath)
				}
			}
		}

		os.Rename(basePath, incrementLogFileName(basePath, 1))

		newFile, err := os.OpenFile(basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatalf("Failed to open new log file: %v", err)
		}

		logger.SetOutput(newFile)
		file = newFile
	}
}
