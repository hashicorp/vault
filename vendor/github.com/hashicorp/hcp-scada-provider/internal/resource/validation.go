// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"

	cloud "github.com/hashicorp/hcp-sdk-go/clients/cloud-shared/v1/models"
)

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
