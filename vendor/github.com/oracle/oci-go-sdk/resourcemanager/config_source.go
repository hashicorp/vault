// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// ConfigSource Location of the zip file that contains the Terraform configuration.
type ConfigSource interface {

	// File path to the directory from which Terraform runs.
	// If not specified, we use the root directory.
	GetWorkingDirectory() *string
}

type configsource struct {
	JsonData         []byte
	WorkingDirectory *string `mandatory:"false" json:"workingDirectory"`
	ConfigSourceType string  `json:"configSourceType"`
}

// UnmarshalJSON unmarshals json
func (m *configsource) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerconfigsource configsource
	s := struct {
		Model Unmarshalerconfigsource
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.WorkingDirectory = s.Model.WorkingDirectory
	m.ConfigSourceType = s.Model.ConfigSourceType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *configsource) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.ConfigSourceType {
	case "ZIP_UPLOAD":
		mm := ZipUploadConfigSource{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetWorkingDirectory returns WorkingDirectory
func (m configsource) GetWorkingDirectory() *string {
	return m.WorkingDirectory
}

func (m configsource) String() string {
	return common.PointerString(m)
}
