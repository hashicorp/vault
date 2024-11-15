// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package resource provides various routines related to a resource.
package resource

import (
	"fmt"

	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
)

// Validate will validate that the Link configuration is valid by making sure that Location, Type and ID are provided.
func Validate(resource cloud.HashicorpCloudLocationLink) error {
	if resource.Location == nil {
		return fmt.Errorf("missing resource location")
	}
	if resource.Type == "" {
		return fmt.Errorf("missing resource type")
	}
	if resource.ID == "" {
		return fmt.Errorf("missing resource ID")
	}
	return nil
}
