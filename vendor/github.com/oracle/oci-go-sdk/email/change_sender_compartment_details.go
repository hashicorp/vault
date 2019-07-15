// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Email Delivery API
//
// API for the Email Delivery service. Use this API to send high-volume, application-generated
// emails. For more information, see Overview of the Email Delivery Service (https://docs.cloud.oracle.com/iaas/Content/Email/Concepts/overview.htm).
//
// **Note:** Write actions (POST, UPDATE, DELETE) may take several minutes to propagate and be reflected by the API. If a subsequent read request fails to reflect your changes, wait a few minutes and try again.
//

package email

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ChangeSenderCompartmentDetails The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment
// into which the resource should be moved.
type ChangeSenderCompartmentDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment
	// into which the sender should be moved.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeSenderCompartmentDetails) String() string {
	return common.PointerString(m)
}
