package main

import (
	"fmt"
	"os"
)

func main() {
	path := os.Args[1]
	command := os.Args[2:]
	env, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Stderr.WriteString(err.Error())
	}
	RunCmd(command, env)
}
