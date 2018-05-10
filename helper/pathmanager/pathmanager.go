package pathmanager

import (
	"strings"
	"sync"

	iradix "github.com/hashicorp/go-immutable-radix"
)

// PathManager is a prefix searchable index of paths
type PathManager struct {
	l     sync.RWMutex
	paths *iradix.Tree
}

// New creates a new path manager
func New() *PathManager {
	return &PathManager{
		paths: iradix.New(),
	}
}

// AddPaths adds path to the paths list
func (p *PathManager) AddPaths(paths []string) {
	p.l.Lock()
	defer p.l.Unlock()

	txn := p.paths.Txn()
	for _, prefix := range paths {
		if len(prefix) == 0 {
			continue
		}

		var exception bool
		if strings.HasPrefix(prefix, "!") {
			prefix = strings.TrimPrefix(prefix, "!")
			exception = true
		}

		// We trim any trailing *, but we don't touch whether it is a trailing
		// slash or not since we want to be able to ignore prefixes that fully
		// specify a file
		txn.Insert([]byte(strings.TrimSuffix(prefix, "*")), exception)
	}
	p.paths = txn.Commit()
}

// RemovePaths removes paths from the paths list
func (p *PathManager) RemovePaths(paths []string) {
	p.l.Lock()
	defer p.l.Unlock()

	txn := p.paths.Txn()
	for _, prefix := range paths {
		if len(prefix) == 0 {
			continue
		}

		// Exceptions aren't stored with the leading ! so strip it
		if strings.HasPrefix(prefix, "!") {
			prefix = strings.TrimPrefix(prefix, "!")
		}

		// We trim any trailing *, but we don't touch whether it is a trailing
		// slash or not since we want to be able to ignore prefixes that fully
		// specify a file
		txn.Delete([]byte(strings.TrimSuffix(prefix, "*")))
	}
	p.paths = txn.Commit()
}

// RemovePathPrefix removes all paths with the given prefix
func (p *PathManager) RemovePathPrefix(prefix string) {
	p.l.Lock()
	defer p.l.Unlock()

	// We trim any trailing *, but we don't touch whether it is a trailing
	// slash or not since we want to be able to ignore prefixes that fully
	// specify a file
	p.paths, _ = p.paths.DeletePrefix([]byte(strings.TrimSuffix(prefix, "*")))
}

// Len returns the number of paths
func (p *PathManager) Len() int {
	return p.paths.Len()
}

// Paths returns the path list
func (p *PathManager) Paths() []string {
	p.l.RLock()
	defer p.l.RUnlock()

	paths := make([]string, 0, p.paths.Len())
	walkFn := func(k []byte, v interface{}) bool {
		paths = append(paths, string(k))
		return false
	}
	p.paths.Root().Walk(walkFn)
	return paths
}

// HasPath returns if the prefix for the path exists regardless if it is a path
// (ending with /) or a prefix for a leaf node
func (p *PathManager) HasPath(path string) bool {
	p.l.RLock()
	defer p.l.RUnlock()

	if _, exceptionRaw, ok := p.paths.Root().LongestPrefix([]byte(path)); ok {
		var exception bool
		if exceptionRaw != nil {
			exception = exceptionRaw.(bool)
		}
		return !exception
	}
	return false
}

// HasExactPath returns if the longest match is an exact match for the
// full path
func (p *PathManager) HasExactPath(path string) bool {
	p.l.RLock()
	defer p.l.RUnlock()

	if val, exceptionRaw, ok := p.paths.Root().LongestPrefix([]byte(path)); ok {
		var exception bool
		if exceptionRaw != nil {
			exception = exceptionRaw.(bool)
		}

		strVal := string(val)
		if strings.HasSuffix(strVal, "/") || strVal == path {
			return !exception
		}
	}
	return false
}
