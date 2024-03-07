package util

import (
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

var Manager *GPUManager

type GPUManager struct {
	MaxPipelinesPerGPU int
	CurrentAllocations map[int]int
	gpuCounter         []int32
	lock               sync.Mutex
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
    gm.lock.Lock()
    defer gm.lock.Unlock()

    for i, count := range gm.gpuCounter {
        if count < int32(gm.MaxPipelinesPerGPU) {
            atomic.AddInt32(&gm.gpuCounter[i], 1)
            return i
        }
    }

    logrus.Errorf("All GPUs are fully allocated")
    return -1 
}

