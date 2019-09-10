package runner

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/blogc/go-blogc"
)

type Task interface {
	GetTag() string
	GetTarget() blogc.File
	Outdated() bool
	Run() error
}

func Runner(tasks []Task) bool {
	for _, task := range tasks {
		if task == nil {
			continue
		}

		if !task.Outdated() {
			continue
		}

		target := ""
		if t := task.GetTarget(); t != nil {
			target = t.Path()
		}

		fmt.Fprintf(os.Stderr, "  %-8s %s\n", strings.ToUpper(task.GetTag()), target)

		if err := task.Run(); err != nil {
			prefix := "go-blogc"
			if len(os.Args) > 0 {
				prefix = path.Base(os.Args[0])
			}

			fmt.Fprintf(os.Stderr, "%s: error: %q: %s", prefix, target, err.Error())

			return false
		}
	}

	return true
}
