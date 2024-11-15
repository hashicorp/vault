/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sl

import (
	"math"
)

var DefaultLimit = 50

// Options contains the individual query parameters that can be applied to a request.
type Options struct {
	Id         *int
	Mask       string
	Filter     string
	Limit      *int
	Offset     *int
	TotalItems int
}

// returns Math.Ciel((TotalItems - Limit) / Limit)
func (opt *Options) GetRemainingAPICalls() int {
	Total := float64(opt.TotalItems)
	Limit := float64(*opt.Limit)
	return int(math.Ceil((Total - Limit) / Limit))
}

// Makes sure the limit is set to something, not 0 or 1. Will set to default if no other limit is set.
func (opt *Options) ValidateLimit() int {
	if opt.Limit == nil || *opt.Limit < 2 {
		opt.Limit = &DefaultLimit
	}
	return *opt.Limit
}

func (opt *Options) SetTotalItems(total int) {
	opt.TotalItems = total
}

func (opt *Options) SetOffset(offset int) {
	opt.Offset = &offset
}

func (opt *Options) SetLimit(limit int) {
	opt.Limit = &limit
}
