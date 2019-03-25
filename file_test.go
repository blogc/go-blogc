package blogc

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestInputFileBytes(t *testing.T) {
	b, err := NewFileBytes([]byte("bola"))
	if err != nil {
		t.Errorf("NewFileBytes failed: %v", err)
	}
	defer b.Close()

	fn := b.Path()
	if !strings.Contains(fn, "blogc_") {
		t.Errorf("Bad tempfile name: %s", fn)
	}

	f, err := ioutil.ReadFile(fn)
	if err != nil {
		t.Errorf("Failed to read tempfile (%s): %v", fn, err)
	}

	if c := bytes.Compare(f, []byte("bola")); c != 0 {
		t.Errorf("Bad tempfile content (%s): %s", fn, f)
	}
}

func TestOutputFileBytes(t *testing.T) {
	b, err := NewFileBytes(nil)
	if err != nil {
		t.Errorf("NewFileBytes failed: %v", err)
	}
	defer b.Close()

	fn := b.Path()
	if !strings.Contains(fn, "blogc_") {
		t.Errorf("Bad tempfile name: %s", fn)
	}

	if err := ioutil.WriteFile(fn, []byte("bola"), 0666); err != nil {
		t.Errorf("Failed to write to tempfile (%s): %v", fn, err)
	}

	f, err := b.Read()
	if err != nil {
		t.Errorf("Failed to read FileBytes content: %v", err)
	}

	if c := bytes.Compare(f, []byte("bola")); c != 0 {
		t.Errorf("Bad FileBytes content (%s): %s", fn, f)
	}
}
