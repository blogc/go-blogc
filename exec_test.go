package blogc

import (
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	cmd := blogcCmd("-v")
	s, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Cmd failed: %v", err)
	}

	if !strings.HasPrefix(string(s), "blogc ") {
		t.Errorf("Cmd failed: version not found")
	}
}
