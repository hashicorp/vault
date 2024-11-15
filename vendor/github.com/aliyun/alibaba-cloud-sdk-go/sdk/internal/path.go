package internal

import (
	"os"
	"runtime"
)

// GetHomePath return home directory according to the system.
// if the environmental virables does not exist, will return empty string
func GetHomePath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	return os.Getenv("HOME")
}
