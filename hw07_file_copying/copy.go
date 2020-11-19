package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	time.Sleep(time.Second)
	file, err := os.OpenFile(fromPath, os.O_RDWR, 0666)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()
	inf, _ := file.Stat()
	buf := bufio.NewReaderSize(file, int(inf.Size()))

	if int(offset) > buf.Size() {
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
	bw := bufio.NewWriter(newFile)

	if limit == 0 || limit > int64(buf.Size()) {
		limit = int64(buf.Size())
	}
	_, err = io.CopyN(bw, buf, limit)
	if err != nil {
		return ErrUnsupportedFile
	}
	return nil
}
