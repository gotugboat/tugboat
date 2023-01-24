package version

import (
	"fmt"
	"runtime"

	"github.com/tonistiigi/go-rosetta"
)

var (
	// Deliberately uninitialized
	version   string
	gitCommit string
)

func arch() string {
	arch := runtime.GOARCH
	if rosetta.Enabled() {
		arch += " (rosetta)"
	}
	return arch
}

func GetFullVersion() string {
	return fmt.Sprintf("%s-%s", GetVersion(), GetCommit())
}

func GetFullVersionWithArch() string {
	return fmt.Sprintf("%s [%s]", GetFullVersion(), arch())
}

func GetCommit() string {
	if gitCommit != "" {
		return gitCommit
	}
	return "unknown"
}

func GetVersion() string {
	if version != "" {
		return version
	}
	return "unknown"
}
