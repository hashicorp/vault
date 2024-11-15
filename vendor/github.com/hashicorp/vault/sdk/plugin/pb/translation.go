// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pb

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	ErrTypeUnknown uint32 = iota
	ErrTypeUserError
	ErrTypeInternalError
	ErrTypeCodedError
	ErrTypeStatusBadRequest
	ErrTypeUnsupportedOperation
	ErrTypeUnsupportedPath
	ErrTypeInvalidRequest
	ErrTypePermissionDenied
	ErrTypeMultiAuthzPending
	ErrTypeUnrecoverable
)

func ProtoErrToErr(e *ProtoError) error {
	if e == nil {
		return nil
	}

	var err error
	switch e.ErrType {
	case ErrTypeUnknown:
		err = errors.New(e.ErrMsg)
	case ErrTypeUserError:
		err = errutil.UserError{Err: e.ErrMsg}
	case ErrTypeInternalError:
		err = errutil.InternalError{Err: e.ErrMsg}
	case ErrTypeCodedError:
		err = logical.CodedError(int(e.ErrCode), e.ErrMsg)
	case ErrTypeStatusBadRequest:
		err = &logical.StatusBadRequest{Err: e.ErrMsg}
	case ErrTypeUnsupportedOperation:
		err = logical.ErrUnsupportedOperation
	case ErrTypeUnsupportedPath:
		err = logical.ErrUnsupportedPath
	case ErrTypeInvalidRequest:
		err = logical.ErrInvalidRequest
	case ErrTypePermissionDenied:
		err = logical.ErrPermissionDenied
	case ErrTypeMultiAuthzPending:
		err = logical.ErrMultiAuthzPending
	case ErrTypeUnrecoverable:
		err = logical.ErrUnrecoverable
	}

	return err
}

func ErrToProtoErr(e error) *ProtoError {
	if e == nil {
		return nil
	}
	pbErr := &ProtoError{
		ErrMsg:  e.Error(),
		ErrType: ErrTypeUnknown,
	}

	switch e.(type) {
	case errutil.UserError:
		pbErr.ErrType = ErrTypeUserError
	case errutil.InternalError:
		pbErr.ErrType = ErrTypeInternalError
	case logical.HTTPCodedError:
		pbErr.ErrType = ErrTypeCodedError
		pbErr.ErrCode = int64(e.(logical.HTTPCodedError).Code())
	case *logical.StatusBadRequest:
		pbErr.ErrType = ErrTypeStatusBadRequest
	}

	switch {
	case e == logical.ErrUnsupportedOperation:
		pbErr.ErrType = ErrTypeUnsupportedOperation
	case e == logical.ErrUnsupportedPath:
		pbErr.ErrType = ErrTypeUnsupportedPath
	case e == logical.ErrInvalidRequest:
		pbErr.ErrType = ErrTypeInvalidRequest
	case e == logical.ErrPermissionDenied:
		pbErr.ErrType = ErrTypePermissionDenied
	case e == logical.ErrMultiAuthzPending:
		pbErr.ErrType = ErrTypeMultiAuthzPending
	case e == logical.ErrUnrecoverable:
		pbErr.ErrType = ErrTypeUnrecoverable
	}

	return pbErr
}

func ErrToString(e error) string {
	if e == nil {
		return ""
	}

	return e.Error()
}

