package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder provides operations to call the filterOperators method.
type ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetQueryParameters list all operators supported in the scoping filters.
type ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetQueryParameters struct {
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
// ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetQueryParameters
}
// NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderInternal instantiates a new ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder and sets the default values.
func NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) {
    m := &ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/synchronization/jobs/{synchronizationJob%2Did}/schema/filterOperators(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder instantiates a new ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder and sets the default values.
func NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list all operators supported in the scoping filters.
// Deprecated: This method is obsolete. Use GetAsFilterOperatorsGetResponse instead.
// returns a ItemSynchronizationJobsItemSchemaFilterOperatorsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-filteroperators?view=graph-rest-1.0
func (m *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(ItemSynchronizationJobsItemSchemaFilterOperatorsResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationJobsItemSchemaFilterOperatorsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationJobsItemSchemaFilterOperatorsResponseable), nil
}
// GetAsFilterOperatorsGetResponse list all operators supported in the scoping filters.
// returns a ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-filteroperators?view=graph-rest-1.0
func (m *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) GetAsFilterOperatorsGetResponse(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable), nil
}
// ToGetRequestInformation list all operators supported in the scoping filters.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder when successful
func (m *ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder) {
    return NewItemSynchronizationJobsItemSchemaFilterOperatorsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
