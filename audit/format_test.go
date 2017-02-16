package audit

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

type noopFormatWriter struct {
}

func (n *noopFormatWriter) WriteRequest(_ io.Writer, _ *AuditRequestEntry) error {
	return nil
}

func (n *noopFormatWriter) WriteResponse(_ io.Writer, _ *AuditResponseEntry) error {
	return nil
}

func TestFormatRequestErrors(t *testing.T) {
	salter, _ := salt.NewSalt(nil, nil)
	config := FormatterConfig{
		Salt: salter,
	}
	formatter := AuditFormatter{
		AuditFormatWriter: &noopFormatWriter{},
	}

	if err := formatter.FormatRequest(ioutil.Discard, config, nil, nil, nil); err == nil {
		t.Fatal("expected error due to nil request")
	}
	if err := formatter.FormatRequest(nil, config, nil, &logical.Request{}, nil); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}

func TestFormatResponseErrors(t *testing.T) {
	salter, _ := salt.NewSalt(nil, nil)
	config := FormatterConfig{
		Salt: salter,
	}
	formatter := AuditFormatter{
		AuditFormatWriter: &noopFormatWriter{},
	}

	if err := formatter.FormatResponse(ioutil.Discard, config, nil, nil, nil, nil); err == nil {
		t.Fatal("expected error due to nil request")
	}
	if err := formatter.FormatResponse(nil, config, nil, &logical.Request{}, nil, nil); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}
