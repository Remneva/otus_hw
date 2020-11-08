package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for i := 0; i < len(env); i++ {
		for k, v := range env {
			v = string(bytes.ReplaceAll([]byte(v), []byte("\x00"), []byte("\n")))
			v = strings.TrimRight(v, "\t")
			os.Setenv(k, v)
			if v == "" {
				os.Unsetenv(k)
			}
			i++
		}
	}

	cm := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	cm.Stdout = os.Stdout
	if err := cm.Run(); err != nil {
		fmt.Println(err)
		os.Stderr.WriteString(err.Error())
		return 1
	}

	return 0
}
