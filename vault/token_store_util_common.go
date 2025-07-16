// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

const sscGenCounterPath string = "core/sscGenCounter/"

type SSCTokenGenerationCounter struct {
	Counter int
}

func (ts *TokenStore) GetSSCTokensGenerationCounter() int {
	return ts.sscTokensGenerationCounter.Counter
}

func (ts *TokenStore) loadSSCTokensGenerationCounter(ctx context.Context) error {
	sscTokensGenerationCounterStorageVal, err := ts.core.barrier.Get(ctx, sscGenCounterPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve SSCTokenGenerationCounter from storage: err %w", err)
	}
	if sscTokensGenerationCounterStorageVal == nil {
		ts.logger.Trace("no token generation counter found in storage")
		ts.sscTokensGenerationCounter = SSCTokenGenerationCounter{Counter: 0}
		return nil
	}
	var sscTokensGenerationCounter SSCTokenGenerationCounter
	err = json.Unmarshal(sscTokensGenerationCounterStorageVal.Value, &sscTokensGenerationCounter)
	if err != nil {
		return fmt.Errorf("malformed token generation counter found in storage: err %w", err)
	}

	ts.logger.Debug("loaded ssct generation counter", "generation", sscTokensGenerationCounter.Counter)
	ts.sscTokensGenerationCounter = sscTokensGenerationCounter
	return nil
}

func (ts *TokenStore) UpdateSSCTokensGenerationCounter(ctx context.Context) error {
	if err := ts.loadSSCTokensGenerationCounter(ctx); err != nil {
		return err
	}
	ts.sscTokensGenerationCounter.Counter += 1
	if ts.sscTokensGenerationCounter.Counter <= 0 {
		// Don't store the 0 value
		ts.logger.Warn("attempt to store non-positive token generation counter was ignored",
			"sscTokensGenerationCounter", ts.sscTokensGenerationCounter.Counter)
	}
	marshalledCtr, err := json.Marshal(ts.sscTokensGenerationCounter)
	if err != nil {
		return err
	}
	err = ts.core.barrier.Put(ctx, &logical.StorageEntry{
		Key:   sscGenCounterPath,
		Value: marshalledCtr,
	})
	if err != nil {
		return err
	}

	ts.logger.Debug("updated ssct generation counter", "generation", ts.sscTokensGenerationCounter.Counter)
	return nil
}
