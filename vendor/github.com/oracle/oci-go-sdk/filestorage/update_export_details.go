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

// UpdateExportDetails Details for updating the export.
type UpdateExportDetails struct {

	// New export options for the export.
	// **Setting to the empty array will make the export invisible to all clients.**
	// Leaving unset will leave the `exportOptions` unchanged.
	ExportOptions []ClientOptions `mandatory:"false" json:"exportOptions"`
}

func (m UpdateExportDetails) String() string {
	return common.PointerString(m)
}
