package api

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var _ GroupsClient = (*ActiveDirectoryApplicationGroupsClient)(nil)

type aadGroupsClient interface {
	AddMember(ctx context.Context, groupObjectID string, parameters graphrbac.GroupAddMemberParameters) (result autorest.Response, err error)
	RemoveMember(ctx context.Context, groupObjectID string, memberObjectID string) (result autorest.Response, err error)
	Get(ctx context.Context, objectID string) (result graphrbac.ADGroup, err error)
	List(ctx context.Context, filter string) (result graphrbac.GroupListResultPage, err error)
}

type ActiveDirectoryApplicationGroupsClient struct {
	BaseURI  string
	TenantID string
	Client   aadGroupsClient
}

func (a ActiveDirectoryApplicationGroupsClient) AddGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) error {
	uri := fmt.Sprintf("%s/%s/directoryObjects/%s", a.BaseURI, a.TenantID, memberObjectID)
	aadParams := graphrbac.GroupAddMemberParameters{
		URL: to.StringPtr(uri),
	}
	_, err := a.Client.AddMember(ctx, groupObjectID, aadParams)
	return err
}

func (a ActiveDirectoryApplicationGroupsClient) RemoveGroupMember(ctx context.Context, groupObjectID string, memberObjectID string) error {
	_, err := a.Client.RemoveMember(ctx, groupObjectID, memberObjectID)
	return err
}

func (a ActiveDirectoryApplicationGroupsClient) GetGroup(ctx context.Context, objectID string) (Group, error) {
	resp, err := a.Client.Get(ctx, objectID)
	if err != nil {
		return Group{}, err
	}

	grp := getGroupFromRBAC(resp)

	return grp, nil
}

func getGroupFromRBAC(resp graphrbac.ADGroup) Group {
	grp := Group{
		ID:          *resp.ObjectID,
		DisplayName: *resp.DisplayName,
	}
	return grp
}

func (a ActiveDirectoryApplicationGroupsClient) ListGroups(ctx context.Context, filter string) ([]Group, error) {
	resp, err := a.Client.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	grps := []Group{}
	for _, aadGrp := range resp.Values() {
		grp := getGroupFromRBAC(aadGrp)
		grps = append(grps, grp)
	}
	return grps, nil
}
