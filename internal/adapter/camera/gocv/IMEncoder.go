package gocv

import (
	"slices"

	eg "gocv.io/x/gocv"
)

type IMEncoder struct {
	FileExt eg.FileExt
}

func (e *IMEncoder) IMEncode(img eg.Mat) ([]byte, error) {
	buf, err := eg.IMEncode(e.FileExt, img)
	defer buf.Close()

	if err != nil {
		return nil, err
	}

	data := slices.Clone(buf.GetBytes())

	return data, nil
}
