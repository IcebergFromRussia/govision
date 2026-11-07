package go4vl

import (
	"context"
	"govision/internal/entity"

	"github.com/vladimirvivien/go4vl/device"
)

type Provider struct {
	device *device.Device
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func NewProvider(device *device.Device) *Provider {
	return &Provider{
		device: device,
	}
}

func (p *Provider) Start(parent context.Context) error {
	p.ctx, p.cancel = context.WithCancel(parent)
	p.done = make(chan struct{})

	if err := p.device.Start(p.ctx); err != nil {
		p.cancel()
		return err
	}
	return nil
}

func (p *Provider) Stream() <-chan *entity.Frame {
	ch := make(chan *entity.Frame, 1)

	go func() {
		defer close(ch) // закроем выходной канал при выходе
		defer close(p.done)

		deviceChanel := p.device.GetFrames()
		for {
			select {
			case <-p.ctx.Done():
				return
			case dFrame := <-deviceChanel:
				eFrame := entity.Frame{
					Data: dFrame.Data,
				}
				select {
				case ch <- &eFrame:
				default:
				}
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
