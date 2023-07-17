package audit

import (
	"os"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/vault/sdk/logical"
)

// Audit subtypes.
const (
	AuditRequestType  auditSubtype = "AuditRequest"
	AuditResponseType auditSubtype = "AuditResponse"
)

// Audit formats.
const (
	AuditFormatJSON  auditFormat = "json"
	AuditFormatJSONx auditFormat = "jsonx"
)

// auditVersion defines the version of audit events.
const auditVersion = "v0.1"

// auditSubtype defines the type of audit event.
type auditSubtype string

// auditFormat defines types of format audit events support.
type auditFormat string

// audit is the audit event.
type audit struct {
	ID             string            `json:"id"`
	Version        string            `json:"version"`
	Subtype        auditSubtype      `json:"subtype"` // the subtype of the audit event.
	Timestamp      time.Time         `json:"timestamp"`
	Data           *logical.LogInput `json:"data"`
	RequiredFormat auditFormat       `json:"format"`
}

type AuditOption func(*AuditOptions) error

type AuditOptions struct {
	event.Options
	withSubtype     auditSubtype
	withFormat      auditFormat
	withFileMode    *os.FileMode
	withPrefix      string
	withFacility    string
	withTag         string
	withSocketType  string
	withMaxDuration time.Duration
}
