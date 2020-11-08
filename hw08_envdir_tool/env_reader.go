package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	list := make(map[string]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "readDir failed")
	}

	for _, f := range files {
		name := f.Name()
		if !strings.Contains(name, "=") {
			file, err := os.OpenFile(dir+"/"+name, os.O_RDONLY, 0755)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			list[name] = scanner.Text()

			if err := scanner.Err(); err != nil {
				fmt.Println(name, " err: ", err)
			}
		}
	}
	return list, nil
}
