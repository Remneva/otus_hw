package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	file, err := os.OpenFile(fromPath, os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "unsupported file")
	}
	defer file.Close()
	inf, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, "getting stat error")
	}
	buf := bufio.NewReaderSize(file, int(inf.Size()))
	if offset > inf.Size() {
		return errors.New("offset exceeds file size")
	}
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "failed to set offset")
	}
	newFile, err := os.Create(toPath)
	if err != nil {
		return errors.Wrap(err, "failed to trying create file")
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
		log.Println(err)
	}
	return nil
}
