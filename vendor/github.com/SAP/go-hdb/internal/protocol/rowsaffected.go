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
	rows    []int32
	_numArg int
}

func (r *rowsAffected) kind() partKind {
	return pkRowsAffected
}

func (r *rowsAffected) setNumArg(numArg int) {
	r._numArg = numArg
}

func (r *rowsAffected) read(rd *bufio.Reader) error {
	if r.rows == nil || r._numArg > cap(r.rows) {
		r.rows = make([]int32, r._numArg)
	} else {
		r.rows = r.rows[:r._numArg]
	}

	for i := 0; i < r._numArg; i++ {
		r.rows[i] = rd.ReadInt32()
		if trace {
			outLogger.Printf("rows affected %d: %d", i, r.rows[i])
		}
	}

	return rd.GetError()
}

func (r *rowsAffected) total() int64 {
	if r.rows == nil {
		return 0
	}

	total := int64(0)
	for _, rows := range r.rows {
		if rows > 0 {
			total += int64(rows)
		}
	}
	return total
}
