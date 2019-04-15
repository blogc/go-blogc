package blogc

import (
	"io"
	"io/ioutil"
	"os/exec"
	"syscall"
)

func blogcCmd(args ...string) *exec.Cmd {
	return exec.Command(blogcBin, args...)
}

func blogcRun(stdinStr string, args ...string) (int, string, string, error) {
	cmd := blogcCmd(args...)

	if stdinStr != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return 0, "", "", err
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, stdinStr)
		}()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", "", err
	}

	if err := cmd.Start(); err != nil {
		return 0, "", "", err
	}

	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return 0, "", "", err
	}

	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return 0, "", "", err
	}

	stdoutStr := string(stdoutBytes)
	stderrStr := string(stderrBytes)

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), stdoutStr, stderrStr, nil
			}
		}
		return 0, "", "", err
	}

	return 0, stdoutStr, stderrStr, nil
}
