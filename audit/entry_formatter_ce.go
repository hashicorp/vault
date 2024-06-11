// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"errors"
)

func (f *entryFormatter) shouldExclude() bool {
	return false
}

func (f *entryFormatter) excludeFields(entry any) (map[string]any, error) {
	return nil, errors.New("enterprise-only feature: audit exclusion")
}
