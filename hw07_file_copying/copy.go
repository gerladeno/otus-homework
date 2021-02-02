package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}
	if fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	var size int64
	switch {
	case limit == 0:
		size = fileInfo.Size() - offset
	case fileInfo.Size()-offset > limit:
		size = limit
	default:
		size = fileInfo.Size() - offset
	}

	to, err := os.Create(toPath)
	if err != nil {
		return err
	}
	_, err = from.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	defer to.Close()
	bar := pb.Full.Start(int(size))
	barReader := bar.NewProxyReader(from)
	_, err = io.CopyN(to, barReader, size)
	bar.Finish()
	return err
}
