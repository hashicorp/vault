// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

func isOrgKey(orgID, projectID string) bool {
	return len(orgID) > 0 && len(projectID) == 0
}

func isProjectKey(orgID, projectID string) bool {
	return len(orgID) == 0 && len(projectID) > 0
}

func isAssignedToProject(orgID, projectID string) bool {
	return len(orgID) > 0 && len(projectID) > 0
}
