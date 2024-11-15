package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemTodoListsItemTasksDeltaRequestBuilder provides operations to call the delta method.
type ItemTodoListsItemTasksDeltaRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemTodoListsItemTasksDeltaRequestBuilderGetQueryParameters get a set of todoTask resources that have been added, deleted, or updated in a specified todoTaskList. A delta function call for todoTask resources in a todoTaskList is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the todoTask in that todoTaskList. This allows you to maintain and synchronize a local store of a user's todoTask resources without having to fetch the entire set from the server every time.  
type ItemTodoListsItemTasksDeltaRequestBuilderGetQueryParameters struct {
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
// ItemTodoListsItemTasksDeltaRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTodoListsItemTasksDeltaRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemTodoListsItemTasksDeltaRequestBuilderGetQueryParameters
}
// NewItemTodoListsItemTasksDeltaRequestBuilderInternal instantiates a new ItemTodoListsItemTasksDeltaRequestBuilder and sets the default values.
func NewItemTodoListsItemTasksDeltaRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTodoListsItemTasksDeltaRequestBuilder) {
    m := &ItemTodoListsItemTasksDeltaRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/todo/lists/{todoTaskList%2Did}/tasks/delta(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemTodoListsItemTasksDeltaRequestBuilder instantiates a new ItemTodoListsItemTasksDeltaRequestBuilder and sets the default values.
func NewItemTodoListsItemTasksDeltaRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTodoListsItemTasksDeltaRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTodoListsItemTasksDeltaRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get a set of todoTask resources that have been added, deleted, or updated in a specified todoTaskList. A delta function call for todoTask resources in a todoTaskList is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the todoTask in that todoTaskList. This allows you to maintain and synchronize a local store of a user's todoTask resources without having to fetch the entire set from the server every time.  
// Deprecated: This method is obsolete. Use GetAsDeltaGetResponse instead.
// returns a ItemTodoListsItemTasksDeltaResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/todotask-delta?view=graph-rest-1.0
func (m *ItemTodoListsItemTasksDeltaRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemTodoListsItemTasksDeltaRequestBuilderGetRequestConfiguration)(ItemTodoListsItemTasksDeltaResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemTodoListsItemTasksDeltaResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemTodoListsItemTasksDeltaResponseable), nil
}
// GetAsDeltaGetResponse get a set of todoTask resources that have been added, deleted, or updated in a specified todoTaskList. A delta function call for todoTask resources in a todoTaskList is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the todoTask in that todoTaskList. This allows you to maintain and synchronize a local store of a user's todoTask resources without having to fetch the entire set from the server every time.  
// returns a ItemTodoListsItemTasksDeltaGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/todotask-delta?view=graph-rest-1.0
func (m *ItemTodoListsItemTasksDeltaRequestBuilder) GetAsDeltaGetResponse(ctx context.Context, requestConfiguration *ItemTodoListsItemTasksDeltaRequestBuilderGetRequestConfiguration)(ItemTodoListsItemTasksDeltaGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemTodoListsItemTasksDeltaGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemTodoListsItemTasksDeltaGetResponseable), nil
}
// ToGetRequestInformation get a set of todoTask resources that have been added, deleted, or updated in a specified todoTaskList. A delta function call for todoTask resources in a todoTaskList is similar to a GET request, except that by appropriately applying state tokens in one or more of these calls, you can query for incremental changes in the todoTask in that todoTaskList. This allows you to maintain and synchronize a local store of a user's todoTask resources without having to fetch the entire set from the server every time.  
// returns a *RequestInformation when successful
func (m *ItemTodoListsItemTasksDeltaRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemTodoListsItemTasksDeltaRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemTodoListsItemTasksDeltaRequestBuilder when successful
func (m *ItemTodoListsItemTasksDeltaRequestBuilder) WithUrl(rawUrl string)(*ItemTodoListsItemTasksDeltaRequestBuilder) {
    return NewItemTodoListsItemTasksDeltaRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
