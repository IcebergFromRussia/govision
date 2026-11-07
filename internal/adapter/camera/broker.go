package camera

import (
	"context"
	"fmt"
	"govision/internal/entity"
	"govision/internal/usecase"
	"sync"
)

type Broker struct {
	provider *usecase.CameraProvider
	subs     map[chan *entity.Frame]struct{}
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
	latest   *entity.Frame
}

func NewBroker(provider *usecase.CameraProvider) *Broker {
	return &Broker{
		provider: provider,
		subs:     make(map[chan *entity.Frame]struct{}),
	}
}

func (b *Broker) Start(parent context.Context) {
	b.ctx, b.cancel = context.WithCancel(parent)
	b.done = make(chan struct{})
	(*b.provider).Start(b.ctx)
	stream := (*b.provider).Stream()

	go func() {
		defer close(b.done)
		for {
			fmt.Println("стримим")
			select {
			case <-b.ctx.Done():
				return
			case frame := <-stream:
				b.latest = frame
				b.broadcast(frame)
			}
		}
	}()
}

func (b *Broker) Stop() {
	if b.cancel != nil {
		b.cancel()
	}

	if b.done != nil {
		<-b.done // ждём, пока Stream реально завершится
	}
	(*b.provider).Stop()
}

func (b *Broker) Subscribe() <-chan *entity.Frame {
	ch := make(chan *entity.Frame, 1)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	return ch
}

func (b *Broker) broadcast(frame *entity.Frame) {
	for ch := range b.subs {
		select {
		case ch <- frame:
		default:
		}
	}
}
