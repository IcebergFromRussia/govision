package usecase

import (
	"context"
	"govision/internal/entity"
)

type CameraProvider interface {
	Stream(pCtx context.Context) (<-chan *entity.Frame, <-chan struct{})
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
