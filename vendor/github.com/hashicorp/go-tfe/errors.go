package tfe

import (
	"errors"
)

// Generic errors applicable to all resources.
var (
	// ErrUnauthorized is returned when a receiving a 401.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrResourceNotFound is returned when a receiving a 404.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrRequiredName is returned when a name option is not present.
	ErrRequiredName = errors.New("name is required")

	// ErrInvalidName is returned when the name option has invalid value.
	ErrInvalidName = errors.New("invalid value for name")

	// ErrMissingDirectory is returned when the path does not have an existing directory.
	ErrMissingDirectory = errors.New("path needs to be an existing directory")
)

// Resource Errors
var (
	// ErrWorkspaceLocked is returned when trying to lock a
	// locked workspace.
	ErrWorkspaceLocked = errors.New("workspace already locked")

	// ErrWorkspaceNotLocked is returned when trying to unlock
	// a unlocked workspace.
	ErrWorkspaceNotLocked = errors.New("workspace already unlocked")

	// ErrInvalidWorkspaceID is returned when the workspace ID is invalid.
	ErrInvalidWorkspaceID = errors.New("invalid value for workspace ID")

	// ErrInvalidWorkspaceValue is returned when workspace value is invalid.
	ErrInvalidWorkspaceValue = errors.New("invalid value for workspace")

	// ErrWorkspacesRequired is returned when the Workspaces are not present.
	ErrWorkspacesRequired = errors.New("workspaces is required")

	// ErrWorkspaceMinLimit is returned when the length of Workspaces is 0.
	ErrWorkspaceMinLimit = errors.New("must provide at least one workspace")

	// ErrMissingTagIdentifier is returned when tag resource identifiers are invalid
	ErrMissingTagIdentifier = errors.New("must specify at least one tag by ID or name")

	// Run/Apply errors

	// ErrInvalidRunID is returned when the run ID is invalid.
	ErrInvalidRunID = errors.New("invalid value for run ID")

	// ErrInvalidApplyID is returned when the apply ID is invalid.
	ErrInvalidApplyID = errors.New("invalid value for apply ID")

	// Organzation errors

	// ErrInvalidOrg is returned when the organization option has an invalid value.
	ErrInvalidOrg = errors.New("invalid value for organization")

	// Agent errors

	// ErrInvalidAgentPoolID is returned when the agent pool ID is invalid.
	ErrInvalidAgentPoolID = errors.New("invalid value for agent pool ID")

	// ErrInvalidAgentTokenID is returned when the agent toek ID is invalid.
	ErrInvalidAgentTokenID = errors.New("invalid value for agent token ID")

	// Token errors

	// ErrAgentTokenDescription is returned when the description is blank.
	ErrAgentTokenDescription = errors.New("agent token description can't be blank")

	// Config errors

	// ErrInvalidConfigVersionID is returned when the configuration version ID is invalid.
	ErrInvalidConfigVersionID = errors.New("invalid value for configuration version ID")

	// Cost Esimation Errors

	// ErrInvalidCostEstimateID is returned when the cost estimate ID is invalid.
	ErrInvalidCostEstimateID = errors.New("invalid value for cost estimate ID")

	// User

	// ErrInvalidUservalue is invalid.
	ErrInvalidUserValue = errors.New("invalid value for user")

	// Settings

	// ErrInvalidSMTPAuth is returned when the smtp auth type is not valid.
	ErrInvalidSMTPAuth = errors.New("invalid smtp auth type")

	// Terraform Versions

	// ErrInvalidTerraformVersionID is returned when the ID for a terraform
	// version is invalid.
	ErrInvalidTerraformVersionID = errors.New("invalid value for terraform version ID")

	// ErrInvalidTerraformVersionType is returned when the type is not valid.
	ErrInvalidTerraformVersionType = errors.New("invalid type for terraform version. Please use 'terraform-version'")
)
