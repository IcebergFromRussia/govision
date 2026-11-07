package gocv

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type FaceDetector struct {
	color      color.RGBA
	classifier ClassifierInterface
	lastRects  []image.Rectangle
}

type ClassifierInterface interface {
	DetectMultiScale(img gocv.Mat) []image.Rectangle
}

func NewFaceDetector(
	color *color.RGBA,
	classifier ClassifierInterface,
) *FaceDetector {
	return &FaceDetector{
		classifier: classifier,
		color:      *color,
	}
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