func LogicalStorageEntryToProtoStorageEntry(e *logical.StorageEntry) *StorageEntry {
	if e == nil {
		return nil
	}

	return &StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoStorageEntryToLogicalStorageEntry(e *StorageEntry) *logical.StorageEntry {
	if e == nil {
		return nil
	}

	return &logical.StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoLeaseOptionsToLogicalLeaseOptions(l *LeaseOptions) (logical.LeaseOptions, error) {
	if l == nil {
		return logical.LeaseOptions{}, nil
	}

	t, err := ptypes.Timestamp(l.IssueTime)
	return logical.LeaseOptions{
		TTL:       time.Duration(l.TTL),
		Renewable: l.Renewable,
		Increment: time.Duration(l.Increment),
		IssueTime: t,
		MaxTTL:    time.Duration(l.MaxTTL),
	}, err
}

func LogicalLeaseOptionsToProtoLeaseOptions(l logical.LeaseOptions) (*LeaseOptions, error) {
	t, err := ptypes.TimestampProto(l.IssueTime)
	if err != nil {
		return nil, err
	}

	return &LeaseOptions{
		TTL:       int64(l.TTL),
		Renewable: l.Renewable,
		Increment: int64(l.Increment),
		IssueTime: t,
		MaxTTL:    int64(l.MaxTTL),
	}, err
}

func ProtoSecretToLogicalSecret(s *Secret) (*logical.Secret, error) {
	if s == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(s.InternalData), &data)
	if err != nil {
		return nil, err
	}

	lease, err := ProtoLeaseOptionsToLogicalLeaseOptions(s.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &logical.Secret{
		LeaseOptions: lease,
		InternalData: data,
		LeaseID:      s.LeaseID,
	}, nil
}

func LogicalSecretToProtoSecret(s *logical.Secret) (*Secret, error) {
	if s == nil {
		return nil, nil
	}

	buf, err := json.Marshal(s.InternalData)
	if err != nil {
		return nil, err
	}

	lease, err := LogicalLeaseOptionsToProtoLeaseOptions(s.LeaseOptions)
	if err != nil {
		return nil, err
	}

	return &Secret{
		LeaseOptions: lease,
		InternalData: string(buf[:]),
		LeaseID:      s.LeaseID,
	}, err
}

func LogicalRequestToProtoRequest(r *logical.Request) (*Request, error) {
	if r == nil {
		return nil, nil
	}

	buf, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	secret, err := LogicalSecretToProtoSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := LogicalAuthToProtoAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	headers := map[string]*Header{}
	for k, v := range r.Headers {
		headers[k] = &Header{Header: v}
	}

	return &Request{
		ID:                       r.ID,
		ReplicationCluster:       r.ReplicationCluster,
		Operation:                string(r.Operation),
		Path:                     r.Path,
		Data:                     string(buf[:]),
		Secret:                   secret,
		Auth:                     auth,
		Headers:                  headers,
		ClientToken:              r.ClientToken,
		ClientTokenAccessor:      r.ClientTokenAccessor,
		DisplayName:              r.DisplayName,
		MountPoint:               r.MountPoint,
		MountType:                r.MountType,
		MountAccessor:            r.MountAccessor,
		WrapInfo:                 LogicalRequestWrapInfoToProtoRequestWrapInfo(r.WrapInfo),
		ClientTokenRemainingUses: int64(r.ClientTokenRemainingUses),
		Connection:               LogicalConnectionToProtoConnection(r.Connection),
		EntityID:                 r.EntityID,
		PolicyOverride:           r.PolicyOverride,
		Unauthenticated:          r.Unauthenticated,
	}, nil
}

func ProtoRequestToLogicalRequest(r *Request) (*logical.Request, error) {
	if r == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}

	secret, err := ProtoSecretToLogicalSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := ProtoAuthToLogicalAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	var headers map[string][]string
	if len(r.Headers) > 0 {
		headers = make(map[string][]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Header
		}
	}

	connection, err := ProtoConnectionToLogicalConnection(r.Connection)
	if err != nil {
		return nil, err
	}

	return &logical.Request{
		ID:                       r.ID,
		ReplicationCluster:       r.ReplicationCluster,
		Operation:                logical.Operation(r.Operation),
		Path:                     r.Path,
		Data:                     data,
		Secret:                   secret,
		Auth:                     auth,
		Headers:                  headers,
		ClientToken:              r.ClientToken,
		ClientTokenAccessor:      r.ClientTokenAccessor,
		DisplayName:              r.DisplayName,
		MountPoint:               r.MountPoint,
		MountType:                r.MountType,
		MountAccessor:            r.MountAccessor,
		WrapInfo:                 ProtoRequestWrapInfoToLogicalRequestWrapInfo(r.WrapInfo),
		ClientTokenRemainingUses: int(r.ClientTokenRemainingUses),
		Connection:               connection,
		EntityID:                 r.EntityID,
		PolicyOverride:           r.PolicyOverride,
		Unauthenticated:          r.Unauthenticated,
	}, nil
}

func LogicalConnectionToProtoConnection(c *logical.Connection) *Connection {
	if c == nil {
		return nil
	}

	return &Connection{
		RemoteAddr:      c.RemoteAddr,
		ConnectionState: TLSConnectionStateToProtoConnectionState(c.ConnState),
	}
}

func ProtoConnectionToLogicalConnection(c *Connection) (*logical.Connection, error) {
	if c == nil {
		return nil, nil
	}

	cs, err := ProtoConnectionStateToTLSConnectionState(c.ConnectionState)
	if err != nil {
		return nil, err
	}

	return &logical.Connection{
		RemoteAddr: c.RemoteAddr,
		ConnState:  cs,
	}, nil
}

func LogicalRequestWrapInfoToProtoRequestWrapInfo(i *logical.RequestWrapInfo) *RequestWrapInfo {
	if i == nil {
		return nil
	}

	return &RequestWrapInfo{
		TTL:      int64(i.TTL),
		Format:   i.Format,
		SealWrap: i.SealWrap,
	}
}

func ProtoRequestWrapInfoToLogicalRequestWrapInfo(i *RequestWrapInfo) *logical.RequestWrapInfo {
	if i == nil {
		return nil
	}

	return &logical.RequestWrapInfo{
		TTL:      time.Duration(i.TTL),
		Format:   i.Format,
		SealWrap: i.SealWrap,
	}
}

func ProtoResponseToLogicalResponse(r *Response) (*logical.Response, error) {
	if r == nil {
		return nil, nil
	}

	secret, err := ProtoSecretToLogicalSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := ProtoAuthToLogicalAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(r.Data), &data)
	if err != nil {
		return nil, err
	}

	wrapInfo, err := ProtoResponseWrapInfoToLogicalResponseWrapInfo(r.WrapInfo)
	if err != nil {
		return nil, err
	}

	var headers map[string][]string
	if len(r.Headers) > 0 {
		headers = make(map[string][]string, len(r.Headers))
		for k, v := range r.Headers {
			headers[k] = v.Header
		}
	}

	return &logical.Response{
		Secret:    secret,
		Auth:      auth,
		Data:      data,
		Redirect:  r.Redirect,
		Warnings:  r.Warnings,
		WrapInfo:  wrapInfo,
		Headers:   headers,
		MountType: r.MountType,
	}, nil
}

