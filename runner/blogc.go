package runner

import (
	"fmt"

	"github.com/blogc/go-blogc"
)

type BlogcTask struct {
	Context *blogc.BuildContext
}

func (b *BlogcTask) GetTag() string {
	return "BLOGC"
}

func (b *BlogcTask) GetTarget() blogc.File {
	if b.Context == nil {
		return nil
	}

	return b.Context.OutputFile
}

func (b *BlogcTask) Outdated() bool {
	if b.Context == nil {
		return false
	}

	return b.Context.NeedsBuild()
}

func (b *BlogcTask) Run() error {
	if b.Context == nil {
		return fmt.Errorf("runner: invalid context")
	}

	return b.Context.Build()
}
