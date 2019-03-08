package blogc

import (
	"errors"
	"fmt"
	"strings"
)

func Version() (string, error) {
	c := blogcCmd("-v")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

type blogcBase struct {
	OutputFile  string
	Definitions []string
	inputFiles  []string
	listing     bool
}

type Entry struct {
	blogcBase
	InputFile string
}

type Listing struct {
	blogcBase
	InputFiles []string
}

func NewEntry(inputFile string, outputFile string, definitions []string) (*Entry, error) {
	rv := &Entry{
		blogcBase: blogcBase{
			OutputFile:  outputFile,
			Definitions: definitions,
			inputFiles:  []string{inputFile},
			listing:     false,
		},
		InputFile: inputFile,
	}
	if err := rv.validate(); err != nil {
		return nil, err
	}
	return rv, nil
}

func NewListing(inputFiles []string, outputFile string, definitions []string) (*Listing, error) {
	rv := &Listing{
		blogcBase: blogcBase{
			OutputFile:  outputFile,
			Definitions: definitions,
			inputFiles:  inputFiles,
			listing:     true,
		},
		InputFiles: inputFiles,
	}
	if err := rv.validate(); err != nil {
		return nil, err
	}
	return rv, nil
}

func (e *blogcBase) validate() error {
	if e.listing {
		if len(e.inputFiles) == 0 {
			return fmt.Errorf("blogc: at least one input file is required")
		}
	} else {
		if len(e.inputFiles) != 1 || e.inputFiles[0] == "" {
			return fmt.Errorf("blogc: input file is required")
		}
	}

	if e.OutputFile == "" {
		return fmt.Errorf("blogc: output file is required")
	}

	return nil
}

func (e *blogcBase) generateCommand(templateFile string, printVar string) []string {
	rv := []string{}

	if e.listing {
		rv = append(rv, "-l", "-i")
	} else {
		rv = append(rv, e.inputFiles[0])
	}

	for _, v := range e.Definitions {
		rv = append(rv, "-D", v)
	}

	if templateFile != "" {
		rv = append(rv, "-o", e.OutputFile, "-t", templateFile)
	} else if printVar != "" {
		rv = append(rv, "-p", printVar)
	}

	return rv
}

func (e *blogcBase) generateStdin() string {
	rv := ""
	if e.listing {
		for _, v := range e.inputFiles {
			rv += fmt.Sprintf("%s\n", v)
		}
	}
	return rv
}

func (e *blogcBase) Build(templateFile string) error {
	cmdArgs := e.generateCommand(templateFile, "")
	statusCode, _, stderr, err := blogcRun(e.generateStdin(), cmdArgs...)
	if err != nil {
		return err
	}

	if statusCode != 0 {
		return errors.New(strings.TrimSpace(stderr))
	}

	return nil
}

func (e *blogcBase) GetVariable(name string) (string, error) {
	cmdArgs := e.generateCommand("", name)
	statusCode, stdout, stderr, err := blogcRun(e.generateStdin(), cmdArgs...)
	if err != nil {
		return "", err
	}

	if statusCode != 0 {
		return "", errors.New(strings.TrimSpace(stderr))
	}

	// remove the last newline, introduced by blogc itself.
	// we don't want to remove any whitespace that is part of the variable.
	if stdout[len(stdout)-1] == byte('\n') {
		return string(stdout[:len(stdout)-1]), nil
	}

	return stdout, nil
}
