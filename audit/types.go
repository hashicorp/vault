package audit

import (
	"os"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// Audit subtypes.
const (
	RequestType  subtype = "AuditRequest"
	ResponseType subtype = "AuditResponse"
)

// Audit formats.
const (
	JSONFormat  format = "json"
	JSONxFormat format = "jsonx"
)

// version defines the version of audit events.
const version = "v0.1"

// subtype defines the type of audit event.
type subtype string

// format defines types of format audit events support.
type format string

// auditEvent is the audit event.
type auditEvent struct {
	ID             string            `json:"id"`
	Version        string            `json:"version"`
	Subtype        subtype           `json:"subtype"` // the subtype of the audit event.
	Timestamp      time.Time         `json:"timestamp"`
	Data           *logical.LogInput `json:"data"`
	RequiredFormat format            `json:"format"`
}

type Option func(*options) error

type options struct {
	withID          string
	withNow         time.Time
	withSubtype     subtype
	withFormat      format
	withFileMode    *os.FileMode
	withPrefix      string
	withFacility    string
	withTag         string
	withSocketType  string
	withMaxDuration time.Duration
}
