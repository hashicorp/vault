// Copyright 2014-2021 Aerospike, Inc.
//
// Portions may be licensed to Aerospike, Inc. under one or more contributor
// license agreements WHICH ARE COMPATIBLE WITH THE APACHE LICENSE, VERSION 2.0.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package aerospike

// PartitionFilter is used in scan/queries.
type PartitionFilter struct {
	begin      int
	count      int
	digest     []byte
	partitions []*partitionStatus
	done       bool
}

// NewPartitionFilterAll creates a partition filter that
// reads all the partitions.
func NewPartitionFilterAll() *PartitionFilter {
	return newPartitionFilter(0, _PARTITIONS)
}

// NewPartitionFilterById creates a partition filter by partition id.
// Partition id is between 0 - 4095
func NewPartitionFilterById(partitionId int) *PartitionFilter {
	return newPartitionFilter(partitionId, 1)
}

// NewPartitionFilterByRange creates a partition filter by partition range.
// begin partition id is between 0 - 4095
// count is the number of partitions, in the range of 1 - 4096 inclusive.
func NewPartitionFilterByRange(begin, count int) *PartitionFilter {
	return newPartitionFilter(begin, count)
}

// NewPartitionFilterByKey creates a partition filter that will return
// records after key's digest in the partition containing the digest.
func NewPartitionFilterByKey(key *Key) *PartitionFilter {
	return &PartitionFilter{begin: key.PartitionId(), count: 1, digest: key.Digest()}
}

func newPartitionFilter(begin, count int) *PartitionFilter {
	return &PartitionFilter{begin: begin, count: count}
}

// IsDone returns - if using ScanPolicy.MaxRecords or QueryPolicy,MaxRecords -
// if the previous paginated scans with this partition filter instance return all records?
func (pf *PartitionFilter) IsDone() bool {
	return pf.done
}
