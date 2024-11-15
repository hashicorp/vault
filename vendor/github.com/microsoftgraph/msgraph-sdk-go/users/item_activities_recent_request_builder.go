package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemActivitiesRecentRequestBuilder provides operations to call the recent method.
type ItemActivitiesRecentRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemActivitiesRecentRequestBuilderGetQueryParameters get recent activities for a given user. This OData function has some default behaviors included to make it operate like a 'most recently used' API. The service queries for the most recent activityHistoryItems, and then pull those related activities. Activities are sorted according to the most recent lastModified on the activityHistoryItem. This means that activities without activityHistoryItems won't be included in the response. The UserActivity.ReadWrite.CreatedByApp permission will also apply extra filtering to the response, so that only activities created by your application are returned. This server-side filtering might result in empty pages if the user is active and other applications have created more recent activities. To get your application's activities, use the nextLink property to paginate.
type ItemActivitiesRecentRequestBuilderGetQueryParameters struct {
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
// ItemActivitiesRecentRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemActivitiesRecentRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemActivitiesRecentRequestBuilderGetQueryParameters
}
// NewItemActivitiesRecentRequestBuilderInternal instantiates a new ItemActivitiesRecentRequestBuilder and sets the default values.
func NewItemActivitiesRecentRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemActivitiesRecentRequestBuilder) {
    m := &ItemActivitiesRecentRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/activities/recent(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemActivitiesRecentRequestBuilder instantiates a new ItemActivitiesRecentRequestBuilder and sets the default values.
func NewItemActivitiesRecentRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemActivitiesRecentRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemActivitiesRecentRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get recent activities for a given user. This OData function has some default behaviors included to make it operate like a 'most recently used' API. The service queries for the most recent activityHistoryItems, and then pull those related activities. Activities are sorted according to the most recent lastModified on the activityHistoryItem. This means that activities without activityHistoryItems won't be included in the response. The UserActivity.ReadWrite.CreatedByApp permission will also apply extra filtering to the response, so that only activities created by your application are returned. This server-side filtering might result in empty pages if the user is active and other applications have created more recent activities. To get your application's activities, use the nextLink property to paginate.
// Deprecated: This method is obsolete. Use GetAsRecentGetResponse instead.
// returns a ItemActivitiesRecentResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/projectrome-get-recent-activities?view=graph-rest-1.0
func (m *ItemActivitiesRecentRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemActivitiesRecentRequestBuilderGetRequestConfiguration)(ItemActivitiesRecentResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemActivitiesRecentResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemActivitiesRecentResponseable), nil
}
// GetAsRecentGetResponse get recent activities for a given user. This OData function has some default behaviors included to make it operate like a 'most recently used' API. The service queries for the most recent activityHistoryItems, and then pull those related activities. Activities are sorted according to the most recent lastModified on the activityHistoryItem. This means that activities without activityHistoryItems won't be included in the response. The UserActivity.ReadWrite.CreatedByApp permission will also apply extra filtering to the response, so that only activities created by your application are returned. This server-side filtering might result in empty pages if the user is active and other applications have created more recent activities. To get your application's activities, use the nextLink property to paginate.
// returns a ItemActivitiesRecentGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/projectrome-get-recent-activities?view=graph-rest-1.0
func (m *ItemActivitiesRecentRequestBuilder) GetAsRecentGetResponse(ctx context.Context, requestConfiguration *ItemActivitiesRecentRequestBuilderGetRequestConfiguration)(ItemActivitiesRecentGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemActivitiesRecentGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemActivitiesRecentGetResponseable), nil
}
// ToGetRequestInformation get recent activities for a given user. This OData function has some default behaviors included to make it operate like a 'most recently used' API. The service queries for the most recent activityHistoryItems, and then pull those related activities. Activities are sorted according to the most recent lastModified on the activityHistoryItem. This means that activities without activityHistoryItems won't be included in the response. The UserActivity.ReadWrite.CreatedByApp permission will also apply extra filtering to the response, so that only activities created by your application are returned. This server-side filtering might result in empty pages if the user is active and other applications have created more recent activities. To get your application's activities, use the nextLink property to paginate.
// returns a *RequestInformation when successful
func (m *ItemActivitiesRecentRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemActivitiesRecentRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemActivitiesRecentRequestBuilder when successful
func (m *ItemActivitiesRecentRequestBuilder) WithUrl(rawUrl string)(*ItemActivitiesRecentRequestBuilder) {
    return NewItemActivitiesRecentRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
