// +build windows

package manager

import (
	"fmt"
	"strings"
)

func prepCommand(command string) ([]string, error) {
	switch len(strings.Fields(command)) {
	case 0:
		return []string{}, nil
	case 1:
		return []string{command}, nil
	}
	return []string{}, fmt.Errorf("only single commands supported on windows")
}
