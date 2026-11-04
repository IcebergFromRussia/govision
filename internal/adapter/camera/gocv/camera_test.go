package gocv

import (
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"gocv.io/x/gocv"
)

type MockVideoCapture struct {
	ok bool
	m  *gocv.Mat
}

func (m *MockVideoCapture) Read(mat *gocv.Mat) bool {
	m.m = mat
	return m.ok
}

func TestCamera_ReadFrameSuccess(t *testing.T) {
	mock := &MockVideoCapture{
		ok: true,
	}
	ticker := make(chan time.Time)
	wg := sync.WaitGroup{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	camera := Camera{
		videoCapture: mock,
		ticker:       ticker,
		wg:           &wg,
		log:          logger,
	}
	img := &gocv.Mat{}

	result := camera.readFrame(img)

	if !result {
		t.Fatal("Expected readFrame to return true")
	}

	if mock.m == nil {
		t.Fatal("MockVideoCapture не получил МАТ")
	}

	if mock.m != img {
		t.Fatal("MockVideoCapture получил не тот экземпляр МАТ")
	}
}
