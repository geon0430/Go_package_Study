package main

import (
	"image"
	"log"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/cuda"
)

var (
	ImgData        = "/go_vms/src/input_images/0001.jpg"
	in_height  = 1080
	in_weight  = 1920

)

func main() {
	img := gocv.IMRead(ImgData, gocv.IMReadColor)
	if img.Empty() {
		log.Println("Error reading image from file")
		return
	}
	defer img.Close()

  DownloadedMat := gocv.NewMat()
  defer DownloadedMat.Close()
  
	gpuMat, resizedGpuMat := cuda.NewGpuMat(), cuda.NewGpuMat()
  defer gpuMat.Close()
  defer resizedGpuMat.Close()

	stream := cuda.NewStream()
	defer stream.Close()


	gpuMat.Upload(img)

	cuda.ResizeWithStream(gpuMat, &resizedGpuMat, image.Pt(in_weight, in_height), 0, 0, cuda.InterpolationLinear, stream)

  resizedGpuMat.DownloadWithStream(&DownloadedMat, stream)

	gocv.IMWrite("output_image.jpg", DownloadedMat)

	log.Println("Success")
}
