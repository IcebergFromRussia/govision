package usecase

import "context"

type WindowWorkUseCase struct {
	window *Window
}

func NewWindowWorkUseCase(window *Window) *WindowWorkUseCase {
	return &WindowWorkUseCase{
		window: window,
	}
}

func (w *WindowWorkUseCase) Start(ctx context.Context) {
	(*w.window).ShowWindow(ctx)
}
