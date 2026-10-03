// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import "context"

type ACMEBillingSystemView interface {
	CreateActivityCountEventForIdentifiers(ctx context.Context, identifiers []string) error
}
