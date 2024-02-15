package main

import (
  "runtime"
  "gocv.io/x/gocv/cuda"
  )

func gocv_gpu_select()
{
  gpu := 1
  runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	cuda.SetDevice(gpu)
  
  stream := cuda.NewStream()
	defer stream.Close()

	Input_resizeGpuMat, Output_resizeGpuMat := cuda.NewGpuMat(), cuda.NewGpuMat()
	defer Input_resizeGpuMat.Close()
	defer Output_resizeGpuMat.Close()

	Output_resizeMat := gocv.NewMat()
	defer Output_resizeMat.Close()

  ```
  ```
  }
