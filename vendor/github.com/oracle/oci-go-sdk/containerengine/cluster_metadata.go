// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Container Engine for Kubernetes API
//
// API for the Container Engine for Kubernetes service. Use this API to build, deploy,
// and manage cloud-native applications. For more information, see
// Overview of Container Engine for Kubernetes (https://docs.cloud.oracle.com/iaas/Content/ContEng/Concepts/contengoverview.htm).
//

package containerengine

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ClusterMetadata The properties that define meta data for a cluster.
type ClusterMetadata struct {

	// The time the cluster was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The user who created the cluster.
	CreatedByUserId *string `mandatory:"false" json:"createdByUserId"`

	// The OCID of the work request which created the cluster.
	CreatedByWorkRequestId *string `mandatory:"false" json:"createdByWorkRequestId"`

	// The time the cluster was deleted.
	TimeDeleted *common.SDKTime `mandatory:"false" json:"timeDeleted"`

	// The user who deleted the cluster.
	DeletedByUserId *string `mandatory:"false" json:"deletedByUserId"`

	// The OCID of the work request which deleted the cluster.
	DeletedByWorkRequestId *string `mandatory:"false" json:"deletedByWorkRequestId"`

	// The time the cluster was updated.
	TimeUpdated *common.SDKTime `mandatory:"false" json:"timeUpdated"`

	// The user who updated the cluster.
	UpdatedByUserId *string `mandatory:"false" json:"updatedByUserId"`

	// The OCID of the work request which updated the cluster.
	UpdatedByWorkRequestId *string `mandatory:"false" json:"updatedByWorkRequestId"`
}

func (m ClusterMetadata) String() string {
	return common.PointerString(m)
}
