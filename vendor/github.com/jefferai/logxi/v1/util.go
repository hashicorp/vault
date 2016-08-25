package log

import (
	"path/filepath"
	"strings"
)

func expandTabs(s string, tabLen int) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "\t")
	buf := pool.Get()
	defer pool.Put(buf)
	for _, part := range parts {
		buf.WriteString(part)
		buf.WriteString(strings.Repeat(" ", tabLen-len(part)%tabLen))
	}
	return buf.String()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func indexOfNonSpace(s string) int {
	if s == "" {
		return -1
	}
	for i, r := range s {
		if r != ' ' {
			return i
		}
	}
	return -1
}

var inLogxiPath = filepath.Join("mgutz", "logxi", "v"+strings.Split(Version, ".")[0])

func isLogxiCode(filename string) bool {
	// need to see errors in tests
	return strings.HasSuffix(filepath.Dir(filename), inLogxiPath) &&
		!strings.HasSuffix(filename, "_test.go")
}
