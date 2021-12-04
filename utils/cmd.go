package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

func getExitCode(err error) (int, error) {
	exitCode := 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if procExit := exiterr.Sys().(syscall.WaitStatus); ok {
			return procExit.ExitStatus(), nil
		}
	}
	return exitCode, fmt.Errorf("获取状态信息失败")
}

func processExitCode(err error) (exitCode int) {
	if err != nil {
		var exiterr error
		if exitCode, exiterr = getExitCode(err); exiterr != nil {
			exitCode = 127
		}
	}
	return
}

func CmdOutExitCodeBytes(name string, arg ...string) ([]byte, int, error) {
	cmd := exec.Command(name, arg...)
	exitCode := 0
	out, err := cmd.CombinedOutput()
	exitCode = processExitCode(err)
	return out, exitCode, err
}

func CmdOutBytes(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
