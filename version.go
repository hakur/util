package util

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	// Major major verion is always with some breaking changes
	Major int
	// Minor minor version is always with some compatible function changes in major version
	Minor int
	// Bugfix bugfix version is always with some bug fix, no any breaking change or function change
	Bugfix string
}

// ParseVersion support version styles are v0.0.0 and 0.0.0
func ParseVersion(versionStr string) (version *Version, err error) {
	version = new(Version)
	versionArr := strings.Split(strings.Trim(versionStr, "v"), ".")

	if len(versionArr) < 3 {
		return nil, fmt.Errorf("version length less than 3")
	}

	if major, err := strconv.Atoi(versionArr[0]); err != nil {
		return nil, err
	} else {
		version.Major = major
	}

	if minor, err := strconv.Atoi(versionArr[1]); err != nil {
		return nil, err
	} else {
		version.Minor = minor
	}

	version.Bugfix = versionArr[2]

	return
}
