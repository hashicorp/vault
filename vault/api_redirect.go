// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	url2 "net/url"
	"strings"
	"sync"

	"github.com/armon/go-radix"
)

type wellKnownRedirectRegistry struct {
	lock  sync.Mutex
	paths *radix.Tree
}

func NewWellKnownRedirects() *wellKnownRedirectRegistry {
	return &wellKnownRedirectRegistry{
		paths: radix.New(),
	}
}

func (reg *wellKnownRedirectRegistry) TryRegister(ctx context.Context, core *Core, mountUUID, src, dest string) error {
	if strings.HasPrefix(dest, "/") {
		return errors.New("redirect targets must be relative")
	}
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

func (reg *wellKnownRedirectRegistry) Find(path string) (*wellKnownRedirect, string) {
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

func (reg *wellKnownRedirectRegistry) DeregisterMount(uuid string) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	reg.paths.Walk(func(k string, v interface{}) bool {
		r := v.(*wellKnownRedirect)
		if r.mountUUID == uuid {
			reg.paths.Delete(k)
			return true
		}
		return false
	})
}

func (reg *wellKnownRedirectRegistry) DeregisterSource(uuid, src string) bool {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	var found bool
	reg.paths.Walk(func(k string, v interface{}) bool {
		r := v.(*wellKnownRedirect)
		if r.mountUUID == uuid && k == src {
			found = true
			reg.paths.Delete(k)
			return true
		}
		return false
	})
	return found
}

type wellKnownRedirect struct {
	c             *Core
	mountUUID     string
	prefix        string
	isPrefixMatch bool
}

func (a *wellKnownRedirect) Destination(remaining string) (string, error) {
	var destPath string
	if a.c == nil {
		// Just for testing
		destPath = a.prefix
	} else {
		m := a.c.mounts.findByMountUUID(a.mountUUID)

		if m == nil {
			return "", fmt.Errorf("cannot find backend with uuid: %s", a.mountUUID)
		}
		var err error
		destPath, err = url2.JoinPath(m.Namespace().Path, m.Path, a.prefix)
		if err != nil {
			return "", err
		}
	}

	u := url2.URL{
		Path: destPath + "/",
	}
	r, err := url2.Parse(remaining)
	if err != nil {
		return "", err
	}
	dest := u.ResolveReference(r)
	dest.Path = strings.TrimSuffix(dest.Path, "/")
	return dest.String(), nil
}
