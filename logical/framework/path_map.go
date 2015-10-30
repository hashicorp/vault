package framework

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

// PathMap can be used to generate a path that stores mappings in the
// storage. It is a structure that also exports functions for querying the
// mappings.
//
// The primary use case for this is for credential providers to do their
// mapping to policies.
type PathMap struct {
	Prefix        string
	Name          string
	Schema        map[string]*FieldSchema
	CaseSensitive bool
	Salt          *salt.Salt

	// Allows the ability to intercept the request and modify it or cancel before
	// the default task is performed here
	Callbacks     map[logical.Operation]OperationFunc

	once sync.Once
}

func (p *PathMap) init() {
	if p.Prefix == "" {
		p.Prefix = "map"
	}

	if p.Schema == nil {
		p.Schema = map[string]*FieldSchema{
			"value": &FieldSchema{
				Type:        TypeString,
				Description: fmt.Sprintf("Value for %s mapping", p.Name),
			},
		}
	}
}

// pathStruct returns the pathStruct for this mapping
func (p *PathMap) pathStruct(k string) *PathStruct {
	p.once.Do(p.init)

	// If we don't care about casing, store everything lowercase
	if !p.CaseSensitive {
		k = strings.ToLower(k)
	}

	// If we have a salt, apply it before lookup
	if p.Salt != nil {
		k = p.Salt.SaltID(k)
	}

	return &PathStruct{
		Name:   fmt.Sprintf("map/%s/%s", p.Name, k),
		Schema: p.Schema,
	}
}

// Get reads a value out of the mapping
func (p *PathMap) Get(s logical.Storage, k string) (map[string]interface{}, error) {
	return p.pathStruct(k).Get(s)
}

// Put writes a value into the mapping
func (p *PathMap) Put(s logical.Storage, k string, v map[string]interface{}) error {
	return p.pathStruct(k).Put(s, v)
}

// Delete removes a value from the mapping
func (p *PathMap) Delete(s logical.Storage, k string) error {
	return p.pathStruct(k).Delete(s)
}

// List reads the keys under a given path
func (p *PathMap) List(s logical.Storage, prefix string) ([]string, error) {
	stripPrefix := fmt.Sprintf("struct/map/%s/", p.Name)
	fullPrefix := fmt.Sprintf("%s%s", stripPrefix, prefix)
	out, err := s.List(fullPrefix)
	if err != nil {
		return nil, err
	}
	stripped := make([]string, len(out))
	for idx, k := range out {
		stripped[idx] = strings.TrimPrefix(k, stripPrefix)
	}
	return stripped, nil
}

// Paths are the paths to append to the Backend paths.
func (p *PathMap) Paths() []*Path {
	p.once.Do(p.init)

	// Build the schema by simply adding the "key"
	schema := make(map[string]*FieldSchema)
	for k, v := range p.Schema {
		schema[k] = v
	}
	schema["key"] = &FieldSchema{
		Type:        TypeString,
		Description: fmt.Sprintf("Key for the %s mapping", p.Name),
	}

	return []*Path{
		&Path{
			Pattern: fmt.Sprintf("%s/%s$", p.Prefix, p.Name),

			Callbacks: map[logical.Operation]OperationFunc{
				logical.ListOperation: p.pathList,
				logical.ReadOperation: p.pathList,
			},

			HelpSynopsis: fmt.Sprintf("Read mappings for %s", p.Name),
		},

		&Path{
			Pattern: fmt.Sprintf(`%s/%s/(?P<key>[-\w]+)`, p.Prefix, p.Name),

			Fields: schema,

			Callbacks: map[logical.Operation]OperationFunc{
				logical.WriteOperation:  p.pathSingleWrite,
				logical.ReadOperation:   p.pathSingleRead,
				logical.DeleteOperation: p.pathSingleDelete,
			},

			HelpSynopsis: fmt.Sprintf("Read/write/delete a single %s mapping", p.Name),
		},
	}
}

func (p *PathMap) pathList(req *logical.Request, d *FieldData) (*logical.Response, error) {
	if p.Callbacks[logical.ListOperation] != nil {
		res, err := p.Callbacks[logical.ListOperation](req, d)
		if res != nil || err != nil {
			return res, err
		}
	}

	keys, err := req.Storage.List(req.Path)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(keys), nil
}

func (p *PathMap) pathSingleRead(req *logical.Request, d *FieldData) (*logical.Response, error) {
	if p.Callbacks[logical.ReadOperation] != nil {
		res, err := p.Callbacks[logical.ReadOperation](req, d)
		if res != nil || err != nil {
			return res, err
		}
	}

	v, err := p.Get(req.Storage, d.Get("key").(string))
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: v,
	}, nil
}

func (p *PathMap) pathSingleWrite(req *logical.Request, d *FieldData) (*logical.Response, error) {
	if p.Callbacks[logical.WriteOperation] != nil {
		res, err := p.Callbacks[logical.WriteOperation](req, d)
		if res != nil || err != nil {
			return res, err
		}
	}

	err := p.Put(req.Storage, d.Get("key").(string), d.Raw)
	return nil, err
}

func (p *PathMap) pathSingleDelete(req *logical.Request, d *FieldData) (*logical.Response, error) {
	if p.Callbacks[logical.DeleteOperation] != nil {
		res, err := p.Callbacks[logical.DeleteOperation](req, d)
		if res != nil || err != nil {
			return res, err
		}
	}

	err := p.Delete(req.Storage, d.Get("key").(string))
	return nil, err
}
