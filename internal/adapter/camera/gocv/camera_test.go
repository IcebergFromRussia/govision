package gocv

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gocv.io/x/gocv"
)

type MockVideoCapture struct {
	ok bool
	m  *gocv.Mat
	c  atomic.Uint32
}

type streamTestCase struct {
	name          string
	ticksToSend   int           // сколько тиков отправить в канал
	cancelAfter   time.Duration // через сколько отменить контекст
	expectedCalls int           // ожидаемое количество вызовов readFrame
}

// provider возвращает набор тестовых случаев для проверки stream.
func provider() []streamTestCase {
	return []streamTestCase{
		{
			name:          "zero ticks, cancel immediately",
			ticksToSend:   0,
			cancelAfter:   10 * time.Millisecond,
			expectedCalls: 0,
		},
		{
			name:          "one tick, cancel after processing",
			ticksToSend:   1,
			cancelAfter:   10 * time.Millisecond,
			expectedCalls: 1,
		},
		{
			name:          "multiple ticks, cancel after processing",
			ticksToSend:   5,
			cancelAfter:   50 * time.Millisecond,
			expectedCalls: 5,
		},
	}
}

func (m *MockVideoCapture) Read(mat *gocv.Mat) bool {
	m.m = mat
	m.c.Add(1)
	return m.ok
}

func TestCamera_ReadFrameFailure(t *testing.T) {
	mock := &MockVideoCapture{
		ok: false,
	}

	camera := createCamera(mock, nil)
	img := &gocv.Mat{}

	result := camera.readFrame(img)

	if result {
		t.Fatal("Expected readFrame to return false")
	}
}

func TestCamera_ReadFrameSuccess(t *testing.T) {
	mock := &MockVideoCapture{
		ok: true,
	}

	camera := createCamera(mock, nil)
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

func TestStream_StreamEquivalenceClasses(t *testing.T) {
	for _, tt := range provider() {
		t.Run(tt.name, func(t *testing.T) {
			runStreamTestCase(t, tt)
		})
	}
}

func runStreamTestCase(t *testing.T, tt streamTestCase) {
	mock := &MockVideoCapture{
		ok: true,
	}
	ticker := make(chan time.Time, tt.ticksToSend)
	camera := createCamera(mock, ticker)

	ctx, cancel := context.WithCancel(context.Background())

	img := &gocv.Mat{}

	done := make(chan struct{})

	go func() {
		camera.stream(ctx, img)
		close(done)
	}()
	for i := 0; i < tt.ticksToSend; i++ {
		ticker <- time.Now()
	}
	<-time.After(tt.cancelAfter)
	cancel()

	select {
	case <-done:
		if int(mock.c.Load()) != tt.expectedCalls {
			t.Fatalf("количество вызовов камеры не равно числу тиков (%d / %d), ", tt.expectedCalls, mock.c.Load())
		}
	case <-time.After(time.Second):
		t.Fatal("stream не завершился после отмены контекста")
	}
}

func createCamera(
	mock *MockVideoCapture,
	ticker chan time.Time,
) *Camera {
	if mock == nil {
		mock = &MockVideoCapture{
			ok: true,
		}
	}
	if ticker == nil {
		ticker = make(chan time.Time)
	}

	wg := sync.WaitGroup{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	camera := Camera{
		videoCapture: mock,
		ticker:       ticker,
		wg:           &wg,
		log:          logger,
	}

	return &camera
}
