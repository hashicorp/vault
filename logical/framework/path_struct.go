package framework

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

// PathStruct can be used to generate a path that stores a struct
// in the storage. This structure is a map[string]interface{} but the
// types are set according to the schema in this structure.
type PathStruct struct {
	Name            string
	Path            string
	Schema          map[string]*FieldSchema
	HelpSynopsis    string
	HelpDescription string

	Read bool
}

// Get reads the structure.
func (p *PathStruct) Get(s logical.Storage) (map[string]interface{}, error) {
	entry, err := s.Get(fmt.Sprintf("struct/%s", p.Name))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result map[string]interface{}
	if err := jsonutil.DecodeJSON(entry.Value, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Put writes the structure.
func (p *PathStruct) Put(s logical.Storage, v map[string]interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return s.Put(&logical.StorageEntry{
		Key:   fmt.Sprintf("struct/%s", p.Name),
		Value: bytes,
	})
}

// Delete removes the structure.
func (p *PathStruct) Delete(s logical.Storage) error {
	return s.Delete(fmt.Sprintf("struct/%s", p.Name))
}

// Paths are the paths to append to the Backend paths.
func (p *PathStruct) Paths() []*Path {
	// The single path we support to read/write this config
	path := &Path{
		Pattern: p.Path,
		Fields:  p.Schema,

		Callbacks: map[logical.Operation]OperationFunc{
			logical.CreateOperation: p.pathWrite,
			logical.UpdateOperation: p.pathWrite,
			logical.DeleteOperation: p.pathDelete,
		},

		ExistenceCheck: p.pathExistenceCheck,

		HelpSynopsis:    p.HelpSynopsis,
		HelpDescription: p.HelpDescription,
	}

	// If we support reads, add that
	if p.Read {
		path.Callbacks[logical.ReadOperation] = p.pathRead
	}

	return []*Path{path}
}

func (p *PathStruct) pathRead(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	v, err := p.Get(req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: v,
	}, nil
}

func (p *PathStruct) pathWrite(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	err := p.Put(req.Storage, d.Raw)
	return nil, err
}

func (p *PathStruct) pathDelete(
	req *logical.Request, d *FieldData) (*logical.Response, error) {
	err := p.Delete(req.Storage)
	return nil, err
}

func (p *PathStruct) pathExistenceCheck(
	req *logical.Request, d *FieldData) (bool, error) {
	v, err := p.Get(req.Storage)
	if err != nil {
		return false, err
	}

	return v != nil, nil
}
