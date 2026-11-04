package gocv

import (
	"errors"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type FaceDetector struct {
	color      color.RGBA
	classifier *gocv.CascadeClassifier
	lastRects  []image.Rectangle
}

func NewFaceDetector(color *color.RGBA) (*FaceDetector, error) {
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load("../gocv/data/haarcascade_frontalface_default.xml") {
		return nil, errors.New("Error reading cascade file: data/haarcascade_frontalface_default.xml")
	}

	return &FaceDetector{
		classifier: &classifier,
		color:      *color,
	}, nil
}

func (fd *FaceDetector) Detect(
	img *gocv.Mat,
) int {
	fd.lastRects = fd.classifier.DetectMultiScale(*img)
	return len(fd.lastRects)
}

func (fd *FaceDetector) Draw(
	img *gocv.Mat,
) {
	for _, r := range fd.lastRects {
		gocv.Rectangle(img, r, fd.color, 3)
	}
}
