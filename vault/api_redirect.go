package vault

import (
	"context"
	"errors"
	"fmt"
	paths "path"
	"strings"
	"sync"

	"github.com/armon/go-radix"
)

type apiRedirectRegistry struct {
	lock  sync.Mutex
	paths *SpecialPathsEntry[*APIRedirect]
}

func NewAPIRedirects() *apiRedirectRegistry {
	return &apiRedirectRegistry{
		paths: &SpecialPathsEntry[*APIRedirect]{
			paths: radix.New(),
		},
	}
}

func (reg *apiRedirectRegistry) TryRegister(ctx context.Context, core *Core, mountUUID, src, dest string) error {
	if strings.HasPrefix(dest, "/") {
		return errors.New("redirect targets must be relative")
	}
	reg.lock.Lock()
	defer reg.lock.Unlock()
	found, _ := reg.paths.Match(src)
	if found {
		return fmt.Errorf("api redirect conflict for %s", src)
	}
	return reg.paths.Add(src, func(b bool) *APIRedirect {
		return &APIRedirect{
			c:             core,
			mountUUID:     mountUUID,
			prefix:        dest,
			isPrefixMatch: b,
		}
	})
}

func (reg *apiRedirectRegistry) Find(path string) *APIRedirect {
	found, e := reg.paths.Match(path)
	if found {
		return e
	}
	return nil
}

func (reg *apiRedirectRegistry) Unregister(uuid string) {
	reg.lock.Lock()
	defer reg.lock.Unlock()
	reg.paths.paths.Walk(func(k string, v interface{}) bool {
		r := v.(*APIRedirect)
		if r.mountUUID == uuid {
			reg.paths.paths.Delete(k)
			return true
		}
		return false
	})
}

type APIRedirect struct {
	c             *Core
	mountUUID     string
	prefix        string
	isPrefixMatch bool
}

func (a *APIRedirect) IsPrefixMatch() bool {
	return a.isPrefixMatch
}

func (a *APIRedirect) Destination() (string, error) {
	m := a.c.mounts.findByMountUUID(a.mountUUID)

	if m == nil {
		return "", fmt.Errorf("cannot find backend with uuid: %s", a.mountUUID)
	}
	return paths.Join(m.Namespace().Path, m.Path, a.prefix), nil
}
