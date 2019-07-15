// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Helper methods for Oracle Cloud Infrastructure Go SDK Samples
//

package helpers

import (
	"os"

	"github.com/oracle/oci-go-sdk/common"
)

var (
	availabilityDomain string
	compartmentID      string
	rootCompartmentID  string
)

// ParseEnvironmentVariables parse shared variables from environment variables, other samples should define their own
// viariables and call this function to initialize shared variables
func ParseEnvironmentVariables() {
	availabilityDomain = os.Getenv("OCI_AVAILABILITY_DOMAIN")
	compartmentID = os.Getenv("OCI_COMPARTMENT_ID")
	rootCompartmentID = os.Getenv("OCI_ROOT_COMPARTMENT_ID")
}

// AvailabilityDomain return the aviailability domain defined in .env.sample file
func AvailabilityDomain() *string {
	return common.String(availabilityDomain)
}

// CompartmentID return the compartment ID defined in .env.sample file
func CompartmentID() *string {
	return common.String(compartmentID)
}

// RootCompartmentID return the root compartment ID defined in .env.sample file
func RootCompartmentID() *string {
	return common.String(rootCompartmentID)
}
