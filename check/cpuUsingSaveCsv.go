package model

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)
func check()
{
  csvFile, err := os.Create("cpu_usage.csv")
	if err != nil {
		logger.Fatalf("Failed to create CSV file: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// CSV 헤더 작성
	writer.Write([]string{"Time", "CPU Usage (%)"})

}

func getCPUUsage() (string, error) {
	cmd := exec.Command("mpstat", "1", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "all") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				idleStr := fields[len(fields)-1]
				idle, err := strconv.ParseFloat(idleStr, 64)
				if err != nil {
					return "", fmt.Errorf("Error parsing idle value: %v", err)
				}
				usage := 100 - idle
				return fmt.Sprintf("%.2f", usage), nil
			}
		}
	}

	return "", fmt.Errorf("CPU usage not found")
}
