package gocv

import (
	"errors"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type FaceDetector struct {
	img        *gocv.Mat
	color      color.RGBA
	classifier *gocv.CascadeClassifier
	lastRects  []image.Rectangle
}

func NewFaceDetector(img *gocv.Mat, color *color.RGBA) (*FaceDetector, error) {
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load("../gocv/data/haarcascade_frontalface_default.xml") {
		return nil, errors.New("Error reading cascade file: data/haarcascade_frontalface_default.xml")
	}

	return &FaceDetector{
		img:        img,
		classifier: &classifier,
		color:      *color,
	}, nil
}

func (fd *FaceDetector) Detect() int {
	fd.lastRects = fd.classifier.DetectMultiScale(*fd.img)
	return len(fd.lastRects)
}

func (fd *FaceDetector) Draw() {
	for _, r := range fd.lastRects {
		gocv.Rectangle(fd.img, r, fd.color, 3)
	}
}
