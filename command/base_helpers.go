package command

import (
	"fmt"
	"strings"
)

var ErrMissingPath = fmt.Errorf("Missing PATH!")

// extractPath extracts the path and list of arguments from the args. If there
// are no extra arguments, the remaining args will be nil.
func extractPath(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, ErrMissingPath
	}

	// Path is always the first argument after all flags
	path := args[0]

	// Strip leading and trailing slashes
	for len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Trim any leading/trailing whitespace
	path = strings.TrimSpace(path)

	// Verify we have a path
	if path == "" {
		return "", nil, ErrMissingPath
	}

	// Splice remaining args
	var remaining []string
	if len(args) > 1 {
		remaining = args[1:]
	}

	return path, remaining, nil
}
