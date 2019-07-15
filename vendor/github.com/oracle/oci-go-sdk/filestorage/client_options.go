// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// File Storage Service API
//
// The API for the File Storage Service.
//

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ClientOptions NFS export options applied to a specified set of
// clients. Only governs access through the associated
// export. Access to the same file system through a different
// export (on the same or different mount target) will be governed
// by that export's export options.
type ClientOptions struct {

	// Clients these options should apply to. Must be a either
	// single IPv4 address or single IPv4 CIDR block.
	// **Note:** Access will also be limited by any applicable VCN
	// security rules and the ability to route IP packets to the
	// mount target. Mount targets do not have Internet-routable IP addresses.
	Source *string `mandatory:"true" json:"source"`

	// If `true`, clients accessing the file system through this
	// export must connect from a privileged source port. If
	// unspecified, defaults to `true`.
	RequirePrivilegedSourcePort *bool `mandatory:"false" json:"requirePrivilegedSourcePort"`

	// Type of access to grant clients using the file system
	// through this export. If unspecified defaults to `READ_ONLY`.
	Access ClientOptionsAccessEnum `mandatory:"false" json:"access,omitempty"`

	// Used when clients accessing the file system through this export
	// have their UID and GID remapped to 'anonymousUid' and
	// 'anonymousGid'. If `ALL`, all users and groups are remapped;
	// if `ROOT`, only the root user and group (UID/GID 0) are
	// remapped; if `NONE`, no remapping is done. If unspecified,
	// defaults to `ROOT`.
	IdentitySquash ClientOptionsIdentitySquashEnum `mandatory:"false" json:"identitySquash,omitempty"`

	// UID value to remap to when squashing a client UID (see
	// identitySquash for more details.) If unspecified, defaults
	// to `65534`.
	AnonymousUid *int64 `mandatory:"false" json:"anonymousUid"`

	// GID value to remap to when squashing a client GID (see
	// identitySquash for more details.) If unspecified defaults
	// to `65534`.
	AnonymousGid *int64 `mandatory:"false" json:"anonymousGid"`
}

func (m ClientOptions) String() string {
	return common.PointerString(m)
}

// ClientOptionsAccessEnum Enum with underlying type: string
type ClientOptionsAccessEnum string

// Set of constants representing the allowable values for ClientOptionsAccessEnum
const (
	ClientOptionsAccessWrite ClientOptionsAccessEnum = "READ_WRITE"
	ClientOptionsAccessOnly  ClientOptionsAccessEnum = "READ_ONLY"
)

var mappingClientOptionsAccess = map[string]ClientOptionsAccessEnum{
	"READ_WRITE": ClientOptionsAccessWrite,
	"READ_ONLY":  ClientOptionsAccessOnly,
}

// GetClientOptionsAccessEnumValues Enumerates the set of values for ClientOptionsAccessEnum
func GetClientOptionsAccessEnumValues() []ClientOptionsAccessEnum {
	values := make([]ClientOptionsAccessEnum, 0)
	for _, v := range mappingClientOptionsAccess {
		values = append(values, v)
	}
	return values
}

// ClientOptionsIdentitySquashEnum Enum with underlying type: string
type ClientOptionsIdentitySquashEnum string

// Set of constants representing the allowable values for ClientOptionsIdentitySquashEnum
const (
	ClientOptionsIdentitySquashNone ClientOptionsIdentitySquashEnum = "NONE"
	ClientOptionsIdentitySquashRoot ClientOptionsIdentitySquashEnum = "ROOT"
	ClientOptionsIdentitySquashAll  ClientOptionsIdentitySquashEnum = "ALL"
)

var mappingClientOptionsIdentitySquash = map[string]ClientOptionsIdentitySquashEnum{
	"NONE": ClientOptionsIdentitySquashNone,
	"ROOT": ClientOptionsIdentitySquashRoot,
	"ALL":  ClientOptionsIdentitySquashAll,
}

// GetClientOptionsIdentitySquashEnumValues Enumerates the set of values for ClientOptionsIdentitySquashEnum
func GetClientOptionsIdentitySquashEnumValues() []ClientOptionsIdentitySquashEnum {
	values := make([]ClientOptionsIdentitySquashEnum, 0)
	for _, v := range mappingClientOptionsIdentitySquash {
		values = append(values, v)
	}
	return values
}
