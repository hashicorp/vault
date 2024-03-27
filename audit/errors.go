package audit

import (
	"errors"
	"fmt"
)

var (
	ErrFilterParameter    = AuditErrorType{errors.New("filter parameter")}
	ErrFallbackParameter  = AuditErrorType{errors.New("fallback parameter")}
	ErrContextDone        = AuditErrorType{errors.New("context error")}
	ErrInvalidParameter   = AuditErrorType{errors.New("invalid parameter")}
	ErrEnterpriseOnly     = AuditErrorType{errors.New("enterprise-only")}
	ErrConfiguration      = AuditErrorType{errors.New("configuration error")}
	ErrConflict           = AuditErrorType{errors.New("audit conflict")}
	ErrUnknown            = AuditErrorType{errors.New("unknown error")}
	ErrPersistence        = AuditErrorType{errors.New("persistence error")}
	ErrBrokerRegistration = AuditErrorType{errors.New("registration error")}
)

type AuditErrorType struct {
	error
}

type AuditError struct {
	msg      string
	op       string
	err      AuditErrorType
	upstream error
}

// NewAuditError is used to create an AuditError which can be used to provide errors
// which are appropriate for internal or external consumption.
func NewAuditError(op string, msg string, err AuditErrorType) *AuditError {
	return &AuditError{
		op:  op,
		msg: msg,
		err: err,
	}
}

// SetUpstream should be used to configure the upstream error that prompted the
// creation of the AuditError.
// The original AuditError is returned after updating the upstream error.
func (e *AuditError) SetUpstream(err error) *AuditError {
	e.upstream = err

	return e
}

// TODO: PW: remove?
func (e *AuditError) Upstream() error {
	return e.upstream
}

func (e *AuditError) Downstream() AuditErrorType {
	return e.err
}

func (e *AuditError) Internal() error {
	err := e.upstream
	if err == nil {
		err = e.err
	}

	return fmt.Errorf("%s: %s: %w", e.op, e.msg, err)
}

func (e *AuditError) External() error {
	return fmt.Errorf("%s: %w", e.msg, e.err)
}

func (e *AuditError) Error() string {
	return e.Internal().Error()
}

func (e *AuditError) String() string {
	return e.Internal().Error()
}
