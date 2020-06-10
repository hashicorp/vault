package stepwise

import (
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestStepwise_makeRequest(t *testing.T) {

}

type mockEnvironment struct {
	s *httptest.Server
}

// func makeRequest(tt TestT, driver StepwiseEnvironment, step Step) (*api.Secret, error) {

func (m *mockEnvironment) Setup() error                  { return nil }
func (m *mockEnvironment) Client() (*api.Client, error)  { return nil, nil }
func (m *mockEnvironment) Teardown() error               { return nil }
func (m *mockEnvironment) Name() string                  { return "" }
func (m *mockEnvironment) ExpandPath(path string) string { return "" }
func (m *mockEnvironment) MountPath() string             { return "" }
func (m *mockEnvironment) RootToken() string             { return "" }
