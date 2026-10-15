package gocv

import (
	"context"
	"fmt"

	"gocv.io/x/gocv"
)

type Window struct {
	window *gocv.Window
	img    *gocv.Mat
}

func NewWindow(window *gocv.Window, img *gocv.Mat) *Window {
	return &Window{
		window: window,
		img:    img,
	}
}

func (w *Window) ShowWindow(ctx context.Context) error {
	for {
		if w.img.Empty() {
			continue
		}

		// show the image in the window, and wait 1 millisecond
		err := w.window.IMShow(*w.img)
		if err != nil {
			fmt.Println(err)
			//return err
		}

		w.window.WaitKey(1)
		select {
		case <-ctx.Done():
			return w.Close()
		default:
		}
	}

}

func (w *Window) Close() error {
	err := w.window.Close()
	return err
}
