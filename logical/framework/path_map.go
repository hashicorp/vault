package framework

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
)

// PathMap can be used to generate a path that stores mappings in the
// storage. It is a structure that also exports functions for querying the
// mappings.
//
// The primary use case for this is for credential providers to do their
// mapping to policies.
type PathMap struct {
	Name string
}

// Get reads a value out of the mapping
func (p *PathMap) Get(s logical.Storage, k string) (string, error) {
	entry, err := s.Get(fmt.Sprintf("map/%s/%s", p.Name, k))
	if err != nil {
		return "", err
	}
	if entry == nil {
		return "", nil
	}

	return string(entry.Value), nil
}

// Put writes a value into the mapping
func (p *PathMap) Put(s logical.Storage, k string, v string) error {
	return s.Put(&logical.StorageEntry{
		Key:   fmt.Sprintf("map/%s/%s", p.Name, k),
		Value: []byte(v),
	})
}

// Paths are the paths to append to the Backend paths.
func (p *PathMap) Paths() []*Path {
	return []*Path{
		&Path{
			Pattern: fmt.Sprintf("map/%s$", p.Name),

			Callbacks: map[logical.Operation]OperationFunc{
				logical.ListOperation: p.pathList,
				logical.ReadOperation: p.pathList,
			},

			HelpSynopsis: fmt.Sprintf("Read mappings for %s", p.Name),
		},

		&Path{
			Pattern: fmt.Sprintf("map/%s/(?P<key>\\w+)", p.Name),

			Fields: map[string]*FieldSchema{
				"key": &FieldSchema{
					Type:        TypeString,
					Description: "Key for the mapping",
				},

				"value": &FieldSchema{
					Type:        TypeString,
					Description: "Value for the mapping",
				},
			},

			Callbacks: map[logical.Operation]OperationFunc{
				logical.WriteOperation: p.pathSingleWrite,
				logical.ReadOperation:  p.pathSingleRead,
			},

			HelpSynopsis: fmt.Sprintf("Read/write a single %s mapping", p.Name),
		},
	}
}

func (p *PathMap) pathList(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	keys, err := req.Storage.List(req.Path)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(keys), nil
}

func (p *PathMap) pathSingleRead(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	v, err := p.Get(req.Storage, d.Get("key").(string))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": v,
		},
	}, nil
}

func (p *PathMap) pathSingleWrite(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	err := p.Put(
		req.Storage,
		d.Get("key").(string), d.Get("value").(string))
	return nil, err
}
