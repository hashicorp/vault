package api

import (
	"fmt"
	"sort"
)

// Namespaces is used to query the namespace endpoints.
type Namespaces struct {
	client *Client
}

// Namespaces returns a new handle on the namespaces.
func (c *Client) Namespaces() *Namespaces {
	return &Namespaces{client: c}
}

// List is used to dump all of the namespaces.
func (n *Namespaces) List(q *QueryOptions) ([]*Namespace, *QueryMeta, error) {
	var resp []*Namespace
	qm, err := n.client.query("/v1/namespaces", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(NamespaceIndexSort(resp))
	return resp, qm, nil
}

// PrefixList is used to do a PrefixList search over namespaces
func (n *Namespaces) PrefixList(prefix string, q *QueryOptions) ([]*Namespace, *QueryMeta, error) {
	if q == nil {
		q = &QueryOptions{Prefix: prefix}
	} else {
		q.Prefix = prefix
	}

	return n.List(q)
}

// Info is used to query a single namespace by its name.
func (n *Namespaces) Info(name string, q *QueryOptions) (*Namespace, *QueryMeta, error) {
	var resp Namespace
	qm, err := n.client.query("/v1/namespace/"+name, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Register is used to register a namespace.
func (n *Namespaces) Register(namespace *Namespace, q *WriteOptions) (*WriteMeta, error) {
	wm, err := n.client.write("/v1/namespace", namespace, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Delete is used to delete a namespace
func (n *Namespaces) Delete(namespace string, q *WriteOptions) (*WriteMeta, error) {
	wm, err := n.client.delete(fmt.Sprintf("/v1/namespace/%s", namespace), nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Namespace is used to serialize a namespace.
type Namespace struct {
	Name        string
	Description string
	CreateIndex uint64
	ModifyIndex uint64
}

// NamespaceIndexSort is a wrapper to sort Namespaces by CreateIndex. We
// reverse the test so that we get the highest index first.
type NamespaceIndexSort []*Namespace

func (n NamespaceIndexSort) Len() int {
	return len(n)
}

func (n NamespaceIndexSort) Less(i, j int) bool {
	return n[i].CreateIndex > n[j].CreateIndex
}

func (n NamespaceIndexSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
