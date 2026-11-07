package storage

import (
	"fmt"
	"govision/internal/entity"
	"os"
	"path/filepath"
	"time"
)

type PicStore struct {
	dirPath   string
	sessionID string
}

func NewPicStore(dirPath string, sessionID string) *PicStore {
	fullPath := filepath.Join(dirPath, sessionID)
	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		fmt.Println("НЕ УДАЛОСЬ СОЗДАТЬ ДИРЕКТОРИЮ")
	}
	return &PicStore{
		fullPath,
		sessionID,
	}
}

func (ps *PicStore) Store(frame *entity.Frame) (string, error) {
	ts := time.Now().UnixNano()
	filename := fmt.Sprintf("%d.jpg", ts)

	fullPath := filepath.Join(ps.dirPath, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return filename, err
	}
	defer file.Close()

	if _, err := file.Write(frame.Data); err != nil {
		return filename, err
	}

	return filename, nil
}
