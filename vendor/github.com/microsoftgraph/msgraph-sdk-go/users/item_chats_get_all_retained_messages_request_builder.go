package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemChatsGetAllRetainedMessagesRequestBuilder provides operations to call the getAllRetainedMessages method.
type ItemChatsGetAllRetainedMessagesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemChatsGetAllRetainedMessagesRequestBuilderGetQueryParameters get all retained messages from all chats that a user is a participant in, including one-on-one chats, group chats, and meeting chats. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
type ItemChatsGetAllRetainedMessagesRequestBuilderGetQueryParameters struct {
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
// ItemChatsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemChatsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemChatsGetAllRetainedMessagesRequestBuilderGetQueryParameters
}
// NewItemChatsGetAllRetainedMessagesRequestBuilderInternal instantiates a new ItemChatsGetAllRetainedMessagesRequestBuilder and sets the default values.
func NewItemChatsGetAllRetainedMessagesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemChatsGetAllRetainedMessagesRequestBuilder) {
    m := &ItemChatsGetAllRetainedMessagesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/chats/getAllRetainedMessages(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemChatsGetAllRetainedMessagesRequestBuilder instantiates a new ItemChatsGetAllRetainedMessagesRequestBuilder and sets the default values.
func NewItemChatsGetAllRetainedMessagesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemChatsGetAllRetainedMessagesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemChatsGetAllRetainedMessagesRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get all retained messages from all chats that a user is a participant in, including one-on-one chats, group chats, and meeting chats. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// Deprecated: This method is obsolete. Use GetAsGetAllRetainedMessagesGetResponse instead.
// returns a ItemChatsGetAllRetainedMessagesResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/chat-getallretainedmessages?view=graph-rest-1.0
func (m *ItemChatsGetAllRetainedMessagesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemChatsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(ItemChatsGetAllRetainedMessagesResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemChatsGetAllRetainedMessagesResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemChatsGetAllRetainedMessagesResponseable), nil
}
// GetAsGetAllRetainedMessagesGetResponse get all retained messages from all chats that a user is a participant in, including one-on-one chats, group chats, and meeting chats. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// returns a ItemChatsGetAllRetainedMessagesGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/chat-getallretainedmessages?view=graph-rest-1.0
func (m *ItemChatsGetAllRetainedMessagesRequestBuilder) GetAsGetAllRetainedMessagesGetResponse(ctx context.Context, requestConfiguration *ItemChatsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(ItemChatsGetAllRetainedMessagesGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemChatsGetAllRetainedMessagesGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemChatsGetAllRetainedMessagesGetResponseable), nil
}
// ToGetRequestInformation get all retained messages from all chats that a user is a participant in, including one-on-one chats, group chats, and meeting chats. To learn more about how to use the Microsoft Teams export APIs to export content, see Export content with the Microsoft Teams export APIs.
// returns a *RequestInformation when successful
func (m *ItemChatsGetAllRetainedMessagesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemChatsGetAllRetainedMessagesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemChatsGetAllRetainedMessagesRequestBuilder when successful
func (m *ItemChatsGetAllRetainedMessagesRequestBuilder) WithUrl(rawUrl string)(*ItemChatsGetAllRetainedMessagesRequestBuilder) {
    return NewItemChatsGetAllRetainedMessagesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
