package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder provides operations to call the getAllRetainedMessages method.
type ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetQueryParameters get all retained messages across all channels in a team. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
type ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetQueryParameters
}
// NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderInternal instantiates a new ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder and sets the default values.
func NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) {
    m := &ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/joinedTeams/{team%2Did}/channels/getAllRetainedMessages(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder instantiates a new ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder and sets the default values.
func NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get all retained messages across all channels in a team. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// Deprecated: This method is obsolete. Use GetAsGetAllRetainedMessagesGetResponse instead.
// returns a ItemJoinedTeamsItemChannelsGetAllRetainedMessagesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/channel-getallretainedmessages?view=graph-rest-1.0
func (m *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(ItemJoinedTeamsItemChannelsGetAllRetainedMessagesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemJoinedTeamsItemChannelsGetAllRetainedMessagesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemJoinedTeamsItemChannelsGetAllRetainedMessagesResponseable), nil
}
// GetAsGetAllRetainedMessagesGetResponse get all retained messages across all channels in a team. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// returns a ItemJoinedTeamsItemChannelsGetAllRetainedMessagesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/channel-getallretainedmessages?view=graph-rest-1.0
func (m *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) GetAsGetAllRetainedMessagesGetResponse(ctx context.Context, requestConfiguration *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(ItemJoinedTeamsItemChannelsGetAllRetainedMessagesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemJoinedTeamsItemChannelsGetAllRetainedMessagesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemJoinedTeamsItemChannelsGetAllRetainedMessagesGetResponseable), nil
}
// ToGetRequestInformation get all retained messages across all channels in a team. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// returns a *RequestInformation when successful
func (m *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder when successful
func (m *ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) WithUrl(rawUrl string)(*ItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder) {
    return NewItemJoinedTeamsItemChannelsGetAllRetainedMessagesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
