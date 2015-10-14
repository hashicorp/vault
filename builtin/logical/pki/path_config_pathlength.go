package pki

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigPathLength(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/pathlength",
		Fields: map[string]*framework.FieldSchema{
			"max_path_length": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: `The maximum allowable path length`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation:  b.pathWritePathLength,
			logical.ReadOperation:   b.pathReadPathLength,
			logical.DeleteOperation: b.pathDeletePathLength,
		},

		HelpSynopsis:    pathConfigPathLengthHelpSyn,
		HelpDescription: pathConfigPathLengthHelpDesc,
	}
}

func getPathLengthEntry(req *logical.Request) (*pathLengthEntry, error) {
	entry, err := req.Storage.Get("pathlength")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var pathLen pathLengthEntry
	if err := entry.DecodeJSON(&pathLen); err != nil {
		return nil, err
	}

	return &pathLen, nil
}

func writePathLengthEntry(req *logical.Request, pathLen *pathLengthEntry) error {
	entry, err := logical.StorageEntryJSON("pathlength", pathLen)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("Unable to marshal entry into JSON")
	}

	err = req.Storage.Put(entry)
	if err != nil {
		return err
	}

	return nil
}

func (b *backend) pathReadPathLength(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	pathLenEntry, err := getPathLengthEntry(req)
	if err != nil {
		return nil, err
	}
	if pathLenEntry == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(pathLenEntry).Map(),
	}

	return resp, nil
}

func (b *backend) pathDeletePathLength(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("pathlength")
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathWritePathLength(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	maxLenInt, ok := data.GetOk("max_path_length")
	if !ok {
		return logical.ErrorResponse("max_path_length must be specified"), nil
	}
	maxLen, ok := maxLenInt.(int)
	if !ok {
		return logical.ErrorResponse("max_path_length value could not be parsed as an integer"), nil
	}
	if maxLen < 0 {
		return logical.ErrorResponse("max_path_length value must be greater than or equal to 0"), nil
	}

	pathLenEntry := &pathLengthEntry{
		MaxPathLength: maxLen,
	}

	return nil, writePathLengthEntry(req, pathLenEntry)
}

type pathLengthEntry struct {
	MaxPathLength int `json:"max_path_length" structs:"max_path_length" mapstructure:"max_path_length"`
}

const pathConfigPathLengthHelpSyn = `
Set the maximum path length for issued certificates.
`

const pathConfigPathLengthHelpDesc = `
This path allows you to set the maximum path length for issued or signed
certificates. The value must be greater than or equal to zero. If a value
is set and an unlimited path length is required, simply delete the value.
`
