package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder provides operations to call the functions method.
type ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetQueryParameters list all the functions currently supported in the attributeMappingSource.
type ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetQueryParameters struct {
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
// ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetQueryParameters
}
// NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilderInternal instantiates a new ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder and sets the default values.
func NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) {
    m := &ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/synchronization/jobs/{synchronizationJob%2Did}/schema/functions(){?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilder instantiates a new ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder and sets the default values.
func NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list all the functions currently supported in the attributeMappingSource.
// Deprecated: This method is obsolete. Use GetAsFunctionsGetResponse instead.
// returns a ItemSynchronizationJobsItemSchemaFunctionsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-functions?view=graph-rest-1.0
func (m *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(ItemSynchronizationJobsItemSchemaFunctionsResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationJobsItemSchemaFunctionsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationJobsItemSchemaFunctionsResponseable), nil
}
// GetAsFunctionsGetResponse list all the functions currently supported in the attributeMappingSource.
// returns a ItemSynchronizationJobsItemSchemaFunctionsGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-synchronizationschema-functions?view=graph-rest-1.0
func (m *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) GetAsFunctionsGetResponse(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(ItemSynchronizationJobsItemSchemaFunctionsGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationJobsItemSchemaFunctionsGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationJobsItemSchemaFunctionsGetResponseable), nil
}
// ToGetRequestInformation list all the functions currently supported in the attributeMappingSource.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder when successful
func (m *ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationJobsItemSchemaFunctionsRequestBuilder) {
    return NewItemSynchronizationJobsItemSchemaFunctionsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
