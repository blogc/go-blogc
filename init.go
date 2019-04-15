package blogc

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	blogcRequiredVersion = "0.16.0"
)

// blogc version number as reported by the blogc binary.
var Version string

// go-blogc version.
var LibraryVersion string

// blogc package version as reported by the blogc binary. Output of `blogc -v`.
var PackageVersion string

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
				fmt.Println("blogc: failed to find \"blogc\" binary in PATH, please install from https://blogc.rgm.io/, or set BLOGC environment variable")
				os.Exit(1)
			}
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// detect versions
	c := blogcCmd("-v")
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	PackageVersion = strings.TrimSpace(string(out))
	pieces := strings.Split(PackageVersion, " ")
	if len(pieces) != 2 || pieces[0] != "blogc" {
		fmt.Printf("blogc: malformed version reported by %q binary: %s\n", blogcBin, PackageVersion)
		os.Exit(1)
	}
	Version = pieces[1]

	// check binary version
	if err := RequiredVersion(blogcRequiredVersion); err != nil {
		panic(err)
	}
}
