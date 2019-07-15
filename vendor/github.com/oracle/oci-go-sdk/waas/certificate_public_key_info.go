// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CertificatePublicKeyInfo The representation of CertificatePublicKeyInfo
type CertificatePublicKeyInfo struct {
	Algorithm *string `mandatory:"false" json:"algorithm"`

	Exponent *int `mandatory:"false" json:"exponent"`

	KeySize *int `mandatory:"false" json:"keySize"`
}

func (m CertificatePublicKeyInfo) String() string {
	return common.PointerString(m)
}
