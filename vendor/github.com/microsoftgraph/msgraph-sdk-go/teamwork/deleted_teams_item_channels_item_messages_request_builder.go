package teamwork

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeletedTeamsItemChannelsItemMessagesRequestBuilder provides operations to manage the messages property of the microsoft.graph.channel entity.
type DeletedTeamsItemChannelsItemMessagesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeletedTeamsItemChannelsItemMessagesRequestBuilderGetQueryParameters a collection of all the messages in the channel. A navigation property. Nullable.
type DeletedTeamsItemChannelsItemMessagesRequestBuilderGetQueryParameters struct {
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
// DeletedTeamsItemChannelsItemMessagesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedTeamsItemChannelsItemMessagesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeletedTeamsItemChannelsItemMessagesRequestBuilderGetQueryParameters
}
// DeletedTeamsItemChannelsItemMessagesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeletedTeamsItemChannelsItemMessagesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByChatMessageId provides operations to manage the messages property of the microsoft.graph.channel entity.
// returns a *DeletedTeamsItemChannelsItemMessagesChatMessageItemRequestBuilder when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) ByChatMessageId(chatMessageId string)(*DeletedTeamsItemChannelsItemMessagesChatMessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if chatMessageId != "" {
        urlTplParams["chatMessage%2Did"] = chatMessageId
    }
    return NewDeletedTeamsItemChannelsItemMessagesChatMessageItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeletedTeamsItemChannelsItemMessagesRequestBuilderInternal instantiates a new DeletedTeamsItemChannelsItemMessagesRequestBuilder and sets the default values.
func NewDeletedTeamsItemChannelsItemMessagesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedTeamsItemChannelsItemMessagesRequestBuilder) {
    m := &DeletedTeamsItemChannelsItemMessagesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/teamwork/deletedTeams/{deletedTeam%2Did}/channels/{channel%2Did}/messages{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewDeletedTeamsItemChannelsItemMessagesRequestBuilder instantiates a new DeletedTeamsItemChannelsItemMessagesRequestBuilder and sets the default values.
func NewDeletedTeamsItemChannelsItemMessagesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeletedTeamsItemChannelsItemMessagesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeletedTeamsItemChannelsItemMessagesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *DeletedTeamsItemChannelsItemMessagesCountRequestBuilder when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) Count()(*DeletedTeamsItemChannelsItemMessagesCountRequestBuilder) {
    return NewDeletedTeamsItemChannelsItemMessagesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delta provides operations to call the delta method.
// returns a *DeletedTeamsItemChannelsItemMessagesDeltaRequestBuilder when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) Delta()(*DeletedTeamsItemChannelsItemMessagesDeltaRequestBuilder) {
    return NewDeletedTeamsItemChannelsItemMessagesDeltaRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get a collection of all the messages in the channel. A navigation property. Nullable.
// returns a ChatMessageCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) Get(ctx context.Context, requestConfiguration *DeletedTeamsItemChannelsItemMessagesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChatMessageCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageCollectionResponseable), nil
}
// Post create new navigation property to messages for teamwork
// returns a ChatMessageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageable, requestConfiguration *DeletedTeamsItemChannelsItemMessagesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChatMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageable), nil
}
// ToGetRequestInformation a collection of all the messages in the channel. A navigation property. Nullable.
// returns a *RequestInformation when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeletedTeamsItemChannelsItemMessagesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to messages for teamwork
// returns a *RequestInformation when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ChatMessageable, requestConfiguration *DeletedTeamsItemChannelsItemMessagesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *DeletedTeamsItemChannelsItemMessagesRequestBuilder when successful
func (m *DeletedTeamsItemChannelsItemMessagesRequestBuilder) WithUrl(rawUrl string)(*DeletedTeamsItemChannelsItemMessagesRequestBuilder) {
    return NewDeletedTeamsItemChannelsItemMessagesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
