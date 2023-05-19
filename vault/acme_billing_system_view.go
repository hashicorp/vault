// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// See comment in command/format.go
const hopeDelim = "♨"

type acmeBillingSystemViewImpl struct {
	extendedSystemView
	logical.ManagedKeySystemView

	core  *Core
	entry *MountEntry
}

var _ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImpl)(nil)

func (c *Core) NewAcmeBillingSystemView(sysView interface{}, managed logical.ManagedKeySystemView) *acmeBillingSystemViewImpl {
	es := sysView.(extendedSystemViewImpl)
	des := es.dynamicSystemView

	return &acmeBillingSystemViewImpl{
		extendedSystemView:   es,
		ManagedKeySystemView: managed,
		core:                 c,
		entry:                des.mountEntry,
	}
}

func (a *acmeBillingSystemViewImpl) CreateActivityCountEventForIdentifiers(ctx context.Context, identifiers []string) error {
	// Fake our clientID from the identifiers, but ensure it is
	// independent of ordering.
	//
	// TODO: Because of prefixing currently handled by AddActivityToFragment,
	// we do not need to ensure it is globally unique.
	sort.Strings(identifiers)
	joinedIdentifiers := "[" + strings.Join(identifiers, "]"+hopeDelim+"[") + "]"
	identifiersHash := sha256.Sum256([]byte(joinedIdentifiers))
	clientID := base64.RawURLEncoding.EncodeToString(identifiersHash[:])

	// Log so users can correlate ACME requests to client count tokens.
	activityType := "acme"
	a.core.activityLogLock.RLock()
	activityLog := a.core.activityLog
	a.core.activityLogLock.RUnlock()
	if activityLog == nil {
		return nil
	}
	activityLog.logger.Debug(fmt.Sprintf("Handling ACME client count event for [%v] -> %v", identifiers, clientID))
	activityLog.AddActivityToFragment(clientID, a.entry.NamespaceID, time.Now().Unix(), activityType, a.entry.Accessor)

	return nil
}
