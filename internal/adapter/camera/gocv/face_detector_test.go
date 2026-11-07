package gocv

import (
	"image"
	"image/color"
	"sync/atomic"
	"testing"

	"gocv.io/x/gocv"
)

type ClassifierMock struct {
	m          *gocv.Mat
	c          atomic.Uint32
	rectangles []image.Rectangle
}

func (cf *ClassifierMock) DetectMultiScale(img gocv.Mat) []image.Rectangle {
	cf.m = &img
	cf.c.Add(1)

	return cf.rectangles
}

func TestFaceDetector_Detect(t *testing.T) {
	fd := createFaceDetector(1)
	img := &gocv.Mat{}

	mock, ok := fd.classifier.(*ClassifierMock)

	if !ok {
		t.Fatalf("expected *ClassifierMock, got %T", fd.classifier)
	}

	result := fd.Detect(img)

	if result != 1 {
		t.Fatalf("face detect failed detected= %d , expected %d)", result, 1)
	}

	if mock.c.Load() != 1 {
		t.Fatalf("face detect failed detection calls= %d , expected %d)", result, 1)
	}
}

func createFaceDetector(rCount uint8) *FaceDetector {
	mock := createClassifierMock(rCount)
	blue := color.RGBA{0, 0, 255, 0}
	return &FaceDetector{
		color:      blue,
		classifier: mock,
	}
}

func createClassifierMock(rCount uint8) *ClassifierMock {
	rectangles := make([]image.Rectangle, rCount)
	for i := range rectangles {
		rectangles[i] = image.Rectangle{
			Min: image.Point{X: 2, Y: 10},
			Max: image.Point{X: 3, Y: 19},
		}
	}

	return &ClassifierMock{
		rectangles: rectangles,
	}
}
