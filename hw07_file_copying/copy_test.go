package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("File.size less than offset", func(t *testing.T) {

		content := []byte("temporary file's content")
		tmpfile, err := ioutil.TempFile("", "test.")
		if err != nil {
			log.Println(err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write(content); err != nil {
			log.Println(err)
		}

		err = Copy(tmpfile.Name(), "/tmp/", 10000, 100)
		if err != nil {
			log.Println(err)
		}
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())

	})

	t.Run("The infinite file unsupported", func(t *testing.T) {
		err := Copy("dev/urandom", "testdata/expected.txt", int64(0), int64(0))
		if err != nil {
			log.Println(err)
		}
		require.EqualError(t, err, ErrUnsupportedFile.Error())

	})
}
