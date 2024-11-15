package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemTeamRequestBuilder provides operations to manage the team property of the microsoft.graph.group entity.
type ItemTeamRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemTeamRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemTeamRequestBuilderGetQueryParameters the team associated with this group.
type ItemTeamRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemTeamRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemTeamRequestBuilderGetQueryParameters
}
// ItemTeamRequestBuilderPutRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamRequestBuilderPutRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllChannels provides operations to manage the allChannels property of the microsoft.graph.team entity.
// returns a *ItemTeamAllChannelsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) AllChannels()(*ItemTeamAllChannelsRequestBuilder) {
    return NewItemTeamAllChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Archive provides operations to call the archive method.
// returns a *ItemTeamArchiveRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Archive()(*ItemTeamArchiveRequestBuilder) {
    return NewItemTeamArchiveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Channels provides operations to manage the channels property of the microsoft.graph.team entity.
// returns a *ItemTeamChannelsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Channels()(*ItemTeamChannelsRequestBuilder) {
    return NewItemTeamChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Clone provides operations to call the clone method.
// returns a *ItemTeamCloneRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Clone()(*ItemTeamCloneRequestBuilder) {
    return NewItemTeamCloneRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CompleteMigration provides operations to call the completeMigration method.
// returns a *ItemTeamCompleteMigrationRequestBuilder when successful
func (m *ItemTeamRequestBuilder) CompleteMigration()(*ItemTeamCompleteMigrationRequestBuilder) {
    return NewItemTeamCompleteMigrationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemTeamRequestBuilderInternal instantiates a new ItemTeamRequestBuilder and sets the default values.
func NewItemTeamRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamRequestBuilder) {
    m := &ItemTeamRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/team{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemTeamRequestBuilder instantiates a new ItemTeamRequestBuilder and sets the default values.
func NewItemTeamRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTeamRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property team for groups
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemTeamRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemTeamRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get the team associated with this group.
// returns a Teamable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemTeamRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemTeamRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
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
// returns a *ItemTeamGroupRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Group()(*ItemTeamGroupRequestBuilder) {
    return NewItemTeamGroupRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IncomingChannels provides operations to manage the incomingChannels property of the microsoft.graph.team entity.
// returns a *ItemTeamIncomingChannelsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) IncomingChannels()(*ItemTeamIncomingChannelsRequestBuilder) {
    return NewItemTeamIncomingChannelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// InstalledApps provides operations to manage the installedApps property of the microsoft.graph.team entity.
// returns a *ItemTeamInstalledAppsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) InstalledApps()(*ItemTeamInstalledAppsRequestBuilder) {
    return NewItemTeamInstalledAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Members provides operations to manage the members property of the microsoft.graph.team entity.
// returns a *ItemTeamMembersRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Members()(*ItemTeamMembersRequestBuilder) {
    return NewItemTeamMembersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.team entity.
// returns a *ItemTeamOperationsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Operations()(*ItemTeamOperationsRequestBuilder) {
    return NewItemTeamOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PermissionGrants provides operations to manage the permissionGrants property of the microsoft.graph.team entity.
// returns a *ItemTeamPermissionGrantsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) PermissionGrants()(*ItemTeamPermissionGrantsRequestBuilder) {
    return NewItemTeamPermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Photo provides operations to manage the photo property of the microsoft.graph.team entity.
// returns a *ItemTeamPhotoRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Photo()(*ItemTeamPhotoRequestBuilder) {
    return NewItemTeamPhotoRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PrimaryChannel provides operations to manage the primaryChannel property of the microsoft.graph.team entity.
// returns a *ItemTeamPrimaryChannelRequestBuilder when successful
func (m *ItemTeamRequestBuilder) PrimaryChannel()(*ItemTeamPrimaryChannelRequestBuilder) {
    return NewItemTeamPrimaryChannelRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Put create a new team under a group. In order to create a team, the group must have a least one owner. If the creation of the team call is delayed, you can retry the call up to three times before you have to wait for 15 minutes due to a propagation delay. If the group was created less than 15 minutes ago, the call might fail with a 404 error code due to replication delays. If the group was created less than 15 minutes ago, it's possible for a call to create a team to fail with a 404 error code, due to ongoing replication delays.The recommended pattern is to retry the Create team call three times, with a 10 second delay between calls.
// returns a Teamable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/team-put-teams?view=graph-rest-1.0
func (m *ItemTeamRequestBuilder) Put(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *ItemTeamRequestBuilderPutRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, error) {
    requestInfo, err := m.ToPutRequestInformation(ctx, body, requestConfiguration);
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
// Schedule provides operations to manage the schedule property of the microsoft.graph.team entity.
// returns a *ItemTeamScheduleRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Schedule()(*ItemTeamScheduleRequestBuilder) {
    return NewItemTeamScheduleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SendActivityNotification provides operations to call the sendActivityNotification method.
// returns a *ItemTeamSendActivityNotificationRequestBuilder when successful
func (m *ItemTeamRequestBuilder) SendActivityNotification()(*ItemTeamSendActivityNotificationRequestBuilder) {
    return NewItemTeamSendActivityNotificationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Tags provides operations to manage the tags property of the microsoft.graph.team entity.
// returns a *ItemTeamTagsRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Tags()(*ItemTeamTagsRequestBuilder) {
    return NewItemTeamTagsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Template provides operations to manage the template property of the microsoft.graph.team entity.
// returns a *ItemTeamTemplateRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Template()(*ItemTeamTemplateRequestBuilder) {
    return NewItemTeamTemplateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property team for groups
// returns a *RequestInformation when successful
func (m *ItemTeamRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemTeamRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the team associated with this group.
// returns a *RequestInformation when successful
func (m *ItemTeamRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemTeamRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPutRequestInformation create a new team under a group. In order to create a team, the group must have a least one owner. If the creation of the team call is delayed, you can retry the call up to three times before you have to wait for 15 minutes due to a propagation delay. If the group was created less than 15 minutes ago, the call might fail with a 404 error code due to replication delays. If the group was created less than 15 minutes ago, it's possible for a call to create a team to fail with a 404 error code, due to ongoing replication delays.The recommended pattern is to retry the Create team call three times, with a 10 second delay between calls.
// returns a *RequestInformation when successful
func (m *ItemTeamRequestBuilder) ToPutRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Teamable, requestConfiguration *ItemTeamRequestBuilderPutRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PUT, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *ItemTeamUnarchiveRequestBuilder when successful
func (m *ItemTeamRequestBuilder) Unarchive()(*ItemTeamUnarchiveRequestBuilder) {
    return NewItemTeamUnarchiveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemTeamRequestBuilder when successful
func (m *ItemTeamRequestBuilder) WithUrl(rawUrl string)(*ItemTeamRequestBuilder) {
    return NewItemTeamRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
