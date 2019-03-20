package mock

import (
	"github.com/hashicorp/vault/command/agent/sink"
)

type mockSink struct {
	token string
}

func NewSink(token string) sink.Sink {
	return &mockSink{
		token: token,
	}
}

func (m *mockSink) WriteToken(token string) error {
	m.token = token
	return nil
}

func (m *mockSink) Token() string {
	return m.token
}
