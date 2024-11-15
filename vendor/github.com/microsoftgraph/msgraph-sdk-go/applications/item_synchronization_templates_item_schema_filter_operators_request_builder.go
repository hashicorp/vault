package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder provides operations to call the filterOperators method.
type ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetQueryParameters list all operators supported in the scoping filters.
type ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetQueryParameters struct {
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
// ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetQueryParameters
}
// NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderInternal instantiates a new ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) {
    m := &ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/synchronization/templates/{synchronizationTemplate%2Did}/schema/filterOperators(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder instantiates a new ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list all operators supported in the scoping filters.
// Deprecated: This method is obsolete. Use GetAsFilterOperatorsGetResponse instead.
// returns a ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-filteroperators?view=graph-rest-1.0
func (m *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseable), nil
}
// GetAsFilterOperatorsGetResponse list all operators supported in the scoping filters.
// returns a ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-filteroperators?view=graph-rest-1.0
func (m *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) GetAsFilterOperatorsGetResponse(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable), nil
}
// ToGetRequestInformation list all operators supported in the scoping filters.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder when successful
func (m *ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder) {
    return NewItemSynchronizationTemplatesItemSchemaFilterOperatorsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
