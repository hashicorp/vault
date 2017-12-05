package pb

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/vault/logical"
)

func ErrToString(e error) string {
	if e == nil {
		return ""
	}

	return e.Error()
}

func LogicalStorageEntryToProtoStorageEntry(e *logical.StorageEntry) *StorageEntry {
	return &StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoStorageEntryToLogicalStorageEntry(e *StorageEntry) *logical.StorageEntry {
	return &logical.StorageEntry{
		Key:      e.Key,
		Value:    e.Value,
		SealWrap: e.SealWrap,
	}
}

func ProtoLeaseOptionsToLogicalLeaseOptions(l *LeaseOptions) (logical.LeaseOptions, error) {
	t, err := ptypes.Timestamp(l.IssueTime)
	return logical.LeaseOptions{
		TTL:       time.Duration(l.TTL),
		Renewable: l.Renewable,
		Increment: time.Duration(l.Increment),
		IssueTime: t,
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
	}, err
}

func ProtoSecretToLogicalSecret(s *Secret) (*logical.Secret, error) {
	if s == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal(s.InternalData, &data)
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
		LeaseID:      s.LeaseId,
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
		InternalData: buf,
		LeaseId:      s.LeaseID,
	}, err
}

func LogicalRequestToProtoRequest(r *logical.Request) (*Request, error) {
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
		headers[k] = &Header{v}
	}

	return &Request{
		Id:                  r.ID,
		ReplicationCluster:  r.ReplicationCluster,
		Operation:           string(r.Operation),
		Path:                r.Path,
		Data:                buf,
		Secret:              secret,
		Auth:                auth,
		Headers:             headers,
		ClientToken:         r.ClientToken,
		ClientTokenAccessor: r.ClientTokenAccessor,
		DisplayName:         r.DisplayName,
		MountPoint:          r.MountPoint,
		MountType:           r.MountType,
		MountAccessor:       r.MountAccessor,
		//		WrapInfo: args.Request.WrapInfo,
		ClientTokenRemainingUses: int64(r.ClientTokenRemainingUses),
		//MFACreds: MFACreds,
		EntityId:        r.EntityID,
		PolicyOverride:  r.PolicyOverride,
		Unauthenticated: r.Unauthenticated,
	}, nil
}

func ProtoRequestToLogicalRequest(r *Request) (*logical.Request, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal(r.Data, &data)
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

	headers := map[string][]string{}
	for k, v := range r.Headers {
		headers[k] = v.Header
	}

	return &logical.Request{
		ID:                  r.Id,
		ReplicationCluster:  r.ReplicationCluster,
		Operation:           logical.Operation(r.Operation),
		Path:                r.Path,
		Data:                data,
		Secret:              secret,
		Auth:                auth,
		Headers:             headers,
		ClientToken:         r.ClientToken,
		ClientTokenAccessor: r.ClientTokenAccessor,
		DisplayName:         r.DisplayName,
		MountPoint:          r.MountPoint,
		MountType:           r.MountType,
		MountAccessor:       r.MountAccessor,
		//		WrapInfo: args.Request.WrapInfo,
		ClientTokenRemainingUses: int(r.ClientTokenRemainingUses),
		//MFACreds: MFACreds,
		EntityID:        r.EntityId,
		PolicyOverride:  r.PolicyOverride,
		Unauthenticated: r.Unauthenticated,
	}, nil
}

func ProtoResponseToLogicalResponse(r *Response) (*logical.Response, error) {
	secret, err := ProtoSecretToLogicalSecret(r.Secret)
	if err != nil {
		return nil, err
	}

	auth, err := ProtoAuthToLogicalAuth(r.Auth)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(r.Data, &data)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Secret:   secret,
		Auth:     auth,
		Data:     data,
		Redirect: r.Redirect,
		Warnings: r.Warnings,
		//	WrapInfo: r.WrapInfo,
	}, nil
}

func LogicalResponseToProtoResp(r *logical.Response) (*Response, error) {
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

	return &Response{
		Secret:   secret,
		Auth:     auth,
		Data:     buf,
		Redirect: r.Redirect,
		Warnings: r.Warnings,
		//	WrapInfo: r.WrapInfo,
	}, nil
}

func LogicalAliasToProtoAlias(a *logical.Alias) *Alias {
	return &Alias{
		MountType:     a.MountType,
		MountAccessor: a.MountAccessor,
		Name:          a.Name,
	}
}

func ProtoAliasToLogicalAlias(a *Alias) *logical.Alias {
	return &logical.Alias{
		MountType:     a.MountType,
		MountAccessor: a.MountAccessor,
		Name:          a.Name,
	}
}

func LogicalAuthToProtoAuth(a *logical.Auth) (*Auth, error) {
	if a == nil {
		return nil, nil
	}

	buf, err := json.Marshal(a.InternalData)
	if err != nil {
		return nil, err
	}

	groupAliases := make([]*Alias, len(a.GroupAliases))
	for i, al := range a.GroupAliases {
		groupAliases[i] = LogicalAliasToProtoAlias(al)
	}

	return &Auth{
		InternalData: buf,
		DisplayName:  a.DisplayName,
		Policies:     a.Policies,
		Metadata:     a.Metadata,
		ClientToken:  a.ClientToken,
		Accessor:     a.Accessor,
		Period:       int64(a.Period),
		NumUses:      int64(a.NumUses),
		EntityId:     a.EntityID,
		Alias:        LogicalAliasToProtoAlias(a.Alias),
		GroupAliases: groupAliases,
	}, nil
}

func ProtoAuthToLogicalAuth(a *Auth) (*logical.Auth, error) {
	if a == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	err := json.Unmarshal(a.InternalData, &data)
	if err != nil {
		return nil, err
	}

	groupAliases := make([]*logical.Alias, len(a.GroupAliases))
	for i, al := range a.GroupAliases {
		groupAliases[i] = ProtoAliasToLogicalAlias(al)
	}

	return &logical.Auth{
		InternalData: data,
		DisplayName:  a.DisplayName,
		Policies:     a.Policies,
		Metadata:     a.Metadata,
		ClientToken:  a.ClientToken,
		Accessor:     a.Accessor,
		Period:       time.Duration(a.Period),
		NumUses:      int(a.NumUses),
		EntityID:     a.EntityId,
		Alias:        ProtoAliasToLogicalAlias(a.Alias),
		GroupAliases: groupAliases,
	}, nil
}
