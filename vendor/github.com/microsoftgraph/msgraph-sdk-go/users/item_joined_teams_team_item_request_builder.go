package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemJoinedTeamsTeamItemRequestBuilder provides operations to manage the joinedTeams property of the microsoft.graph.user entity.
type ItemJoinedTeamsTeamItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemJoinedTeamsTeamItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemJoinedTeamsTeamItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemJoinedTeamsTeamItemRequestBuilderGetQueryParameters get joinedTeams from users
type ItemJoinedTeamsTeamItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemJoinedTeamsTeamItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemJoinedTeamsTeamItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemJoinedTeamsTeamItemRequestBuilderGetQueryParameters
}
// ItemJoinedTeamsTeamItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemJoinedTeamsTeamItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllChannels provides operations to manage the allChannels property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemAllChannelsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) AllChannels()(*ItemJoinedTeamsItemAllChannelsRequestBuilder) {
    return NewItemJoinedTeamsItemAllChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Archive provides operations to call the archive method.
// returns a *ItemJoinedTeamsItemArchiveRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Archive()(*ItemJoinedTeamsItemArchiveRequestBuilder) {
    return NewItemJoinedTeamsItemArchiveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemChannelsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Channels()(*ItemJoinedTeamsItemChannelsRequestBuilder) {
    return NewItemJoinedTeamsItemChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Clone provides operations to call the clone method.
// returns a *ItemJoinedTeamsItemCloneRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Clone()(*ItemJoinedTeamsItemCloneRequestBuilder) {
    return NewItemJoinedTeamsItemCloneRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CompleteMigration provides operations to call the completeMigration method.
// returns a *ItemJoinedTeamsItemCompleteMigrationRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) CompleteMigration()(*ItemJoinedTeamsItemCompleteMigrationRequestBuilder) {
    return NewItemJoinedTeamsItemCompleteMigrationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemJoinedTeamsTeamItemRequestBuilderInternal instantiates a new ItemJoinedTeamsTeamItemRequestBuilder and sets the default values.
func NewItemJoinedTeamsTeamItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemJoinedTeamsTeamItemRequestBuilder) {
    m := &ItemJoinedTeamsTeamItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemJoinedTeamsTeamItemRequestBuilder instantiates a new ItemJoinedTeamsTeamItemRequestBuilder and sets the default values.
func NewItemJoinedTeamsTeamItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemJoinedTeamsTeamItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemJoinedTeamsTeamItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property joinedTeams for users
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get get joinedTeams from users
// returns a Teamable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// Group provides operations to manage the group property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemGroupRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Group()(*ItemJoinedTeamsItemGroupRequestBuilder) {
    return NewItemJoinedTeamsItemGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemIncomingChannelsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) IncomingChannels()(*ItemJoinedTeamsItemIncomingChannelsRequestBuilder) {
    return NewItemJoinedTeamsItemIncomingChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemInstalledAppsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) InstalledApps()(*ItemJoinedTeamsItemInstalledAppsRequestBuilder) {
    return NewItemJoinedTeamsItemInstalledAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemMembersRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Members()(*ItemJoinedTeamsItemMembersRequestBuilder) {
    return NewItemJoinedTeamsItemMembersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemOperationsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Operations()(*ItemJoinedTeamsItemOperationsRequestBuilder) {
    return NewItemJoinedTeamsItemOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property joinedTeams in users
// returns a Teamable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTeamFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable), nil
}
// PermissionGrants provides operations to manage the permissionGrants property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemPermissionGrantsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) PermissionGrants()(*ItemJoinedTeamsItemPermissionGrantsRequestBuilder) {
    return NewItemJoinedTeamsItemPermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Photo provides operations to manage the photo property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemPhotoRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Photo()(*ItemJoinedTeamsItemPhotoRequestBuilder) {
    return NewItemJoinedTeamsItemPhotoRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemPrimaryChannelRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) PrimaryChannel()(*ItemJoinedTeamsItemPrimaryChannelRequestBuilder) {
    return NewItemJoinedTeamsItemPrimaryChannelRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemScheduleRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Schedule()(*ItemJoinedTeamsItemScheduleRequestBuilder) {
    return NewItemJoinedTeamsItemScheduleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
// returns a *ItemJoinedTeamsItemSendActivityNotificationRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) SendActivityNotification()(*ItemJoinedTeamsItemSendActivityNotificationRequestBuilder) {
    return NewItemJoinedTeamsItemSendActivityNotificationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemTagsRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Tags()(*ItemJoinedTeamsItemTagsRequestBuilder) {
    return NewItemJoinedTeamsItemTagsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
// returns a *ItemJoinedTeamsItemTemplateRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Template()(*ItemJoinedTeamsItemTemplateRequestBuilder) {
    return NewItemJoinedTeamsItemTemplateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property joinedTeams for users
// returns a *RequestInformation when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get joinedTeams from users
// returns a *RequestInformation when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPatchRequestInformation update the navigation property joinedTeams in users
// returns a *RequestInformation when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *ItemJoinedTeamsTeamItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// Unarchive provides operations to call the unarchive method.
// returns a *ItemJoinedTeamsItemUnarchiveRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) Unarchive()(*ItemJoinedTeamsItemUnarchiveRequestBuilder) {
    return NewItemJoinedTeamsItemUnarchiveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemJoinedTeamsTeamItemRequestBuilder when successful
func (m *ItemJoinedTeamsTeamItemRequestBuilder) WithUrl(rawUrl string)(*ItemJoinedTeamsTeamItemRequestBuilder) {
    return NewItemJoinedTeamsTeamItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
