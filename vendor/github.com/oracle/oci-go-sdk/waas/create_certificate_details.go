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

// CreateCertificateDetails The data used to create a new SSL certificate.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateCertificateDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment in which to create the SSL certificate.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The data of the SSL certificate.
	CertificateData *string `mandatory:"true" json:"certificateData"`

	// The private key of the SSL certificate.
	PrivateKeyData *string `mandatory:"true" json:"privateKeyData"`

	// A user-friendly name for the SSL certificate. The name can be changed and does not need to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Set to true if the SSL certificate is self-signed.
	IsTrustVerificationDisabled *bool `mandatory:"false" json:"isTrustVerificationDisabled"`

	// A simple key-value pair without any defined schema.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// A key-value pair with a defined schema that restricts the values of tags. These predefined keys are scoped to namespaces.
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateCertificateDetails) String() string {
	return common.PointerString(m)
}
