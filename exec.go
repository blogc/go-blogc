package blogc

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/hashicorp/go-version"
)

const (
	blogcRequiredVersion = "0.15.1"
)

var (
	blogcBin = "blogc"
)

func init() {

	if bin, ok := os.LookupEnv("BLOGC"); ok {
		blogcBin = bin
	}

	// check if binary exists
	var err error
	blogcBin, err = exec.LookPath(blogcBin)
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok {
			if execErr.Err == exec.ErrNotFound {
				panic("blogc: failed to find \"blogc\" binary in PATH, please install from https://blogc.rgm.io/, or set BLOGC environment variable")
			}
		}
		panic(err)
	}

	// check binary version
	vStr, err := Version()
	if err != nil {
		panic(err)
	}
	pieces := strings.Split(vStr, " ")
	if len(pieces) != 2 {
		panic(fmt.Sprintf("blogc: malformed version reported by %q binary: %s", blogcBin, vStr))
	}

	v, err := version.NewVersion(pieces[1])
	if err != nil {
		panic(err)
	}

	r, err := version.NewVersion(blogcRequiredVersion)
	if err != nil {
		panic(err)
	}

	if c := r.Compare(v); c > 0 {
		panic(fmt.Sprintf("blogc: version %q or greater required, got %q", blogcRequiredVersion, pieces[1]))
	}
}

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
