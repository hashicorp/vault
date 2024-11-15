// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type GroupsClient interface {
	AddGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) error
	RemoveGroupMember(ctx context.Context, groupObjectID, memberObjectID string) error
	GetGroup(ctx context.Context, objectID string) (result Group, err error)
	ListGroups(ctx context.Context, filter string) (result []Group, err error)
}

type Group struct {
	ID          string
	DisplayName string
}

func (c *MSGraphClient) AddGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) error {
	req := models.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", memberObjectID)
	req.SetOdataId(&odataId)

	return c.client.Groups().ByGroupId(groupObjectID).Members().Ref().Post(ctx, req, nil)
}

func (c *MSGraphClient) RemoveGroupMember(ctx context.Context, groupObjectID, memberObjectID string) error {
	return c.client.Groups().ByGroupId(groupObjectID).Members().ByDirectoryObjectId(memberObjectID).Ref().Delete(ctx, nil)
}

func (c *MSGraphClient) GetGroup(ctx context.Context, groupID string) (Group, error) {
	resp, err := c.client.Groups().ByGroupId(groupID).Get(ctx, nil)
	if err != nil {
		return Group{}, err
	}

	return getGroupResponse(resp), nil
}

func (c *MSGraphClient) ListGroups(ctx context.Context, filter string) ([]Group, error) {
	req := &groups.GroupsRequestBuilderGetQueryParameters{
		Filter: &filter,
	}
	configuration := &groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: req,
	}

	groupList, err := c.client.Groups().Get(ctx, configuration)
	if err != nil {
		return nil, err
	}

	var g []Group
	for _, group := range groupList.GetValue() {
		g = append(g, getGroupResponse(group))
	}

	return g, nil
}

func getGroupResponse(group models.Groupable) Group {
	if group != nil {
		return Group{
			ID:          ptrToString(group.GetId()),
			DisplayName: ptrToString(group.GetDisplayName()),
		}
	}

	// return zero-value result if app in nil
	// or fields can't be dereferenced
	return Group{
		ID:          "",
		DisplayName: "",
	}
}
