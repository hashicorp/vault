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

// CertificateSubjectName The representation of CertificateSubjectName
type CertificateSubjectName struct {
	Country *string `mandatory:"false" json:"country"`

	StateProvince *string `mandatory:"false" json:"stateProvince"`

	Locality *string `mandatory:"false" json:"locality"`

	Organization *string `mandatory:"false" json:"organization"`

	OrganizationalUnit *string `mandatory:"false" json:"organizationalUnit"`

	CommonName *string `mandatory:"false" json:"commonName"`

	EmailAddress *string `mandatory:"false" json:"emailAddress"`
}

func (m CertificateSubjectName) String() string {
	return common.PointerString(m)
}