func ProtoResponseWrapInfoToLogicalResponseWrapInfo(i *ResponseWrapInfo) (*wrapping.ResponseWrapInfo, error) {
	if i == nil {
		return nil, nil
	}

	t, err := ptypes.Timestamp(i.CreationTime)
	if err != nil {
		return nil, err
	}

	return &wrapping.ResponseWrapInfo{
		TTL:             time.Duration(i.TTL),
		Token:           i.Token,
		Accessor:        i.Accessor,
		CreationTime:    t,
		WrappedAccessor: i.WrappedAccessor,
		WrappedEntityID: i.WrappedEntityID,
		Format:          i.Format,
		CreationPath:    i.CreationPath,
		SealWrap:        i.SealWrap,
	}, nil
}

func LogicalResponseWrapInfoToProtoResponseWrapInfo(i *wrapping.ResponseWrapInfo) (*ResponseWrapInfo, error) {
	if i == nil {
		return nil, nil
	}

	t, err := ptypes.TimestampProto(i.CreationTime)
	if err != nil {
		return nil, err
	}

	return &ResponseWrapInfo{
		TTL:             int64(i.TTL),
		Token:           i.Token,
		Accessor:        i.Accessor,
		CreationTime:    t,
		WrappedAccessor: i.WrappedAccessor,
		WrappedEntityID: i.WrappedEntityID,
		Format:          i.Format,
		CreationPath:    i.CreationPath,
		SealWrap:        i.SealWrap,
	}, nil
}

func LogicalResponseToProtoResponse(r *logical.Response) (*Response, error) {
	if r == nil {
		return nil, nil
	}

	secret, err := LogicalSecretToProtoSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := LogicalAuthToProtoAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	wrapInfo, err := LogicalResponseWrapInfoToProtoResponseWrapInfo(r.WrapInfo)
	if err != nil {
		return nil, err
	}

	headers := map[string]*Header{}
	for k, v := range r.Headers {
		headers[k] = &Header{Header: v}
	}

	return &Response{
		Secret:    secret,
		Auth:      auth,
		Data:      string(buf[:]),
		Redirect:  r.Redirect,
		Warnings:  r.Warnings,
		WrapInfo:  wrapInfo,
		Headers:   headers,
		MountType: r.MountType,
	}, nil
}

