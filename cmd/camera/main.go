package main

import (
	"context"
	"fmt"
	"govision/internal/adapter/camera"
	internalGocv "govision/internal/adapter/camera/gocv"
	"govision/internal/adapter/storage"
	"govision/internal/usecase"
	"image/color"
	"log/slog"
	"os"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	img := gocv.NewMat()
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load("../gocv/data/haarcascade_frontalface_default.xml") {
		panic("Error reading cascade file: data/haarcascade_frontalface_default.xml")
	}

	faceDetector := internalGocv.NewFaceDetector(&blue, &classifier)

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	tiker := time.Tick(30 * time.Millisecond)

	c := internalGocv.NewCamera(
		tiker,
		webcam,
		logger,
	)

	var provider usecase.CameraProvider = internalGocv.NewProvider(
		faceDetector,
		c,
	)
	var broker = camera.NewBroker(
		&provider,
	)

	picRepository := storage.NewPicStore(baseDir, sessionID)

	frameSaver := usecase.NewFrameSaverUseCase(
		broker,
		picRepository,
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
