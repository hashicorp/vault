package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemInsightsUsedItemResourceRequestBuilder provides operations to manage the resource property of the microsoft.graph.usedInsight entity.
type ItemInsightsUsedItemResourceRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemInsightsUsedItemResourceRequestBuilderGetQueryParameters used for navigating to the item that was used. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
type ItemInsightsUsedItemResourceRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemInsightsUsedItemResourceRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemInsightsUsedItemResourceRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemInsightsUsedItemResourceRequestBuilderGetQueryParameters
}
// NewItemInsightsUsedItemResourceRequestBuilderInternal instantiates a new ItemInsightsUsedItemResourceRequestBuilder and sets the default values.
func NewItemInsightsUsedItemResourceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemInsightsUsedItemResourceRequestBuilder) {
    m := &ItemInsightsUsedItemResourceRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/insights/used/{usedInsight%2Did}/resource{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemInsightsUsedItemResourceRequestBuilder instantiates a new ItemInsightsUsedItemResourceRequestBuilder and sets the default values.
func NewItemInsightsUsedItemResourceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemInsightsUsedItemResourceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemInsightsUsedItemResourceRequestBuilderInternal(urlParams, requestAdapter)
}
// Get used for navigating to the item that was used. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
// returns a Entityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemInsightsUsedItemResourceRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemInsightsUsedItemResourceRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEntityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable), nil
}
// ToGetRequestInformation used for navigating to the item that was used. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
// returns a *RequestInformation when successful
func (m *ItemInsightsUsedItemResourceRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemInsightsUsedItemResourceRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemInsightsUsedItemResourceRequestBuilder when successful
func (m *ItemInsightsUsedItemResourceRequestBuilder) WithUrl(rawUrl string)(*ItemInsightsUsedItemResourceRequestBuilder) {
    return NewItemInsightsUsedItemResourceRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
