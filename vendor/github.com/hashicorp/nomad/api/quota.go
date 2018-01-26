package api

import (
	"fmt"
	"sort"
)

// Quotas is used to query the quotas endpoints.
type Quotas struct {
	client *Client
}

// Quotas returns a new handle on the quotas.
func (c *Client) Quotas() *Quotas {
	return &Quotas{client: c}
}

// List is used to dump all of the quota specs
func (q *Quotas) List(qo *QueryOptions) ([]*QuotaSpec, *QueryMeta, error) {
	var resp []*QuotaSpec
	qm, err := q.client.query("/v1/quotas", &resp, qo)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(QuotaSpecIndexSort(resp))
	return resp, qm, nil
}

// PrefixList is used to do a PrefixList search over quota specs
func (q *Quotas) PrefixList(prefix string, qo *QueryOptions) ([]*QuotaSpec, *QueryMeta, error) {
	if qo == nil {
		qo = &QueryOptions{Prefix: prefix}
	} else {
		qo.Prefix = prefix
	}

	return q.List(qo)
}

// ListUsage is used to dump all of the quota usages
func (q *Quotas) ListUsage(qo *QueryOptions) ([]*QuotaUsage, *QueryMeta, error) {
	var resp []*QuotaUsage
	qm, err := q.client.query("/v1/quota-usages", &resp, qo)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(QuotaUsageIndexSort(resp))
	return resp, qm, nil
}

// PrefixList is used to do a PrefixList search over quota usages
func (q *Quotas) PrefixListUsage(prefix string, qo *QueryOptions) ([]*QuotaUsage, *QueryMeta, error) {
	if qo == nil {
		qo = &QueryOptions{Prefix: prefix}
	} else {
		qo.Prefix = prefix
	}

	return q.ListUsage(qo)
}

// Info is used to query a single quota spec by its name.
func (q *Quotas) Info(name string, qo *QueryOptions) (*QuotaSpec, *QueryMeta, error) {
	var resp QuotaSpec
	qm, err := q.client.query("/v1/quota/"+name, &resp, qo)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Usage is used to query a single quota usage by its name.
func (q *Quotas) Usage(name string, qo *QueryOptions) (*QuotaUsage, *QueryMeta, error) {
	var resp QuotaUsage
	qm, err := q.client.query("/v1/quota/usage/"+name, &resp, qo)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Register is used to register a quota spec.
func (q *Quotas) Register(spec *QuotaSpec, qo *WriteOptions) (*WriteMeta, error) {
	wm, err := q.client.write("/v1/quota", spec, nil, qo)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Delete is used to delete a quota spec
func (q *Quotas) Delete(quota string, qo *WriteOptions) (*WriteMeta, error) {
	wm, err := q.client.delete(fmt.Sprintf("/v1/quota/%s", quota), nil, qo)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// QuotaSpec specifies the allowed resource usage across regions.
type QuotaSpec struct {
	// Name is the name for the quota object
	Name string

	// Description is an optional description for the quota object
	Description string

	// Limits is the set of quota limits encapsulated by this quota object. Each
	// limit applies quota in a particular region and in the future over a
	// particular priority range and datacenter set.
	Limits []*QuotaLimit

	// Raft indexes to track creation and modification
	CreateIndex uint64
	ModifyIndex uint64
}

// QuotaLimit describes the resource limit in a particular region.
type QuotaLimit struct {
	// Region is the region in which this limit has affect
	Region string

	// RegionLimit is the quota limit that applies to any allocation within a
	// referencing namespace in the region. A value of zero is treated as
	// unlimited and a negative value is treated as fully disallowed. This is
	// useful for once we support GPUs
	RegionLimit *Resources

	// Hash is the hash of the object and is used to make replication efficient.
	Hash []byte
}

// QuotaUsage is the resource usage of a Quota
type QuotaUsage struct {
	Name        string
	Used        map[string]*QuotaLimit
	CreateIndex uint64
	ModifyIndex uint64
}

// QuotaSpecIndexSort is a wrapper to sort QuotaSpecs by CreateIndex. We
// reverse the test so that we get the highest index first.
type QuotaSpecIndexSort []*QuotaSpec

func (q QuotaSpecIndexSort) Len() int {
	return len(q)
}

func (q QuotaSpecIndexSort) Less(i, j int) bool {
	return q[i].CreateIndex > q[j].CreateIndex
}

func (q QuotaSpecIndexSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// QuotaUsageIndexSort is a wrapper to sort QuotaUsages by CreateIndex. We
// reverse the test so that we get the highest index first.
type QuotaUsageIndexSort []*QuotaUsage

func (q QuotaUsageIndexSort) Len() int {
	return len(q)
}

func (q QuotaUsageIndexSort) Less(i, j int) bool {
	return q[i].CreateIndex > q[j].CreateIndex
}

func (q QuotaUsageIndexSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// QuotaLimitSort is a wrapper to sort QuotaLimits
type QuotaLimitSort []*QuotaLimit

func (q QuotaLimitSort) Len() int {
	return len(q)
}

func (q QuotaLimitSort) Less(i, j int) bool {
	return q[i].Region < q[j].Region
}

func (q QuotaLimitSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}
