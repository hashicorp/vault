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
	paths *radix.Tree
}

func NewAPIRedirects() *apiRedirectRegistry {
	return &apiRedirectRegistry{
		paths: radix.New(),
	}
}

func (reg *apiRedirectRegistry) TryRegister(ctx context.Context, core *Core, mountUUID, src, dest string) error {
	if strings.HasPrefix(dest, "/") {
		return errors.New("redirect targets must be relative")
	}
	reg.lock.Lock()
	defer reg.lock.Unlock()
	_, _, found := reg.paths.LongestPrefix(src)
	if found {
		return fmt.Errorf("api redirect conflict for %s", src)
	}
	_, found = reg.paths.Insert(src, &APIRedirect{
		c:         core,
		mountUUID: mountUUID,
		prefix:    dest,
	})
	if found {
		panic("somehow had a duplicate even though lock was held")
	}
	return nil
}

func (reg *apiRedirectRegistry) Find(path string) *APIRedirect {
	_, r, found := reg.paths.LongestPrefix(path)
	if found {
		return r.(*APIRedirect)
	}
	return nil
}

type APIRedirect struct {
	c         *Core
	mountUUID string
	prefix    string
}

func (a *APIRedirect) Destination() (string, error) {
	m := a.c.mounts.findByMountUUID(a.mountUUID)

	if m == nil {
		return "", fmt.Errorf("cannot find backend with uuid: %s", a.mountUUID)
	}
	return paths.Join(m.Namespace().Path, m.Path, a.prefix), nil
}
