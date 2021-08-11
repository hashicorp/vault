// +build !enterprise

package vault

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	passwordPolicySubPath = "password_policy/"
)

// retrievePasswordPolicy retrieves a password policy from the logical storage
func (d dynamicSystemView) retrievePasswordPolicy(ctx context.Context, policyName string) (*passwordPolicyConfig, error) {
	storage := d.core.systemBarrierView.SubView(passwordPolicySubPath)
	entry, err := storage.Get(ctx, policyName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	policyCfg := &passwordPolicyConfig{}
	err = json.Unmarshal(entry.Value, &policyCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored data: %w", err)
	}

	return policyCfg, nil
}
