package gocv

import (
	"context"
	"fmt"
	"govision/internal/entity"
	"slices"

	"gocv.io/x/gocv"
)

type Provider struct {
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	faceDetector *FaceDetector
	camera       *Camera
}

func NewProvider(
	faceDetector *FaceDetector,
	camera *Camera,
) *Provider {
	return &Provider{
		faceDetector: faceDetector,
		camera:       camera,
	}
}

func (p *Provider) Start(parent context.Context) error {
	p.ctx, p.cancel = context.WithCancel(parent)
	p.done = make(chan struct{})

	return nil
}

func (p *Provider) Stream() <-chan *entity.Frame {
	ch := make(chan *entity.Frame, 1)

	img := gocv.NewMat()
	go p.camera.stream(p.ctx, &img)

	go func() {
		defer close(ch) // закроем выходной канал при выходе
		defer close(p.done)
		defer img.Close()

		for {
			clone := img.Clone()
			defer clone.Close()

			// кодируем в JPEG/PNG через gocv
			buf, err := gocv.IMEncode(gocv.JPEGFileExt, clone)

			if err != nil {
				fmt.Println("НЕ ПОЛУЧАЕТСЯ ФОРМАТИРОВАТЬ КАРТИНКУ") //todo log this
				continue
			}
			detected := p.faceDetector.Detect(&clone)
			fmt.Println(fmt.Sprintf("Обнаружено %d", detected)) //todo log this

			data := slices.Clone(buf.GetBytes())
			buf.Close()

			eFrame := entity.Frame{
				Data:    data,
				HasFace: detected > 0,
			}

			fmt.Printf("длинна картинки %d\n", len(eFrame.Data)) //todo log this
			select {
			case <-p.ctx.Done():
				return
			case ch <- &eFrame:
				fmt.Println("отправляем картинку") //todo log this
			default:
			}
		}
	}()

	return ch
}

func (p *Provider) Stop() {
	if p.cancel != nil {
		p.cancel() // сигнал остановки
	}
	if p.done != nil {
		<-p.done // ждём, пока Stream реально завершится
	}
}
