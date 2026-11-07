package gocv

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log/slog"
	"testing"
	"time"

	"gocv.io/x/gocv"
)

type IMEncoderMock struct {
	len int
}

func (ie *IMEncoderMock) IMEncode(img gocv.Mat) ([]byte, error) {
	i := image.NewRGBA(image.Rect(0, 0, 1, 1))
	i.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, i, &jpeg.Options{Quality: 90})
	bytes := buf.Bytes()
	ie.len = len(bytes)
	return bytes, nil
}

func TestProvider_StreamSuccessful(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ticker := make(chan time.Time)
	p := createProvider(ticker)

	ieMock := p.encoder.(*IMEncoderMock)

	ch, done := p.Stream(ctx)
	ticker <- time.Now()

	<-time.After(10 * time.Millisecond)
	cancel()

	select {
	case <-done:
		frame := <-ch
		if !frame.HasFace {
			t.Fatalf("лица не обнаружены")
		}
		if len(frame.Data) != ieMock.len {
			t.Fatalf("изображдение неверно декодировано")
		}
	case <-time.After(time.Second):
		t.Fatal("stream не завершился после отмены контекста")
	}
}

func createProvider(
	ticker chan time.Time,
) *Provider {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	camera := createCamera(nil, nil)
	ie := &IMEncoderMock{}

	if ticker == nil {
		ticker = make(chan time.Time)
	}

	p := Provider{
		faceDetector: createFaceDetector(1),
		camera:       camera,
		encoder:      ie,
		log:          logger,
		ticker:       ticker,
	}

	return &p
}
