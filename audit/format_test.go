package audit

import (
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

type noopFormatWriter struct {
	salt     *salt.Salt
	SaltFunc func() (*salt.Salt, error)
}

func (n *noopFormatWriter) WriteRequest(_ io.Writer, _ *AuditRequestEntry) error {
	return nil
}

func (n *noopFormatWriter) WriteResponse(_ io.Writer, _ *AuditResponseEntry) error {
	return nil
}

func (n *noopFormatWriter) Salt(ctx context.Context) (*salt.Salt, error) {
	if n.salt != nil {
		return n.salt, nil
	}
	var err error
	n.salt, err = salt.NewSalt(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	return n.salt, nil
}

func TestFormatRequestErrors(t *testing.T) {
	config := FormatterConfig{}
	formatter := AuditFormatter{
		AuditFormatWriter: &noopFormatWriter{},
	}

	if err := formatter.FormatRequest(context.Background(), ioutil.Discard, config, &LogInput{}); err == nil {
		t.Fatal("expected error due to nil request")
	}

	in := &LogInput{
		Request: &logical.Request{},
	}
	if err := formatter.FormatRequest(context.Background(), nil, config, in); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}

func TestFormatResponseErrors(t *testing.T) {
	config := FormatterConfig{}
	formatter := AuditFormatter{
		AuditFormatWriter: &noopFormatWriter{},
	}

	if err := formatter.FormatResponse(context.Background(), ioutil.Discard, config, &LogInput{}); err == nil {
		t.Fatal("expected error due to nil request")
	}

	in := &LogInput{
		Request: &logical.Request{},
	}
	if err := formatter.FormatResponse(context.Background(), nil, config, in); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}
