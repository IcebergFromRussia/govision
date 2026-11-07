package gocv

import (
	"context"
	"time"

	"log/slog"

	"gocv.io/x/gocv"
)

type Camera struct {
	videoCapture VideoCaptureInterface
	ticker       <-chan time.Time
	log          *slog.Logger
}

type VideoCaptureInterface interface {
	Read(m *gocv.Mat) bool
}

func NewCamera(
	ticker <-chan time.Time,
	videoCapture VideoCaptureInterface,
	log *slog.Logger,
) *Camera {
	return &Camera{
		ticker:       ticker,
		videoCapture: videoCapture,
		log:          log,
	}
}

func (c *Camera) readFrame(
	img *gocv.Mat,
) bool {
	ok := c.videoCapture.Read(img)
	if !ok {
		c.log.Info("Camera: failed to read frame")
	}

	return ok
}

func (c *Camera) stream(
	ctx context.Context,
	img *gocv.Mat,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.ticker:
			c.readFrame(img)
		}
	}
}
