// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package crdb

import (
	"context"
	"fmt"
	"sync"
)

// WriteSkewTest abstracts the operations that needs to be performed by a
// particular framework for the purposes of TestExecuteTx. This allows the test
// to be written once and run for any framework supported by this library.
type WriteSkewTest interface {
	Init(context.Context) error
	ExecuteTx(ctx context.Context, fn func(tx interface{}) error) error
	GetBalances(ctx context.Context, tx interface{}) (bal1, bal2 int, err error)
	UpdateBalance(ctx context.Context, tx interface{}, acct, delta int) error
}

// ExecuteTxGenericTest represents the structure of a test for the ExecuteTx
// function. The actual database operations are abstracted by framework; the
// idea is that tests for different frameworks implement that interface and then
// invoke this test.
//
// The test interleaves two transactions such that one of them will require a
// restart because of write skew.
func ExecuteTxGenericTest(ctx context.Context, framework WriteSkewTest) error {
	framework.Init(ctx)
	// wg is used as a barrier, blocking each transaction after it performs the
	// initial read until they both read.
	var wg sync.WaitGroup
	wg.Add(2)
	runTxn := func(iter *int) <-chan error {
		errCh := make(chan error, 1)
		go func() {
			*iter = 0
			errCh <- framework.ExecuteTx(ctx, func(tx interface{}) (retErr error) {
				defer func() {
					if retErr == nil {
						return
					}
					// Wrap the error so that we test the library's unwrapping.
					retErr = testError{cause: retErr}
				}()

				*iter++
				bal1, bal2, err := framework.GetBalances(ctx, tx)
				if err != nil {
					return err
				}
				// If this is the first iteration, wait for the other tx to also read.
				if *iter == 1 {
					wg.Done()
					wg.Wait()
				}
				// Now, subtract from one account and give to the other.
				if bal1 > bal2 {
					if err := framework.UpdateBalance(ctx, tx, 1, -100); err != nil {
						return err
					}
					if err := framework.UpdateBalance(ctx, tx, 2, +100); err != nil {
						return err
					}
				} else {
					if err := framework.UpdateBalance(ctx, tx, 1, +100); err != nil {
						return err
					}
					if err := framework.UpdateBalance(ctx, tx, 2, -100); err != nil {
						return err
					}
				}
				return nil
			})
		}()
		return errCh
	}

	var iters1, iters2 int
	txn1Err := runTxn(&iters1)
	txn2Err := runTxn(&iters2)
	if err := <-txn1Err; err != nil {
		return fmt.Errorf("expected success in txn1; got %s", err)
	}
	if err := <-txn2Err; err != nil {
		return fmt.Errorf("expected success in txn2; got %s", err)
	}
	if iters1+iters2 <= 2 {
		return fmt.Errorf("expected at least one retry between the competing transactions; "+
			"got txn1=%d, txn2=%d", iters1, iters2)
	}

	var bal1, bal2 int
	err := framework.ExecuteTx(ctx, func(txi interface{}) error {
		var err error
		bal1, bal2, err = framework.GetBalances(ctx, txi)
		return err
	})
	if err != nil {
		return err
	}
	if bal1 != 100 || bal2 != 100 {
		return fmt.Errorf("expected balances to be restored without error; "+
			"got acct1=%d, acct2=%d: %s", bal1, bal2, err)
	}
	return nil
}

type testError struct {
	cause error
}

func (t testError) Error() string {
	return "test error"
}

func (t testError) Unwrap() error {
	return t.cause
}
