package blogc

import (
	"io/ioutil"
	"os"
)

type File interface {
	Path() string
}

type FilePath string

func (f FilePath) Path() string {
	return string(f)
}

type FileBytes struct {
	path string
}

func NewFileBytes(in []byte) (*FileBytes, error) {
	f, err := ioutil.TempFile("", "blogc_")
	if err != nil {
		return nil, err
	}

	filename := f.Name()

	if c := len(in); c > 0 {
		if _, err := f.Write(in); err != nil {
			os.Remove(filename)
			return nil, err
		}
	}

	if err := f.Close(); err != nil {
		os.Remove(filename)
		return nil, err
	}

	return &FileBytes{path: filename}, nil
}

func (f *FileBytes) Path() string {
	return f.path
}

func (f *FileBytes) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *FileBytes) Close() {
	os.Remove(f.path)
}
