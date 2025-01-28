// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/vault/helper/activationflags"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	paramFeatureName = "feature_name"
	descFeatureName  = "The name of the feature to be activated."
	summaryList      = "Returns the available and activated activation-flagged features."
	summaryUpdate    = "Activate a flagged feature."

	prefixActivationFlags         = "activation-flags"
	verbActivationFlagsActivate   = "activate"
	verbActivationFlagsDeactivate = "deactivate"

	fieldActivated   = "activated"
	fieldUnactivated = "unactivated"

	helpSynopsis    = "Returns information about Vault's features that require a one-time activation step."
	helpDescription = `
This path responds to the following HTTP methods.
	GET /
		Returns the available and activated activation-flags.

	PUT|POST /<feature-name>/activate
		Activates the specified feature. Cannot be undone.`

	activationFlagTest = "activation-test"
)

// These variables should only be mutated during initialization or server construction.
// It is unsafe to modify them once the Vault core is running.
var (
	readActivationFlag = func(ctx context.Context, b *SystemBackend, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
		return b.readActivationFlag(ctx, req, fd)
	}

	writeActivationFlag = func(ctx context.Context, b *SystemBackend, req *logical.Request, fd *framework.FieldData, isActivate bool) (*logical.Response, error) {
		return b.writeActivationFlagWrite(ctx, req, fd, isActivate)
	}
)

func (b *SystemBackend) activationFlagsPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: fmt.Sprintf("%s$", prefixActivationFlags),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationVerb:   "read",
				OperationSuffix: prefixActivationFlags,
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleActivationFlagRead,
					Summary:  summaryList,
				},
			},
			HelpSynopsis:    helpSynopsis,
			HelpDescription: helpDescription,
		},
		{
			Pattern: fmt.Sprintf("%s/%s/%s", prefixActivationFlags, activationFlagTest, verbActivationFlagsActivate),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: prefixActivationFlags,
				OperationVerb:   verbActivationFlagsActivate,
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.handleActivationFlagsActivate,
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
					Summary:                     summaryUpdate,
				},
			},
			HelpSynopsis:    helpSynopsis,
			HelpDescription: helpDescription,
		},
		{
			Pattern: fmt.Sprintf("%s/%s/%s", prefixActivationFlags, activationflags.IdentityDeduplication, verbActivationFlagsActivate),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: prefixActivationFlags,
				OperationVerb:   verbActivationFlagsActivate,
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.handleActivationFlagsActivate,
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
					Summary:                     summaryUpdate,
				},
			},
			HelpSynopsis:    helpSynopsis,
			HelpDescription: helpDescription,
		},
	}
}

func (b *SystemBackend) handleActivationFlagRead(ctx context.Context, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
	return readActivationFlag(ctx, b, req, fd)
}

func (b *SystemBackend) handleActivationFlagsActivate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return writeActivationFlag(ctx, b, req, data, true)
}

func (b *SystemBackend) readActivationFlag(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	activationFlags, err := b.Core.FeatureActivationFlags.Get(ctx)
	if err != nil {
		return nil, err
	}

	return b.activationFlagsToResponse(activationFlags), nil
}

func (b *SystemBackend) writeActivationFlagWrite(ctx context.Context, req *logical.Request, _ *framework.FieldData, isActivate bool) (*logical.Response, error) {
	// We need to manually parse out the feature_name from the path because we can't use FieldSchema parameters
	// in the path to make generic endpoints. We need each activation-flag path to be a separate endpoint.
	// Path starts out as activation-flags/<feature_name>/verb
	// Removes activation-flags/ from the path
	trimPrefix := strings.TrimPrefix(req.Path, prefixActivationFlags+"/")
	// Removes /verb from the path
	featureName := trimPrefix[:strings.LastIndex(trimPrefix, "/")]

	err := b.Core.FeatureActivationFlags.Write(ctx, featureName, isActivate)
	if err != nil {
		return nil, fmt.Errorf("failed to write new activation flags: %w", err)
	}

	// We read back the value after writing it to storage so that we can try forcing a cache update right away.
	// If this fails, it's still okay to proceed as the write has been successful and the cache will get updated
	// at the time of an endpoint getting called. However, we can only return the one feature name we just activated
	// in the response since the read to retrieve any others did not succeed.
	activationFlags, err := b.Core.FeatureActivationFlags.Get(ctx)
	if err != nil {
		resp := b.activationFlagsToResponse([]string{featureName})
		return resp, fmt.Errorf("failed to read activation-flags back after write: %w", err)
	}

	return b.activationFlagsToResponse(activationFlags), nil
}

func (b *SystemBackend) activationFlagsToResponse(activationFlags []string) *logical.Response {
	slices.Sort(activationFlags)
	return &logical.Response{
		Data: map[string]interface{}{
			fieldActivated: activationFlags,
		},
	}
}
