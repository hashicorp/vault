package teams

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder provides operations to manage the sharedWithTeams property of the microsoft.graph.channel entity.
type ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetQueryParameters a collection of teams with which a channel is shared.
type ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetQueryParameters
}
// ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllowedMembers provides operations to manage the allowedMembers property of the microsoft.graph.sharedWithChannelTeamInfo entity.
// returns a *ItemPrimaryChannelSharedWithTeamsItemAllowedMembersRequestBuilder when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) AllowedMembers()(*ItemPrimaryChannelSharedWithTeamsItemAllowedMembersRequestBuilder) {
    return NewItemPrimaryChannelSharedWithTeamsItemAllowedMembersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderInternal instantiates a new ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder and sets the default values.
func NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) {
    m := &ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/teams/{team%2Did}/primaryChannel/sharedWithTeams/{sharedWithChannelTeamInfo%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder instantiates a new ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder and sets the default values.
func NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property sharedWithTeams for teams
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get a collection of teams with which a channel is shared.
// returns a SharedWithChannelTeamInfoable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSharedWithChannelTeamInfoFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable), nil
}
// Patch update the navigation property sharedWithTeams in teams
// returns a SharedWithChannelTeamInfoable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSharedWithChannelTeamInfoFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable), nil
}
// Team provides operations to manage the team property of the microsoft.graph.teamInfo entity.
// returns a *ItemPrimaryChannelSharedWithTeamsItemTeamRequestBuilder when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) Team()(*ItemPrimaryChannelSharedWithTeamsItemTeamRequestBuilder) {
    return NewItemPrimaryChannelSharedWithTeamsItemTeamRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property sharedWithTeams for teams
// returns a *RequestInformation when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation a collection of teams with which a channel is shared.
// returns a *RequestInformation when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property sharedWithTeams in teams
// returns a *RequestInformation when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SharedWithChannelTeamInfoable, requestConfiguration *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder when successful
func (m *ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) WithUrl(rawUrl string)(*ItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder) {
    return NewItemPrimaryChannelSharedWithTeamsSharedWithChannelTeamInfoItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
