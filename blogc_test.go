package blogc

import (
	"regexp"
	"testing"
)

var v = regexp.MustCompile(`^blogc [0-9a-f-.]+(-dirty)?$`)

func TestVersion(t *testing.T) {
	s, err := Version()
	if err != nil {
		t.Errorf("Version failed: %v", err)
	}

	if !v.MatchString(s) {
		t.Errorf("Version failed: version not found: %s", s)
	}
}
