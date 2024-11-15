package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemMailFoldersItemMessagesDeltaRequestBuilder provides operations to call the delta method.
type ItemMailFoldersItemMessagesDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemMailFoldersItemMessagesDeltaRequestBuilderGetQueryParameters get a set of messages that have been added, deleted, or updated in a specified folder. A delta function call for messages in a folder is similar to a GET request, except that by appropriatelyapplying state tokens in one or more of these calls, you can [query for incremental changes in the messages inthat folder](/graph/delta-query-messages). This allows you to maintain and synchronize a local store of a user's messages withouthaving to fetch the entire set of messages from the server every time.
type ItemMailFoldersItemMessagesDeltaRequestBuilderGetQueryParameters struct {
    // A custom query option to filter the delta response based on the type of change. Supported values are created, updated or deleted.
    ChangeType *string `uriparametername:"changeType"`
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
// ItemMailFoldersItemMessagesDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemMailFoldersItemMessagesDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemMailFoldersItemMessagesDeltaRequestBuilderGetQueryParameters
}
// NewItemMailFoldersItemMessagesDeltaRequestBuilderInternal instantiates a new ItemMailFoldersItemMessagesDeltaRequestBuilder and sets the default values.
func NewItemMailFoldersItemMessagesDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemMailFoldersItemMessagesDeltaRequestBuilder) {
    m := &ItemMailFoldersItemMessagesDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/mailFolders/{mailFolder%2Did}/messages/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top,changeType*}", pathParameters),
    }
    return m
}
// NewItemMailFoldersItemMessagesDeltaRequestBuilder instantiates a new ItemMailFoldersItemMessagesDeltaRequestBuilder and sets the default values.
func NewItemMailFoldersItemMessagesDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemMailFoldersItemMessagesDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemMailFoldersItemMessagesDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a set of messages that have been added, deleted, or updated in a specified folder. A delta function call for messages in a folder is similar to a GET request, except that by appropriatelyapplying state tokens in one or more of these calls, you can [query for incremental changes in the messages inthat folder](/graph/delta-query-messages). This allows you to maintain and synchronize a local store of a user's messages withouthaving to fetch the entire set of messages from the server every time.
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemMailFoldersItemMessagesDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/message-delta?view=graph-rest-1.0
func (m *ItemMailFoldersItemMessagesDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemMailFoldersItemMessagesDeltaRequestBuilderGetRequestConfiguration)(ItemMailFoldersItemMessagesDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemMailFoldersItemMessagesDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemMailFoldersItemMessagesDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a set of messages that have been added, deleted, or updated in a specified folder. A delta function call for messages in a folder is similar to a GET request, except that by appropriatelyapplying state tokens in one or more of these calls, you can [query for incremental changes in the messages inthat folder](/graph/delta-query-messages). This allows you to maintain and synchronize a local store of a user's messages withouthaving to fetch the entire set of messages from the server every time.
// returns a ItemMailFoldersItemMessagesDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/message-delta?view=graph-rest-1.0
func (m *ItemMailFoldersItemMessagesDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemMailFoldersItemMessagesDeltaRequestBuilderGetRequestConfiguration)(ItemMailFoldersItemMessagesDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemMailFoldersItemMessagesDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemMailFoldersItemMessagesDeltaGetResponseable), nil
}
// ToGetRequestInformation get a set of messages that have been added, deleted, or updated in a specified folder. A delta function call for messages in a folder is similar to a GET request, except that by appropriatelyapplying state tokens in one or more of these calls, you can [query for incremental changes in the messages inthat folder](/graph/delta-query-messages). This allows you to maintain and synchronize a local store of a user's messages withouthaving to fetch the entire set of messages from the server every time.
// returns a *RequestInformation when successful
func (m *ItemMailFoldersItemMessagesDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemMailFoldersItemMessagesDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemMailFoldersItemMessagesDeltaRequestBuilder when successful
func (m *ItemMailFoldersItemMessagesDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemMailFoldersItemMessagesDeltaRequestBuilder) {
    return NewItemMailFoldersItemMessagesDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
