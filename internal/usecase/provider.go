package usecase

import (
	"context"
	"govision/internal/entity"
)

type FrameBroker interface {
	Start(ctx context.Context)
	Subscribe() <-chan *entity.Frame
	Stop()
}

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

type PicRepository interface {
	Store(frame *entity.Frame) (string, error)
}
