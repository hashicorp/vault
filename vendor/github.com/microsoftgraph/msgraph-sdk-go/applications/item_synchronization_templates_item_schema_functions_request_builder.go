package applications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder provides operations to call the functions method.
type ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetQueryParameters list all the functions currently supported in the attributeMappingSource.
type ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetQueryParameters struct {
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
// ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetQueryParameters
}
// NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderInternal instantiates a new ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) {
    m := &ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/applications/{application%2Did}/synchronization/templates/{synchronizationTemplate%2Did}/schema/functions(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder instantiates a new ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list all the functions currently supported in the attributeMappingSource.
// Deprecated: This method is obsolete. Use GetAsFunctionsGetResponse instead.
// returns a ItemSynchronizationTemplatesItemSchemaFunctionsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-functions?view=graph-rest-1.0
func (m *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(ItemSynchronizationTemplatesItemSchemaFunctionsResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationTemplatesItemSchemaFunctionsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationTemplatesItemSchemaFunctionsResponseable), nil
}
// GetAsFunctionsGetResponse list all the functions currently supported in the attributeMappingSource.
// returns a ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-functions?view=graph-rest-1.0
func (m *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) GetAsFunctionsGetResponse(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationTemplatesItemSchemaFunctionsGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable), nil
}
// ToGetRequestInformation list all the functions currently supported in the attributeMappingSource.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder when successful
func (m *ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder) {
    return NewItemSynchronizationTemplatesItemSchemaFunctionsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
