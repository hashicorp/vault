package framework

import (
	"context"
	"fmt"
	"strings"
	"sync"

	saltpkg "github.com/hashicorp/vault/helper/salt"
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
	Salt          *saltpkg.Salt
	SaltFunc      func(context.Context) (*saltpkg.Salt, error)

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
func (p *PathMap) pathStruct(ctx context.Context, s logical.Storage, k string) (*PathStruct, error) {
	p.once.Do(p.init)

	// If we don't care about casing, store everything lowercase
	if !p.CaseSensitive {
		k = strings.ToLower(k)
	}

	// The original key before any salting
	origKey := k

	// If we have a salt, apply it before lookup
	salt := p.Salt
	var err error
	if p.SaltFunc != nil {
		salt, err = p.SaltFunc(ctx)
		if err != nil {
			return nil, err
		}
	}
	if salt != nil {
		k = "s" + salt.SaltIDHashFunc(k, saltpkg.SHA256Hash)
	}

	finalName := fmt.Sprintf("map/%s/%s", p.Name, k)
	ps := &PathStruct{
		Name:   finalName,
		Schema: p.Schema,
	}

	if !strings.HasPrefix(origKey, "s") && k != origKey {
		// Ensure that no matter what happens what is returned is the final
		// path
		defer func() {
			ps.Name = finalName
		}()

		//
		// Check for unsalted version and upgrade if so
		//

		// Generate the unsalted name
		unsaltedName := fmt.Sprintf("map/%s/%s", p.Name, origKey)
		// Set the path struct to use the unsalted name
		ps.Name = unsaltedName

		val, err := ps.Get(ctx, s)
		if err != nil {
			return nil, err
		}
		// If not nil, we have an unsalted entry -- upgrade it
		if val != nil {
			// Set the path struct to use the desired final name
			ps.Name = finalName
			err = ps.Put(ctx, s, val)
			if err != nil {
				return nil, err
			}
			// Set it back to the old path and delete
			ps.Name = unsaltedName
			err = ps.Delete(ctx, s)
			if err != nil {
				return nil, err
			}
			// We'll set this in the deferred function but doesn't hurt here
			ps.Name = finalName
		}

		//
		// Check for SHA1 hashed version and upgrade if so
		//

		// Generate the SHA1 hash suffixed path name
		sha1SuffixedName := fmt.Sprintf("map/%s/%s", p.Name, salt.SaltID(origKey))

		// Set the path struct to use the SHA1 hash suffixed path name
		ps.Name = sha1SuffixedName

		val, err = ps.Get(ctx, s)
		if err != nil {
			return nil, err
		}
		// If not nil, we have an SHA1 hash suffixed entry -- upgrade it
		if val != nil {
			// Set the path struct to use the desired final name
			ps.Name = finalName
			err = ps.Put(ctx, s, val)
			if err != nil {
				return nil, err
			}
			// Set it back to the old path and delete
			ps.Name = sha1SuffixedName
			err = ps.Delete(ctx, s)
			if err != nil {
				return nil, err
			}
			// We'll set this in the deferred function but doesn't hurt here
			ps.Name = finalName
		}
	}

	return ps, nil
}

// Get reads a value out of the mapping
func (p *PathMap) Get(ctx context.Context, s logical.Storage, k string) (map[string]interface{}, error) {
	ps, err := p.pathStruct(ctx, s, k)
	if err != nil {
		return nil, err
	}
	return ps.Get(ctx, s)
}

// Put writes a value into the mapping
func (p *PathMap) Put(ctx context.Context, s logical.Storage, k string, v map[string]interface{}) error {
	ps, err := p.pathStruct(ctx, s, k)
	if err != nil {
		return err
	}
	return ps.Put(ctx, s, v)
}

// Delete removes a value from the mapping
func (p *PathMap) Delete(ctx context.Context, s logical.Storage, k string) error {
	ps, err := p.pathStruct(ctx, s, k)
	if err != nil {
		return err
	}
	return ps.Delete(ctx, s)
}

// List reads the keys under a given path
func (p *PathMap) List(ctx context.Context, s logical.Storage, prefix string) ([]string, error) {
	stripPrefix := fmt.Sprintf("struct/map/%s/", p.Name)
	fullPrefix := fmt.Sprintf("%s%s", stripPrefix, prefix)
	out, err := s.List(ctx, fullPrefix)
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
			Pattern: fmt.Sprintf("%s/%s/?$", p.Prefix, p.Name),

			Callbacks: map[logical.Operation]OperationFunc{
				logical.ListOperation: p.pathList(),
				logical.ReadOperation: p.pathList(),
			},

			HelpSynopsis: fmt.Sprintf("Read mappings for %s", p.Name),
		},

		&Path{
			Pattern: fmt.Sprintf(`%s/%s/(?P<key>[-\w]+)`, p.Prefix, p.Name),

			Fields: schema,

			Callbacks: map[logical.Operation]OperationFunc{
				logical.CreateOperation: p.pathSingleWrite(),
				logical.ReadOperation:   p.pathSingleRead(),
				logical.UpdateOperation: p.pathSingleWrite(),
				logical.DeleteOperation: p.pathSingleDelete(),
			},

			HelpSynopsis: fmt.Sprintf("Read/write/delete a single %s mapping", p.Name),

			ExistenceCheck: p.pathSingleExistenceCheck(),
		},
	}
}

func (p *PathMap) pathList() OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *FieldData) (*logical.Response, error) {
		keys, err := p.List(ctx, req.Storage, "")
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(keys), nil
	}
}

func (p *PathMap) pathSingleRead() OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *FieldData) (*logical.Response, error) {
		v, err := p.Get(ctx, req.Storage, d.Get("key").(string))
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: v,
		}, nil
	}
}

func (p *PathMap) pathSingleWrite() OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *FieldData) (*logical.Response, error) {
		err := p.Put(ctx, req.Storage, d.Get("key").(string), d.Raw)
		return nil, err
	}
}

func (p *PathMap) pathSingleDelete() OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *FieldData) (*logical.Response, error) {
		err := p.Delete(ctx, req.Storage, d.Get("key").(string))
		return nil, err
	}
}

func (p *PathMap) pathSingleExistenceCheck() ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, d *FieldData) (bool, error) {
		v, err := p.Get(ctx, req.Storage, d.Get("key").(string))
		if err != nil {
			return false, err
		}
		return v != nil, nil
	}
}
