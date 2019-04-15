package blogc

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
)

func RequiredVersion(v string) error {
	actualVersion, err := version.NewVersion(Version)
	if err != nil {
		return err
	}

	reqVersion, err := version.NewVersion(v)
	if err != nil {
		return err
	}

	if reqVersion.Compare(actualVersion) > 0 {
		return fmt.Errorf("blogc: version %q or greater required, got %q", v, Version)
	}

	return nil
}

type BuildContext struct {
	Listing          bool
	GlobalVariables  []string
	InputFiles       []File
	ListingEntryFile File
	OutputFile       File
	TemplateFile     File
}

func (e *BuildContext) NeedsBuild() bool {
	st, err := os.Stat(e.OutputFile.Path())
	if err != nil {
		return true
	}
	mtime := st.ModTime()

	files := e.InputFiles
	if e.TemplateFile != nil {
		files = append(files, e.TemplateFile)
	}
	if e.Listing && e.ListingEntryFile != nil {
		files = append(files, e.ListingEntryFile)
	}

	for _, f := range files {
		st, err := os.Stat(f.Path())
		if err != nil {
			// file not found. recomend a new build, so blogc can generate
			// useful error message
			return true
		}

		if mtime.Before(st.ModTime()) {
			return true
		}
	}

	return false
}

func (e *BuildContext) generateCommand(printVar string) []string {
	rv := []string{}

	if e.Listing {
		rv = append(rv, "-l", "-i")
		if e.ListingEntryFile != nil {
			rv = append(rv, "-e", e.ListingEntryFile.Path())
		}
	} else {
		rv = append(rv, e.InputFiles[0].Path())
	}

	for _, v := range e.GlobalVariables {
		rv = append(rv, "-D", v)
	}

	if printVar != "" {
		rv = append(rv, "-p", printVar)
	} else if e.OutputFile != nil && e.TemplateFile != nil {
		rv = append(rv, "-o", e.OutputFile.Path(), "-t", e.TemplateFile.Path())
	}

	return rv
}

func (e *BuildContext) generateStdin() string {
	rv := ""
	if e.Listing {
		for _, v := range e.InputFiles {
			rv += fmt.Sprintf("%s\n", v.Path())
		}
	}
	return rv
}

func (e *BuildContext) validateInputFiles() error {
	if e.Listing {
		if len(e.InputFiles) == 0 {
			return fmt.Errorf("blogc: at least one input file is required")
		}
	} else {
		if len(e.InputFiles) != 1 {
			return fmt.Errorf("blogc: one input file is required")
		}
		if e.ListingEntryFile != nil {
			return fmt.Errorf("blogc: listing entry is only supported by listing mode")
		}
	}
	return nil
}

func (e *BuildContext) Build() error {
	if err := e.validateInputFiles(); err != nil {
		return err
	}

	if e.OutputFile == nil {
		return fmt.Errorf("blogc: output file is required")
	}

	if e.TemplateFile == nil {
		return fmt.Errorf("blogc: template file is required")
	}

	cmdArgs := e.generateCommand("")
	statusCode, _, stderr, err := blogcRun(e.generateStdin(), cmdArgs...)
	if err != nil {
		return err
	}

	if statusCode != 0 {
		return errors.New(strings.TrimSpace(stderr))
	}

	return nil
}

func (e *BuildContext) GetEvaluatedVariable(name string) (string, bool, error) {
	if name == "" {
		return "", false, fmt.Errorf("blogc: variable name is required")
	}

	if err := e.validateInputFiles(); err != nil {
		return "", false, err
	}

	cmdArgs := e.generateCommand(name)
	statusCode, stdout, stderr, err := blogcRun(e.generateStdin(), cmdArgs...)
	if err != nil {
		return "", false, err
	}

	if statusCode != 0 {
		if statusCode == 78 { // EX_CONFIG, as of blogc-0.15.2, that is our minimum version
			return "", false, nil
		}
		return "", false, errors.New(strings.TrimSpace(stderr))
	}

	// remove the last newline, introduced by blogc itself.
	// we don't want to remove any whitespace that is part of the variable.
	if stdout[len(stdout)-1] == byte('\n') {
		return string(stdout[:len(stdout)-1]), true, nil
	}

	return stdout, true, nil
}