func LogicalAuthToProtoAuth(a *logical.Auth) (*Auth, error) {
	if a == nil {
		return nil, nil
	}

	buf, err := json.Marshal(a.InternalData)
	if err != nil {
		return nil, err
	}

	lo, err := LogicalLeaseOptionsToProtoLeaseOptions(a.LeaseOptions)
	if err != nil {
		return nil, err
	}

	boundCIDRs := make([]string, len(a.BoundCIDRs))
	for i, cidr := range a.BoundCIDRs {
		boundCIDRs[i] = cidr.String()
	}

	return &Auth{
		LeaseOptions:     lo,
		TokenType:        uint32(a.TokenType),
		InternalData:     string(buf[:]),
		DisplayName:      a.DisplayName,
		Policies:         a.Policies,
		TokenPolicies:    a.TokenPolicies,
		IdentityPolicies: a.IdentityPolicies,
		NoDefaultPolicy:  a.NoDefaultPolicy,
		Metadata:         a.Metadata,
		ClientToken:      a.ClientToken,
		Accessor:         a.Accessor,
		Period:           int64(a.Period),
		NumUses:          int64(a.NumUses),
		EntityID:         a.EntityID,
		Alias:            a.Alias,
		GroupAliases:     a.GroupAliases,
		BoundCIDRs:       boundCIDRs,
		ExplicitMaxTTL:   int64(a.ExplicitMaxTTL),
	}, nil
}

func ProtoAuthToLogicalAuth(a *Auth) (*logical.Auth, error) {
	if a == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(a.InternalData), &data)
	if err != nil {
		return nil, err
	}

	lo, err := ProtoLeaseOptionsToLogicalLeaseOptions(a.LeaseOptions)
	if err != nil {
		return nil, err
	}

	boundCIDRs, err := parseutil.ParseAddrs(a.BoundCIDRs)
	if err != nil {
		return nil, err
	}
	if len(boundCIDRs) == 0 {
		// On inbound auths, if auth.BoundCIDRs is empty, it will be nil.
		// Let's match that behavior outbound.
		boundCIDRs = nil
	}

	return &logical.Auth{
		LeaseOptions:     lo,
		TokenType:        logical.TokenType(a.TokenType),
		InternalData:     data,
		DisplayName:      a.DisplayName,
		Policies:         a.Policies,
		TokenPolicies:    a.TokenPolicies,
		IdentityPolicies: a.IdentityPolicies,
		NoDefaultPolicy:  a.NoDefaultPolicy,
		Metadata:         a.Metadata,
		ClientToken:      a.ClientToken,
		Accessor:         a.Accessor,
		Period:           time.Duration(a.Period),
		NumUses:          int(a.NumUses),
		EntityID:         a.EntityID,
		Alias:            a.Alias,
		GroupAliases:     a.GroupAliases,
		BoundCIDRs:       boundCIDRs,
		ExplicitMaxTTL:   time.Duration(a.ExplicitMaxTTL),
	}, nil
}

func LogicalTokenEntryToProtoTokenEntry(t *logical.TokenEntry) *TokenEntry {
	if t == nil {
		return nil
	}

	boundCIDRs := make([]string, len(t.BoundCIDRs))
	for i, cidr := range t.BoundCIDRs {
		boundCIDRs[i] = cidr.String()
	}

	return &TokenEntry{
		ID:                 t.ID,
		Accessor:           t.Accessor,
		Parent:             t.Parent,
		Policies:           t.Policies,
		InlinePolicy:       t.InlinePolicy,
		Path:               t.Path,
		Meta:               t.Meta,
		InternalMeta:       t.InternalMeta,
		DisplayName:        t.DisplayName,
		NumUses:            int64(t.NumUses),
		CreationTime:       t.CreationTime,
		TTL:                int64(t.TTL),
		ExplicitMaxTTL:     int64(t.ExplicitMaxTTL),
		Role:               t.Role,
		Period:             int64(t.Period),
		EntityID:           t.EntityID,
		NoIdentityPolicies: t.NoIdentityPolicies,
		BoundCIDRs:         boundCIDRs,
		NamespaceID:        t.NamespaceID,
		CubbyholeID:        t.CubbyholeID,
		Type:               uint32(t.Type),
		ExternalID:         t.ExternalID,
	}
}

func ProtoTokenEntryToLogicalTokenEntry(t *TokenEntry) (*logical.TokenEntry, error) {
	if t == nil {
		return nil, nil
	}

	boundCIDRs, err := parseutil.ParseAddrs(t.BoundCIDRs)
	if err != nil {
		return nil, err
	}
	if len(boundCIDRs) == 0 {
		// On inbound auths, if auth.BoundCIDRs is empty, it will be nil.
		// Let's match that behavior outbound.
		boundCIDRs = nil
	}

	return &logical.TokenEntry{
		ID:                 t.ID,
		Accessor:           t.Accessor,
		Parent:             t.Parent,
		Policies:           t.Policies,
		InlinePolicy:       t.InlinePolicy,
		Path:               t.Path,
		Meta:               t.Meta,
		InternalMeta:       t.InternalMeta,
		DisplayName:        t.DisplayName,
		NumUses:            int(t.NumUses),
		CreationTime:       t.CreationTime,
		TTL:                time.Duration(t.TTL),
		ExplicitMaxTTL:     time.Duration(t.ExplicitMaxTTL),
		Role:               t.Role,
		Period:             time.Duration(t.Period),
		EntityID:           t.EntityID,
		NoIdentityPolicies: t.NoIdentityPolicies,
		BoundCIDRs:         boundCIDRs,
		NamespaceID:        t.NamespaceID,
		CubbyholeID:        t.CubbyholeID,
		Type:               logical.TokenType(t.Type),
		ExternalID:         t.ExternalID,
	}, nil
}

