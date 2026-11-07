package usecase

import (
	"context"
	"fmt"
	"govision/internal/entity"
)

type FrameSaverUseCase struct {
	broker        FrameBroker
	picRepository PicRepositoryInterface
}

type PicRepositoryInterface interface {
	Store(frame *entity.Frame) (string, error)
}

type FrameBroker interface {
	Start(ctx context.Context)
	Subscribe() <-chan *entity.Frame
	Stop()
}

func NewFrameSaverUseCase(
	broker FrameBroker,
	picRepository PicRepositoryInterface,
) *FrameSaverUseCase {
	return &FrameSaverUseCase{
		broker,
		picRepository,
	}
}

func (f *FrameSaverUseCase) SaveFrame(ctx context.Context) {
	frameChan := f.broker.Subscribe()

	fmt.Println("ожидаем каритинку")

	for {
		select {
		case <-ctx.Done():
			return
		case frame := <-frameChan:
			if !frame.HasFace {
				continue
			}
			name, err := f.picRepository.Store(frame)
			if err != nil {
				fmt.Println("ошибка во время сохранения файла")

				fmt.Println(err)
				continue
			}
			fmt.Println("Сохранено изображение" + name)
		}
	}
}
