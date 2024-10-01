// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/armon/go-radix"
)

type wellKnownRedirect struct {
	c         *Core
	mountUUID string
	prefix    string
}

type wellKnownRedirectRegistry struct {
	lock  sync.RWMutex
	paths *radix.Tree
}

func NewWellKnownRedirects() *wellKnownRedirectRegistry {
	return &wellKnownRedirectRegistry{
		paths: radix.New(),
	}
}

// Attempt to register a mapping from /.well-known/_src_ to /v1/_mount-path_/_dest_
func (reg *wellKnownRedirectRegistry) TryRegister(ctx context.Context, core *Core, mountUUID, src, dest string) error {
	if strings.HasPrefix(dest, "/") {
		return errors.New("redirect targets must be relative")
	}
	src = strings.TrimSuffix(src, "/")
	reg.lock.Lock()
	defer reg.lock.Unlock()
	_, _, found := reg.paths.LongestPrefix(src)
	if found {
		return fmt.Errorf("api redirect conflict for %s", src)
	}
	reg.paths.Insert(src, &wellKnownRedirect{
		c:         core,
		mountUUID: mountUUID,
		prefix:    dest,
	})
	return nil
}

// Find any relevant redirects for a given source path
func (reg *wellKnownRedirectRegistry) Find(path string) (*wellKnownRedirect, string) {
	reg.lock.RLock()
	defer reg.lock.RUnlock()

	s, a, found := reg.paths.LongestPrefix(path)
	if found {
		remaining := strings.TrimPrefix(path, s)
		if len(remaining) > 0 {
			switch remaining[0] {
			case '/':
				remaining = remaining[1:]
			case '?':
			default:
				// This isn't an exact path match
				return nil, ""
			}
		}
		return a.(*wellKnownRedirect), remaining
	}
	return nil, ""
}

// Remove all redirects for a given mount
func (reg *wellKnownRedirectRegistry) DeregisterMount(mountUuid string) {
	reg.lock.Lock()
	defer reg.lock.Unlock()

	var toDelete []string
	reg.paths.Walk(func(k string, v interface{}) bool {
		r := v.(*wellKnownRedirect)
		if r.mountUUID == mountUuid {
			toDelete = append(toDelete, k)
		}
		return false
	})
	for _, d := range toDelete {
		reg.paths.Delete(d)
	}
}

// Remove a specific redirect for a mount
func (reg *wellKnownRedirectRegistry) DeregisterSource(mountUuid, src string) bool {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	var found bool
	reg.paths.Walk(func(k string, v interface{}) bool {
		r := v.(*wellKnownRedirect)
		if r.mountUUID == mountUuid && k == src {
			found = true
			reg.paths.Delete(k)
			return true
		}
		return false
	})
	return found
}

// List returns a map keyed by the registered label and associated well known redirect
func (reg *wellKnownRedirectRegistry) List() map[string]*wellKnownRedirect {
	reg.lock.RLock()
	defer reg.lock.RUnlock()

	labels := map[string]*wellKnownRedirect{}
	reg.paths.Walk(func(s string, v interface{}) bool {
		labels[s] = v.(*wellKnownRedirect)
		return false
	})

	return labels
}

// Get returns a well known redirect for a specific registered label
func (reg *wellKnownRedirectRegistry) Get(label string) (*wellKnownRedirect, bool) {
	reg.lock.RLock()
	defer reg.lock.RUnlock()

	if v, ok := reg.paths.Get(label); ok {
		return v.(*wellKnownRedirect), true
	}

	return nil, false
}

// Construct the full destination of the redirect, including any remaining path past the src
func (a *wellKnownRedirect) Destination(remaining string) (string, error) {
	var destPath string
	if a.c == nil {
		// Just for testing
		destPath = a.prefix
	} else {
		m := a.c.router.MatchingMountByUUID(a.mountUUID)

		if m == nil {
			return "", fmt.Errorf("cannot find backend with uuid: %s", a.mountUUID)
		}
		var err error
		destPath, err = url.JoinPath(m.Namespace().Path, m.Path, a.prefix)
		if err != nil {
			return "", err
		}
	}

	u := url.URL{
		Path: destPath + "/",
	}
	r, err := url.Parse(remaining)
	if err != nil {
		return "", err
	}
	dest := u.ResolveReference(r)
	dest.Path = strings.TrimSuffix(dest.Path, "/")
	return dest.String(), nil
}
