// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parsing

import "encoding/asn1"

// Asn1UnmarshallNoTrailing is a wrapper around asn1.Unmarshal that ensures there
// is no trailing data in the input returning an error if there is.
func Asn1UnmarshallNoTrailing(b []byte, val any) error {
	rest, err := asn1.Unmarshal(b, val)
	if err != nil {
		return err
	} else if len(rest) != 0 {
		return asn1.SyntaxError{Msg: "trailing data"}
	}
	return nil
}
