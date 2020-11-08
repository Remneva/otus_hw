package main

import (
	"fmt"
	"github.com/dchest/safefile"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("чтение файла из директории", func(t *testing.T) {

		err := os.Mkdir("/tmp/testdir", 0755)
		defer os.RemoveAll("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}

		tmpfile, err := safefile.Create("/tmp/testdir/S.txt", 0755)
		if err != nil {
			fmt.Println(err)
		}
		defer tmpfile.Close()

		content := []byte("temporary")
		if _, err := tmpfile.Write(content); err != nil {
			log.Fatal(err)
		}

		if err := tmpfile.Commit(); err != nil {
			log.Fatal(err)
		}

		var actual Environment
		actual, err = ReadDir("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}
		for k, v := range actual {
			actual[k] = v
		}
		var expected Environment
		expected = map[string]string{
			"S.txt": "temporary",
		}
		result := reflect.DeepEqual(expected, actual)
		assert.True(t, result)
	})

	t.Run("чтение невалидного файла", func(t *testing.T) {

		err := os.Mkdir("/tmp/testdir", 0755)
		defer os.RemoveAll("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}
		tmpfile, err := safefile.Create("/tmp/testdir/S=q.txt", 0755)
		if err != nil {
			fmt.Println(err)
		}
		defer tmpfile.Close()

		content := []byte("temporary")
		if _, err := tmpfile.Write(content); err != nil {
			log.Fatal(err)
		}

		if err := tmpfile.Commit(); err != nil {
			log.Fatal(err)
		}

		var actual Environment
		actual, err = ReadDir("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}
		for k, v := range actual {
			actual[k] = v
		}
		var expected Environment
		expected = map[string]string{}

		result := reflect.DeepEqual(expected, actual)
		assert.True(t, result)
	})

	t.Run("чтение файла состоящего из нескольких строк", func(t *testing.T) {

		err := os.Mkdir("/tmp/testdir", 0755)
		defer os.RemoveAll("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}
		tmpfile, err := safefile.Create("/tmp/testdir/file.txt", 0755)
		if err != nil {
			fmt.Println(err)
		}
		defer tmpfile.Close()

		content1 := []byte("the first line\nthe second line\nthe third line")
		if _, err := tmpfile.Write(content1); err != nil {
			log.Fatal(err)
		}

		if err := tmpfile.Commit(); err != nil {
			log.Fatal(err)
		}

		var actual Environment
		actual, err = ReadDir("/tmp/testdir")
		if err != nil {
			fmt.Println(err)
		}
		for k, v := range actual {
			actual[k] = v
		}
		var expected Environment
		expected = map[string]string{
			"file.txt": "the first line",
		}
		result := reflect.DeepEqual(expected, actual)
		assert.True(t, result)
	})

}
