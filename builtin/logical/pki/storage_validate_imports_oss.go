// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

// Common Criteria is an Ent-Only Feature
func enforceCommonCriteriaOnUploadedCAs(sc *storageContext, issuerKeyMap map[string]string, createdIssuers []string) error {
	return nil
}
