package usecase

import (
	"context"
	"fmt"
)

type FrameSaverUseCase struct {
	broker        *FrameBroker
	picRepository *PicRepository
}

func NewFrameSaverUseCase(
	broker *FrameBroker,
	picRepository *PicRepository,
) *FrameSaverUseCase {
	return &FrameSaverUseCase{
		broker,
		picRepository,
	}
}

func (f *FrameSaverUseCase) SaveFrame(ctx context.Context) {
	frameChan := (*f.broker).Subscribe()

	fmt.Println("ожидаем каритинку")

	for {
		select {
		case <-ctx.Done():
			return
		case frame := <-frameChan:
			if !frame.HasFace {
				continue
			}
			name, err := (*f.picRepository).Store(frame)
			if err != nil {
				fmt.Println("ошибка во время сохранения файла")

				fmt.Println(err)
				continue
			}
			fmt.Println("Сохранено изображение" + name)
		}
	}
}
