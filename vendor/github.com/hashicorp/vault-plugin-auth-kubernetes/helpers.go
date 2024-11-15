// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"net/http"
	"strings"
)

func setRequestHeader(req *http.Request, bearer string) {
	bearer = strings.TrimSpace(bearer)

	// Set the JWT as the Bearer token
	req.Header.Set("Authorization", bearer)

	// Set the MIME type headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}
