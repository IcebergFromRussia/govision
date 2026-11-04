package usecase

import (
	"context"
	"govision/internal/entity"
)

type CameraProvider interface {
	Start(ctx context.Context) error
	Stream() <-chan *entity.Frame
	Stop()
}

type FrameRepository interface {
	Save(*entity.Frame) error
}

type Window interface {
	ShowWindow(ctx context.Context) error
	Close() error
}

type FaceDetector interface {
	Detect() int
	Draw()
}
