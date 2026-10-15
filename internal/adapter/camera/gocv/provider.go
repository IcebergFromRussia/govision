package gocv

import (
	"context"
	"fmt"
	"govision/internal/entity"

	"gocv.io/x/gocv"
)

type Provider struct {
	videoCapture *gocv.VideoCapture
	img          *gocv.Mat
	ctx          context.Context
	cancel       context.CancelFunc
	done         chan struct{}
	faceDetector *FaceDetector
}

func NewProvider(
	videoCapture *gocv.VideoCapture,
	img *gocv.Mat,
	faceDetector *FaceDetector,
) *Provider {
	return &Provider{
		videoCapture: videoCapture,
		img:          img,
		faceDetector: faceDetector,
	}
}

func (p *Provider) Start(parent context.Context) error {
	p.ctx, p.cancel = context.WithCancel(parent)
	p.done = make(chan struct{})

	return nil
}

func (p *Provider) Stream() <-chan *entity.Frame {
	ch := make(chan *entity.Frame, 1)

	go func() {
		defer close(ch) // закроем выходной канал при выходе
		defer close(p.done)
		defer p.img.Close()

		for {
			if ok := p.videoCapture.Read(p.img); !ok {
				fmt.Printf("cannot read device \n")
				return
			}
			if p.img.Empty() {
				fmt.Println("пропускаем картинку")
				continue
			}

			clone := p.img.Clone()

			// кодируем в JPEG/PNG через gocv
			buf, err := gocv.IMEncode(gocv.JPEGFileExt, clone)
			if err != nil {
				fmt.Println("НЕ ПОЛУЧАЕТСЯ ФОРМАТИРОВАТЬ КАРТИНКУ")
				continue
			}
			detected := p.faceDetector.Detect()
			//detected := 1
			fmt.Println(fmt.Sprintf("Обнаружено %d", detected))

			eFrame := entity.Frame{
				Data:    buf.GetBytes(),
				HasFace: detected > 0,
			}
			//buf.Close() без закрытого буфера возможна утечка памяти, но так ломаются файлы для сохранения
			//clone.Close()

			fmt.Printf("длинна картинки %d\n", len(eFrame.Data))
			select {
			case <-p.ctx.Done():
				return
			case ch <- &eFrame:
				fmt.Println("отправляем картинку")
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