func TLSConnectionStateToProtoConnectionState(connState *tls.ConnectionState) *ConnectionState {
	if connState == nil {
		return nil
	}

	var verifiedChains []*CertificateChain

	if lvc := len(connState.VerifiedChains); lvc > 0 {
		verifiedChains = make([]*CertificateChain, lvc)
		for i, vc := range connState.VerifiedChains {
			verifiedChains[i] = CertificateChainToProtoCertificateChain(vc)
		}
	}

	return &ConnectionState{
		Version:                     uint32(connState.Version),
		HandshakeComplete:           connState.HandshakeComplete,
		DidResume:                   connState.DidResume,
		CipherSuite:                 uint32(connState.CipherSuite),
		NegotiatedProtocol:          connState.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  connState.NegotiatedProtocolIsMutual,
		ServerName:                  connState.ServerName,
		PeerCertificates:            CertificateChainToProtoCertificateChain(connState.PeerCertificates),
		VerifiedChains:              verifiedChains,
		SignedCertificateTimestamps: connState.SignedCertificateTimestamps,
		OcspResponse:                connState.OCSPResponse,
		TlsUnique:                   connState.TLSUnique,
	}
}

func ProtoConnectionStateToTLSConnectionState(cs *ConnectionState) (*tls.ConnectionState, error) {
	if cs == nil {
		return nil, nil
	}

	var (
		err              error
		peerCertificates []*x509.Certificate
		verifiedChains   [][]*x509.Certificate
	)

	if peerCertificates, err = ProtoCertificateChainToCertificateChain(cs.PeerCertificates); err != nil {
		return nil, err
	}

	if lvc := len(cs.VerifiedChains); lvc > 0 {
		verifiedChains = make([][]*x509.Certificate, lvc)
		for i, vc := range cs.VerifiedChains {
			if verifiedChains[i], err = ProtoCertificateChainToCertificateChain(vc); err != nil {
				return nil, err
			}
		}
	}

	connState := &tls.ConnectionState{
		Version:                     uint16(cs.Version),
		HandshakeComplete:           cs.HandshakeComplete,
		DidResume:                   cs.DidResume,
		CipherSuite:                 uint16(cs.CipherSuite),
		NegotiatedProtocol:          cs.NegotiatedProtocol,
		NegotiatedProtocolIsMutual:  cs.NegotiatedProtocolIsMutual,
		ServerName:                  cs.ServerName,
		PeerCertificates:            peerCertificates,
		VerifiedChains:              verifiedChains,
		SignedCertificateTimestamps: cs.SignedCertificateTimestamps,
		OCSPResponse:                cs.OcspResponse,
		TLSUnique:                   cs.TlsUnique,
	}

	return connState, nil
}

func CertificateChainToProtoCertificateChain(chain []*x509.Certificate) *CertificateChain {
	if len(chain) == 0 {
		return nil
	}

	cc := &CertificateChain{Certificates: make([]*Certificate, len(chain))}

	for i, c := range chain {
		cc.Certificates[i] = X509CertificateToProtoCertificate(c)
	}

	return cc
}

func ProtoCertificateChainToCertificateChain(cc *CertificateChain) ([]*x509.Certificate, error) {
	if cc == nil || len(cc.Certificates) == 0 {
		return nil, nil
	}

	certs := make([]*x509.Certificate, len(cc.Certificates))

	for i, c := range cc.Certificates {
		var err error
		if certs[i], err = ProtoCertificateToX509Certificate(c); err != nil {
			return nil, err
		}
	}

	return certs, nil
}

func X509CertificateToProtoCertificate(cert *x509.Certificate) *Certificate {
	if cert == nil {
		return nil
	}

	return &Certificate{Asn1Data: cert.Raw}
}

func ProtoCertificateToX509Certificate(c *Certificate) (*x509.Certificate, error) {
	if c == nil {
		return nil, nil
	}

	return x509.ParseCertificate(c.Asn1Data)
}
