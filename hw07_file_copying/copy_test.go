package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("File.size less than offset", func(t *testing.T) {

		content := []byte("temporary file's content")
		tmpfile, err := ioutil.TempFile("", "test.")
		if err != nil {
			require.NoError(t, err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write(content); err != nil {
			require.NoError(t, err)
		}

		err = Copy(tmpfile.Name(), "/tmp/", 10000, 100)
		if err != nil {
			log.Println(err)
		}
		require.EqualError(t, err, "offset exceeds file size")
	})

	t.Run("The infinite file unsupported", func(t *testing.T) {
		err := Copy("dev/urandom", "testdata/expected.txt", int64(0), int64(0))
		if err != nil {
			log.Println(err)
		}
		require.EqualError(t, err, "unsupported file: open dev/urandom: no such file or directory")
	})

	t.Run("Filesize less than limit", func(t *testing.T) {
		err := Copy("/dev/null", "testdata/expected.txt", int64(0), int64(0))
		if err != nil {
			log.Println(err)
		}
		require.NoError(t, err)
	})

	t.Run("Success copy", func(t *testing.T) {

		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			require.NoError(t, err)
		}
		fmt.Println("dir: ", dir)
		content := []byte("Hello world")
		tmpfile, err := ioutil.TempFile("", "test.")
		if err != nil {
			log.Println(err)
			require.NoError(t, err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.Write(content); err != nil {
			fmt.Println("write file err: ", err)
			log.Println(err)
			require.NoError(t, err)
		}
		fmt.Println("tmpfile: ", tmpfile.Name())
		err = Copy(tmpfile.Name(), "testdata/expected.txt", 0, 0)
		if err != nil {
			log.Println(err)
		}
		file, err := ioutil.ReadFile("testdata/expected.txt")
		if err != nil {
			fmt.Println("read file err: ", err)
			require.NoError(t, err)
		}
		actual := string(file)
		expected := "Hello world"
		result := reflect.DeepEqual(expected, actual)
		assert.True(t, result)
	})
}
