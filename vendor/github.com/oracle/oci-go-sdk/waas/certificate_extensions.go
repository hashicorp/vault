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

// CertificateExtensions The representation of CertificateExtensions
type CertificateExtensions struct {
	Name *string `mandatory:"false" json:"name"`

	IsCritical *bool `mandatory:"false" json:"isCritical"`

	Value *string `mandatory:"false" json:"value"`
}

func (m CertificateExtensions) String() string {
	return common.PointerString(m)
}
