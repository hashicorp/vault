// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"

	"github.com/aerospike/aerospike-client-go/v5/logger"
)

type werrGroup struct {
	sem  *semaphore.Weighted
	ctx  context.Context
	wg   sync.WaitGroup
	el   sync.Mutex
	errs Error

	// function to defer; used for recordset signals
	f func()
}

func newWeightedErrGroup(maxConcurrency int) *werrGroup {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	return &werrGroup{
		sem: semaphore.NewWeighted(int64(maxConcurrency)),
		ctx: context.Background(),
	}
}

func (weg *werrGroup) execute(cmd command) {
	weg.wg.Add(1)

	if err := weg.sem.Acquire(weg.ctx, 1); err != nil {
		logger.Logger.Error("Constraint Semaphore failed: %s", err.Error())
	}

	go func() {
		defer weg.sem.Release(1)
		defer weg.wg.Done()
		if weg.f != nil {
			defer weg.f()
		}

		if err := cmd.Execute(); err != nil {
			weg.el.Lock()
			weg.errs = chainErrors(err, weg.errs)
			weg.el.Unlock()
		}
	}()
}

func (weg *werrGroup) wait() Error {
	weg.wg.Wait()
	return weg.errs
}
