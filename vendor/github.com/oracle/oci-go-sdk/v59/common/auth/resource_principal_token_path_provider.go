// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	imdsPathTemplate = "/20180711/resourcePrincipalToken/{id}"
	instanceIDURL    = `http://169.254.169.254/opc/v2/instance/id`

	//ResourcePrincipalTokenPath path for retrieving the Resource Principal Token
	ResourcePrincipalTokenPath = "OCI_RESOURCE_PRINCIPAL_RPT_PATH"
	//ResourceID OCID for the resource for Resource Principal
	ResourceID = "OCI_RESOURCE_PRINCIPAL_RPT_ID"
)

// PathProvider is an interface that returns path and resource ID
type PathProvider interface {
	Path() (*string, error)
	ResourceID() (*string, error)
}

// StringRptPathProvider is a simple path provider that takes a string and returns it
type StringRptPathProvider struct {
	path       string
	resourceID string
}

// Path returns the resource principal token path
func (pp StringRptPathProvider) Path() (*string, error) {
	return &pp.path, nil
}

// ResourceID returns the resource associated with the resource principal
func (pp StringRptPathProvider) ResourceID() (*string, error) {
	return &pp.resourceID, nil
}

// ImdsRptPathProvider sets the path from a default value and the resource ID from instance metadata
type ImdsRptPathProvider struct{}

// Path returns the resource principal token path
func (pp ImdsRptPathProvider) Path() (*string, error) {
	path := imdsPathTemplate
	return &path, nil
}

// ResourceID returns the resource associated with the resource principal
func (pp ImdsRptPathProvider) ResourceID() (*string, error) {
	instanceID, err := getInstanceIDFromMetadata()
	return &instanceID, err
}

// EnvRptPathProvider sets the path and resource ID from environment variables
type EnvRptPathProvider struct{}

// Path returns the resource principal token path
func (pp EnvRptPathProvider) Path() (*string, error) {
	path := requireEnv(ResourcePrincipalTokenPath)
	if path == nil {
		return nil, fmt.Errorf("missing %s env var", ResourcePrincipalTokenPath)
	}
	return path, nil
}

// ResourceID returns the resource associated with the resource principal
func (pp EnvRptPathProvider) ResourceID() (*string, error) {
	rpID := requireEnv(ResourceID)
	if rpID == nil {
		return nil, fmt.Errorf("missing %s env var", ResourceID)
	}
	return rpID, nil
}

//DefaultRptPathProvider path provider makes sure the behavior happens with the correct fallback.
//
//For the path,
//Use the contents of the OCI_RESOURCE_PRINCIPAL_RPT_PATH environment variable, if set.
//Otherwise, use the current path: "/20180711/resourcePrincipalToken/{id}"
//
//For the resource id,
//Use the contents of the OCI_RESOURCE_PRINCIPAL_RPT_ID environment variable, if set.
//Otherwise, use IMDS to get the instance id
//
//This path provider is used when the caller doesn't provide a specific path provider to the resource principals signer
type DefaultRptPathProvider struct {
	path       string
	resourceID string
}

// Path returns the resource principal token path
func (pp DefaultRptPathProvider) Path() (*string, error) {
	path := requireEnv(ResourcePrincipalTokenPath)
	if path == nil {
		rpPath := imdsPathTemplate
		return &rpPath, nil
	}
	return path, nil
}

// ResourceID returns the resource associated with the resource principal
func (pp DefaultRptPathProvider) ResourceID() (*string, error) {
	rpID := requireEnv(ResourceID)
	if rpID == nil {
		instanceID, err := getInstanceIDFromMetadata()
		if err != nil {
			return nil, err
		}
		return &instanceID, nil
	}
	return rpID, nil
}

func getInstanceIDFromMetadata() (instanceID string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", instanceIDURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer Oracle")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	return bodyString, nil
}
