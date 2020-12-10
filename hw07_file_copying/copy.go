package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	file, err := os.OpenFile(fromPath, os.O_RDWR, 0666)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()
	inf, _ := file.Stat()
	buf := bufio.NewReaderSize(file, int(inf.Size()))
	if offset > inf.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)

	if err != nil {
		return fmt.Errorf("failed to set offset: %w", err)
	}
	newFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to trying create file: %w", err)
	}
	defer newFile.Close()

	if limit == 0 || limit > int64(buf.Size()) {
		limit = int64(buf.Size())
	}
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(buf)
	_, err = io.CopyN(newFile, barReader, limit)
	bar.Finish()

	if err != nil {
		return ErrUnsupportedFile
	}
	return nil
}
