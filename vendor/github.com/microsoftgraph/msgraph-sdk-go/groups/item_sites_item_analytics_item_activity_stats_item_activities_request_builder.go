package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder provides operations to manage the activities property of the microsoft.graph.itemActivityStat entity.
type ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetQueryParameters exposes the itemActivities represented in this itemActivityStat resource.
type ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetQueryParameters struct {
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
// ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetQueryParameters
}
// ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByItemActivityId provides operations to manage the activities property of the microsoft.graph.itemActivityStat entity.
// returns a *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesItemActivityItemRequestBuilder when successful
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) ByItemActivityId(itemActivityId string)(*ItemSitesItemAnalyticsItemActivityStatsItemActivitiesItemActivityItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if itemActivityId != "" {
        urlTplParams["itemActivity%2Did"] = itemActivityId
    }
    return NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesItemActivityItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderInternal instantiates a new ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder and sets the default values.
func NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) {
    m := &ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}/analytics/itemActivityStats/{itemActivityStat%2Did}/activities{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder instantiates a new ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder and sets the default values.
func NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesCountRequestBuilder when successful
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) Count()(*ItemSitesItemAnalyticsItemActivityStatsItemActivitiesCountRequestBuilder) {
    return NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get exposes the itemActivities represented in this itemActivityStat resource.
// returns a ItemActivityCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateItemActivityCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityCollectionResponseable), nil
}
// Post create new navigation property to activities for groups
// returns a ItemActivityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityable, requestConfiguration *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateItemActivityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityable), nil
}
// ToGetRequestInformation exposes the itemActivities represented in this itemActivityStat resource.
// returns a *RequestInformation when successful
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to activities for groups
// returns a *RequestInformation when successful
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ItemActivityable, requestConfiguration *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder when successful
func (m *ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) WithUrl(rawUrl string)(*ItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder) {
    return NewItemSitesItemAnalyticsItemActivityStatsItemActivitiesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
