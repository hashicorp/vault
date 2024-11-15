// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"fmt"
	"io"

	// Ignore gosec warning "Blocklisted import crypto/md5: weak cryptographic primitive". We need
	// to use MD5 here to implement the SCRAM specification.
	/* #nosec G501 */
	"crypto/md5"
)

const defaultAuthDB = "admin"

func mongoPasswordDigest(username, password string) string {
	// Ignore gosec warning "Use of weak cryptographic primitive". We need to use MD5 here to
	// implement the SCRAM specification.
	/* #nosec G401 */
	h := md5.New()
	_, _ = io.WriteString(h, username)
	_, _ = io.WriteString(h, ":mongo:")
	_, _ = io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
