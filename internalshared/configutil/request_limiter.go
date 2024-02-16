// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

type RequestLimiter struct {
	UnusedKeys UnusedKeyMap `hcl:",unusedKeyPositions"`

	Disable    bool        `hcl:"-"`
	DisableRaw interface{} `hcl:"disable"`
}

func (r *RequestLimiter) Validate(source string) []ConfigError {
	return ValidateUnusedFields(r.UnusedKeys, source)
}

func (r *RequestLimiter) GoString() string {
	return fmt.Sprintf("*%#v", *r)
}

var DefaultRequestLimiter = &RequestLimiter{
	Disable: true,
}

func parseRequestLimiter(result *SharedConfig, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'request_limiter' block is permitted")
	}

	result.RequestLimiter = DefaultRequestLimiter

	// Get our one item
	item := list.Items[0]

	if err := hcl.DecodeObject(&result.RequestLimiter, item.Val); err != nil {
		return multierror.Prefix(err, "request_limiter:")
	}

	result.RequestLimiter.Disable = true
	if result.RequestLimiter.DisableRaw != nil {
		var err error
		if result.RequestLimiter.Disable, err = parseutil.ParseBool(result.RequestLimiter.DisableRaw); err != nil {
			return err
		}
		result.RequestLimiter.DisableRaw = nil
	}

	return nil
}
