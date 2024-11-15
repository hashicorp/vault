// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wrapping

type WrapperType string

// These values define known types of Wrappers
const (
	WrapperTypeUnknown         WrapperType = "unknown"
	WrapperTypeAead            WrapperType = "aead"
	WrapperTypeAliCloudKms     WrapperType = "alicloudkms"
	WrapperTypeAwsKms          WrapperType = "awskms"
	WrapperTypeAzureKeyVault   WrapperType = "azurekeyvault"
	WrapperTypeGcpCkms         WrapperType = "gcpckms"
	WrapperTypeHsmAuto         WrapperType = "hsm-auto"
	WrapperTypeHuaweiCloudKms  WrapperType = "huaweicloudkms"
	WrapperTypeOciKms          WrapperType = "ocikms"
	WrapperTypePkcs11          WrapperType = "pkcs11"
	WrapperTypePooled          WrapperType = "pooled"
	WrapperTypeShamir          WrapperType = "shamir"
	WrapperTypeTencentCloudKms WrapperType = "tencentcloudkms"
	WrapperTypeTransit         WrapperType = "transit"
	WrapperTypeTest            WrapperType = "test-auto"
)

func (t WrapperType) String() string {
	return string(t)
}

type AeadType uint32

// These values define supported types of AEADs
const (
	AeadTypeUnknown AeadType = iota
	AeadTypeAesGcm
)

func (t AeadType) String() string {
	switch t {
	case AeadTypeAesGcm:
		return "aes-gcm"
	default:
		return "unknown"
	}
}

func AeadTypeMap(t string) AeadType {
	switch t {
	case "aes-gcm":
		return AeadTypeAesGcm
	default:
		return AeadTypeUnknown
	}
}

type HashType uint32

// These values define supported types of hashes
const (
	HashTypeUnknown HashType = iota
	HashTypeSha256
)

func (t HashType) String() string {
	switch t {
	case HashTypeSha256:
		return "sha256"
	default:
		return "unknown"
	}
}

func HashTypeMap(t string) HashType {
	switch t {
	case "sha256":
		return HashTypeSha256
	default:
		return HashTypeUnknown
	}
}
