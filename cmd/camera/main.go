package main

import (
	"context"
	"fmt"
	"govision/internal/adapter/camera"
	internalGocv "govision/internal/adapter/camera/gocv"
	"govision/internal/adapter/storage"
	"govision/internal/usecase"
	"image/color"
	"time"
	"uuid"

	//"github.com/vladimirvivien/go4vl/device"

	"gocv.io/x/gocv"
)

func main() {
	//dev, err := device.Open("/dev/video0", device.WithBufferSize(1))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer dev.Close()

	// set to use a video capture device 0
	deviceID := 0
	blue := color.RGBA{0, 0, 255, 0}
	sessionID := uuid.New().String() // github.com/google/uuid
	baseDir := "./pics"

	img := gocv.NewMat()

	faceDetector, err := internalGocv.NewFaceDetector(&img, &blue)

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	var provider usecase.CameraProvider = internalGocv.NewProvider(
		webcam,
		&img,
		faceDetector,
	)
	var broker usecase.FrameBroker = camera.NewBroker(
		&provider,
	)

	var picRepository usecase.PicRepository = storage.NewPicStore(baseDir, sessionID)

	frameSaver := usecase.NewFrameSaverUseCase(
		&broker,
		&picRepository,
	)
	ctx := context.Background()

	broker.Start(ctx)

	w := gocv.NewWindow("Hello")

	window := internalGocv.NewWindow(w, &img)

	go frameSaver.SaveFrame(ctx)

	window.ShowWindow(ctx)
	time.Sleep(5 * time.Second)

	broker.Stop()

}
