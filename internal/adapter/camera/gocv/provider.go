package gocv

import (
	"context"
	"govision/internal/entity"
	"log/slog"
	"time"

	"gocv.io/x/gocv"
)

type Provider struct {
	faceDetector *FaceDetector
	camera       *Camera
	encoder      IMEncoderInterface
	log          *slog.Logger
	ticker       <-chan time.Time
}

type IMEncoderInterface interface {
	IMEncode(img gocv.Mat) ([]byte, error)
}

func NewProvider(
	faceDetector *FaceDetector,
	encoder IMEncoderInterface,
	camera *Camera,
	log *slog.Logger,
) *Provider {
	return &Provider{
		faceDetector: faceDetector,
		camera:       camera,
		log:          log,
		encoder:      encoder,
	}
}

func (p *Provider) Stream(pCtx context.Context) (<-chan *entity.Frame, <-chan struct{}) {
	ch := make(chan *entity.Frame, 1)

	img := gocv.NewMat()
	ctx, cancel := context.WithCancel(pCtx)
	go p.camera.stream(ctx, &img)

	done := make(chan struct{})

	go func() {
		defer close(ch) // закроем выходной канал при выходе
		defer img.Close()
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				cancel()
				return
			case <-p.ticker:
				eFrame := p.getFrame(&img)
				select {
				case ch <- eFrame:
				default:
				}
			}
		}
	}()

	return ch, done
}

func (p *Provider) getFrame(img *gocv.Mat) *entity.Frame {
	clone := img.Clone()
	defer clone.Close()

	data, err := p.encoder.IMEncode(clone)

	if err != nil {
		p.log.Info("НЕ ПОЛУЧАЕТСЯ ФОРМАТИРОВАТЬ КАРТИНКУ")
		return nil
	}
	detected := p.faceDetector.Detect(&clone)

	p.log.Info("Обнаружено ебало", slog.Int("detected", detected))

	return &entity.Frame{
		Data:    data,
		HasFace: detected > 0,
	}
}
