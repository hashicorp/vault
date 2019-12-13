package service

import (
	"log"
	"time"

	"github.com/hashicorp/gokrb5/keytab"
	"github.com/hashicorp/gokrb5/types"
)

// Settings defines service side configuration settings.
type Settings struct {
	Keytab             *keytab.Keytab
	ktprinc            *types.PrincipalName
	sname              string
	requireHostAddr    bool
	disablePACDecoding bool
	cAddr              types.HostAddress
	maxClockSkew       time.Duration
	logger             *log.Logger
}

// NewSettings creates a new service Settings.
func NewSettings(kt *keytab.Keytab, settings ...func(*Settings)) *Settings {
	s := new(Settings)
	s.Keytab = kt
	for _, set := range settings {
		set(s)
	}
	return s
}

// RequireHostAddr used to configure service side to required host addresses to be specified in Kerberos tickets.
//
// s := NewSettings(kt, RequireHostAddr(true))
func RequireHostAddr(b bool) func(*Settings) {
	return func(s *Settings) {
		s.requireHostAddr = b
	}
}

// RequireHostAddr indicates if the service should require the host address to be included in the ticket.
func (s *Settings) RequireHostAddr() bool {
	return s.requireHostAddr
}

// DecodePAC used to configure service side to enable/disable PAC decoding if the PAC is present.
// Defaults to enabled if not specified.
//
// s := NewSettings(kt, DecodePAC(false))
func DecodePAC(b bool) func(*Settings) {
	return func(s *Settings) {
		s.disablePACDecoding = !b
	}
}

// DecodePAC indicates whether the service should decode any PAC information present in the ticket.
func (s *Settings) DecodePAC() bool {
	return !s.disablePACDecoding
}

// ClientAddress used to configure service side with the clients host address to be used during validation.
//
// s := NewSettings(kt, ClientAddress(h))
func ClientAddress(h types.HostAddress) func(*Settings) {
	return func(s *Settings) {
		s.cAddr = h
	}
}

// ClientAddress returns the client host address which has been provided to the service.
func (s *Settings) ClientAddress() types.HostAddress {
	return s.cAddr
}

// Logger used to configure service side with a logger.
//
// s := NewSettings(kt, Logger(l))
func Logger(l *log.Logger) func(*Settings) {
	return func(s *Settings) {
		s.logger = l
	}
}

// Logger returns the logger instances configured for the service. If none is configured nill will be returned.
func (s *Settings) Logger() *log.Logger {
	return s.logger
}

// KeytabPrincipal used to override the principal name used to find the key in the keytab.
//
// s := NewSettings(kt, KeytabPrincipal("someaccount"))
func KeytabPrincipal(p string) func(*Settings) {
	return func(s *Settings) {
		pn, _ := types.ParseSPNString(p)
		s.ktprinc = &pn
	}
}

// KeytabPrincipal returns the principal name used to find the key in the keytab if it has been overridden.
func (s *Settings) KeytabPrincipal() *types.PrincipalName {
	return s.ktprinc
}

// MaxClockSkew used to configure service side with the maximum acceptable clock skew
// between the service and the issue time of kerberos tickets
//
// s := NewSettings(kt, MaxClockSkew(d))
func MaxClockSkew(d time.Duration) func(*Settings) {
	return func(s *Settings) {
		s.maxClockSkew = d
	}
}

// MaxClockSkew returns the maximum acceptable clock skew between the service and the issue time of kerberos tickets.
// If none is defined a duration of 5 minutes is returned.
func (s *Settings) MaxClockSkew() time.Duration {
	if s.maxClockSkew.Nanoseconds() == 0 {
		return time.Duration(5) * time.Minute
	}
	return s.maxClockSkew
}

// SName used provide a specific service name to the service settings.
//
// s := NewSettings(kt, SName("HTTP/some.service.com"))
func SName(sname string) func(*Settings) {
	return func(s *Settings) {
		s.sname = sname
	}
}

// SName returns the specific service name to the service.
func (s *Settings) SName() string {
	return s.sname
}
