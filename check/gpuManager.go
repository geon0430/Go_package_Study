package util

import (
	"os/exec"
	"strings"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var Manager *GPUManager

type GPUManager struct {
	MaxPipelinesPerGPU int
	CurrentAllocations map[int]int
	gpuCounter         []int32
}

func InitGPUManagerWithSystemGPUs(maxPipelinesPerGPU int) (int, error) {
	gpuCount, err := getSystemGPUsCount()
	if err != nil {
		logrus.Errorf("Failed to get system GPUs count: %v", err)
		return 0, err
	}

	logrus.Infof("%d GPUs in the system", gpuCount)

	Manager = NewGPUManager(gpuCount, maxPipelinesPerGPU)

	return gpuCount, nil
}

func getSystemGPUsCount() (int, error) {
	cmd := exec.Command("nvidia-smi", "-L")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	gpuCount := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
	return gpuCount, nil
}

func NewGPUManager(gpuCount, maxPipelinesPerGPU int) *GPUManager {
	return &GPUManager{
		MaxPipelinesPerGPU: maxPipelinesPerGPU,
		CurrentAllocations: make(map[int]int, gpuCount),
		gpuCounter:         make([]int32, gpuCount),
	}
}

func (gm *GPUManager) AssignGPUNumber() int {
	for {
		for i := 0; i < len(gm.gpuCounter); i++ {
			if atomic.LoadInt32(&gm.gpuCounter[i]) < int32(gm.MaxPipelinesPerGPU) {
				atomic.AddInt32(&gm.gpuCounter[i], 1)
				return i
			}
		}
	}
}

