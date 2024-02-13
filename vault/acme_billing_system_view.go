// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-kms-wrapping/entropy/v2"
	"github.com/hashicorp/vault/sdk/logical"
)

// See comment in command/format.go
const hopeDelim = "♨"

// Modifying this constant will result in previous activities not being
// associated with future activities of this same type. A storage
// migration would need to occur to fix this.
//
// Suggested naming precedence: <plugin name>-<client type>. Since this
// comes from the builtin PKI plugin, and we're counting ACME certificates,
// this becomes pki-acme.
const ACMEActivityType = "pki-acme"

// acmeBillingImpl is the (single) implementation of the actual client
// counting interface. It needs a reference to core (as per discussions
// with Mike, in the future the activityLog reference will no longer be
// static but may be replaced throughout the lifecycle of a core instance)
// and a reference to the mount that is being counted.
type acmeBillingImpl struct {
	core  *Core
	entry *MountEntry
}

var _ logical.ACMEBillingSystemView = (*acmeBillingImpl)(nil)

// Due to unfortunate layering of system view interfaces, there are three
// possible sets of interfaces we need to layer with this ACME interface:
//
// 1. Everything: a managed key system view, an entropy sourcer, and an
//    extended system view.
// 2. Managed keys without an entropy sourcer.
// 3. Just extended system view.
//
// Unfortunately, just using acmeBillingSystemViewImpl is not sufficient:
// because of the embedded interfaces, even when these are nil, the
// implementation will claim to support these interfaces, (thus, their cast
// being accepted by the toolchain), but when the caller goes to use the new
// value, they are hit with a nil panic.
//
// This is unfortunate.
//
// To avoid this, we create three possible implementations and use the lowest
// common denominator of the caller of NewAcmeBillingSystemView(...) to assume
// it is _just_ an extendedSystemView implementation.

// Scenario 1 above.
type acmeBillingSystemViewImpl struct {
	extendedSystemView
	logical.ManagedKeySystemView
	entropy.Sourcer
	acmeBillingImpl
}

var (
	_ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImpl)(nil)
	_ extendedSystemView            = (*acmeBillingSystemViewImpl)(nil)
	_ logical.ManagedKeySystemView  = (*acmeBillingSystemViewImpl)(nil)
	_ entropy.Sourcer               = (*acmeBillingSystemViewImpl)(nil)
)

// Scenario 2 above.
type acmeBillingSystemViewImplNoSourcer struct {
	extendedSystemView
	logical.ManagedKeySystemView
	acmeBillingImpl
}

var (
	_ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImplNoSourcer)(nil)
	_ extendedSystemView            = (*acmeBillingSystemViewImplNoSourcer)(nil)
	_ logical.ManagedKeySystemView  = (*acmeBillingSystemViewImplNoSourcer)(nil)
)

// Scenario 3 above.
type acmeBillingSystemViewImplNoManagedKeys struct {
	extendedSystemView
	acmeBillingImpl
}

var (
	_ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImplNoManagedKeys)(nil)
	_ extendedSystemView            = (*acmeBillingSystemViewImplNoManagedKeys)(nil)
)

// NewAcmeBillingSystemView creates the appropriate implementation based on
// the passed arguments, mapping them to the above scenarios. We further
// restrict sysView to have a dynamicSystemView implementation, to get the
// mount entry out of.
func (c *Core) NewAcmeBillingSystemView(sysView interface{}) extendedSystemView {
	es := sysView.(extendedSystemViewImpl)
	des := es.dynamicSystemView

	// Scenario 3.
	return &acmeBillingSystemViewImplNoManagedKeys{
		extendedSystemView: es,
		acmeBillingImpl: acmeBillingImpl{
			core:  c,
			entry: des.mountEntry,
		},
	}
}

func (a *acmeBillingImpl) CreateActivityCountEventForIdentifiers(ctx context.Context, identifiers []string) error {
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
	a.core.activityLogLock.RLock()
	activityLog := a.core.activityLog
	a.core.activityLogLock.RUnlock()
	if activityLog == nil {
		return nil
	}
	activityLog.logger.Debug(fmt.Sprintf("Handling ACME client count event for [%v] -> %v", identifiers, clientID))
	activityLog.AddActivityToFragment(clientID, a.entry.NamespaceID, time.Now().Unix(), ACMEActivityType, a.entry.Accessor)

	return nil
}
