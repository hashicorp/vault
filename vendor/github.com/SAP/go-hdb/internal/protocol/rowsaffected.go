/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import (
	"github.com/SAP/go-hdb/internal/bufio"
)

//rows affected
const (
	raSuccessNoInfo   = -2
	raExecutionFailed = -3
)

//rows affected
type rowsAffected struct {
	sums    []int32
	_numArg int
}

func (r *rowsAffected) kind() partKind {
	return pkRowsAffected
}

func (r *rowsAffected) setNumArg(numArg int) {
	r._numArg = numArg
}

func (r *rowsAffected) read(rd *bufio.Reader) error {
	if r.sums == nil || r._numArg > cap(r.sums) {
		r.sums = make([]int32, r._numArg)
	} else {
		r.sums = r.sums[:r._numArg]
	}

	var err error

	for i := 0; i < r._numArg; i++ {
		r.sums[i], err = rd.ReadInt32()
		if err != nil {
			return err
		}
	}

	if trace {
		outLogger.Printf("rows affected %v", r.sums)
	}

	return nil
}

func (r *rowsAffected) total() int64 {
	if r.sums == nil {
		return 0
	}

	total := int64(0)
	for _, sum := range r.sums {
		total += int64(sum)
	}
	return total
}
