package testhelpers

import (
	"github.com/hashicorp/vault/sdk/logical"
)

// RequestFactory helps write a common testing pattern in a concise way:
//
//	rf := testhelpers.NewRequestFactory()
//	b.HandleRequest(ctx, rf.Read("some/path"))
//	b.HandleRequest(ctx, rf.Update("some/path", map[string]interface{
//	    "param": "value",
//	})
//
// This automates all of the boilerplate of creating a logical.Storage, and
// setting that storage on every created request.
type RequestFactory struct {
	Storage logical.Storage
}

// NewRequestFactory creates a RequestFactory with the common default of an InmemStorage already in place.
// You should use it most of the time, although the
func NewRequestFactory() *RequestFactory {
	return &RequestFactory{Storage: &logical.InmemStorage{}}
}

func (r *RequestFactory) Read(path string) *logical.Request {
	return &logical.Request{
		Operation: logical.ReadOperation,
		Path:      path,
		Storage:   r.Storage,
	}
}

func (r *RequestFactory) List(path string) *logical.Request {
	return &logical.Request{
		Operation: logical.ListOperation,
		Path:      path,
		Storage:   r.Storage,
	}
}

func (r *RequestFactory) Delete(path string) *logical.Request {
	return &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      path,
		Storage:   r.Storage,
	}
}

func (r *RequestFactory) Create(path string, data map[string]interface{}) *logical.Request {
	return &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      path,
		Data:      data,
		Storage:   r.Storage,
	}
}

func (r *RequestFactory) Update(path string, data map[string]interface{}) *logical.Request {
	return &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      path,
		Data:      data,
		Storage:   r.Storage,
	}
}
