// Copyright 2013-2020 Aerospike, Inc.
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

// ScanPolicy encapsulates parameters used in scan operations.
type ScanPolicy struct {
	MultiPolicy

	// ScanPercent determines percent of data to scan.
	// Valid integer range is 1 to 100.
	// Default is 100.
	ScanPercent int //= 100;

	// ConcurrentNodes determines how to issue scan requests (in parallel or sequentially).
	ConcurrentNodes bool //= true;
}

// NewScanPolicy creates a new ScanPolicy instance with default values.
func NewScanPolicy() *ScanPolicy {
	mp := *NewMultiPolicy()
	mp.TotalTimeout = 0

	return &ScanPolicy{
		MultiPolicy:     mp,
		ScanPercent:     100,
		ConcurrentNodes: true,
	}
}
