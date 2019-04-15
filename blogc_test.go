package blogc

import (
	"regexp"
	"testing"
)

var (
	v  = regexp.MustCompile(`^[0-9a-f-.]+(-dirty)?$`)
	pv = regexp.MustCompile(`^blogc [0-9a-f-.]+(-dirty)?$`)
)

func TestVersion(t *testing.T) {
	if !v.MatchString(Version) {
		t.Errorf("Version failed: version not found: %s", Version)
	}
}

func TestPackageVersion(t *testing.T) {
	if !pv.MatchString(PackageVersion) {
		t.Errorf("PackageVersion failed: version not found: %s", PackageVersion)
	}
}
