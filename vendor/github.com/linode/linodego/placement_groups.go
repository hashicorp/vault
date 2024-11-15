package linodego

import "context"

// PlacementGroupType is an enum that determines the affinity policy
// for Linodes in a placement group.
type PlacementGroupType string

const (
	PlacementGroupTypeAntiAffinityLocal PlacementGroupType = "anti_affinity:local"
)

// PlacementGroupPolicy is an enum for the policy that determines whether a
// Linode can be assigned to a Placement Group.
type PlacementGroupPolicy string

const (
	PlacementGroupPolicyStrict   PlacementGroupPolicy = "strict"
	PlacementGroupPolicyFlexible PlacementGroupPolicy = "flexible"
)

// PlacementGroupMember represents a single Linode assigned to a
// placement group.
type PlacementGroupMember struct {
	LinodeID    int  `json:"linode_id"`
	IsCompliant bool `json:"is_compliant"`
}

// PlacementGroup represents a Linode placement group.
// NOTE: Placement Groups may not currently be available to all users.
type PlacementGroup struct {
	ID                   int                    `json:"id"`
	Label                string                 `json:"label"`
	Region               string                 `json:"region"`
	PlacementGroupType   PlacementGroupType     `json:"placement_group_type"`
	PlacementGroupPolicy PlacementGroupPolicy   `json:"placement_group_policy"`
	IsCompliant          bool                   `json:"is_compliant"`
	Members              []PlacementGroupMember `json:"members"`
}

// PlacementGroupCreateOptions represents the options to use
// when creating a placement group.
type PlacementGroupCreateOptions struct {
	Label                string               `json:"label"`
	Region               string               `json:"region"`
	PlacementGroupType   PlacementGroupType   `json:"placement_group_type"`
	PlacementGroupPolicy PlacementGroupPolicy `json:"placement_group_policy"`
}

// PlacementGroupUpdateOptions represents the options to use
// when updating a placement group.
type PlacementGroupUpdateOptions struct {
	Label string `json:"label,omitempty"`
}

// PlacementGroupAssignOptions represents options used when
// assigning Linodes to a placement group.
type PlacementGroupAssignOptions struct {
	Linodes       []int `json:"linodes"`
	CompliantOnly *bool `json:"compliant_only,omitempty"`
}

// PlacementGroupUnAssignOptions represents options used when
// unassigning Linodes from a placement group.
type PlacementGroupUnAssignOptions struct {
	Linodes []int `json:"linodes"`
}

// ListPlacementGroups lists placement groups under the current account
// matching the given list options.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) ListPlacementGroups(
	ctx context.Context,
	options *ListOptions,
) ([]PlacementGroup, error) {
	return getPaginatedResults[PlacementGroup](
		ctx,
		c,
		"placement/groups",
		options,
	)
}

// GetPlacementGroup gets a placement group with the specified ID.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) GetPlacementGroup(
	ctx context.Context,
	id int,
) (*PlacementGroup, error) {
	return doGETRequest[PlacementGroup](
		ctx,
		c,
		formatAPIPath("placement/groups/%d", id),
	)
}

// CreatePlacementGroup creates a placement group with the specified options.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) CreatePlacementGroup(
	ctx context.Context,
	options PlacementGroupCreateOptions,
) (*PlacementGroup, error) {
	return doPOSTRequest[PlacementGroup](
		ctx,
		c,
		"placement/groups",
		options,
	)
}

// UpdatePlacementGroup updates a placement group with the specified ID using the provided options.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) UpdatePlacementGroup(
	ctx context.Context,
	id int,
	options PlacementGroupUpdateOptions,
) (*PlacementGroup, error) {
	return doPUTRequest[PlacementGroup](
		ctx,
		c,
		formatAPIPath("placement/groups/%d", id),
		options,
	)
}

// AssignPlacementGroupLinodes assigns the specified Linodes to the given
// placement group.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) AssignPlacementGroupLinodes(
	ctx context.Context,
	id int,
	options PlacementGroupAssignOptions,
) (*PlacementGroup, error) {
	return doPOSTRequest[PlacementGroup](
		ctx,
		c,
		formatAPIPath("placement/groups/%d/assign", id),
		options,
	)
}

// UnassignPlacementGroupLinodes un-assigns the specified Linodes from the given
// placement group.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) UnassignPlacementGroupLinodes(
	ctx context.Context,
	id int,
	options PlacementGroupUnAssignOptions,
) (*PlacementGroup, error) {
	return doPOSTRequest[PlacementGroup](
		ctx,
		c,
		formatAPIPath("placement/groups/%d/unassign", id),
		options,
	)
}

// DeletePlacementGroup deletes a placement group with the specified ID.
// NOTE: Placement Groups may not currently be available to all users.
func (c *Client) DeletePlacementGroup(
	ctx context.Context,
	id int,
) error {
	return doDELETERequest(
		ctx,
		c,
		formatAPIPath("placement/groups/%d", id),
	)
}
