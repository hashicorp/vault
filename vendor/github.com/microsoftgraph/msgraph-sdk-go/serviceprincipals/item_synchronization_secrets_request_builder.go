package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSynchronizationSecretsRequestBuilder builds and executes requests for operations under \servicePrincipals\{servicePrincipal-id}\synchronization\secrets
type ItemSynchronizationSecretsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSynchronizationSecretsRequestBuilderPutRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSynchronizationSecretsRequestBuilderPutRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSynchronizationSecretsRequestBuilderInternal instantiates a new ItemSynchronizationSecretsRequestBuilder and sets the default values.
func NewItemSynchronizationSecretsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationSecretsRequestBuilder) {
    m := &ItemSynchronizationSecretsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/synchronization/secrets", pathParameters),
    }
    return m
}
// NewItemSynchronizationSecretsRequestBuilder instantiates a new ItemSynchronizationSecretsRequestBuilder and sets the default values.
func NewItemSynchronizationSecretsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSynchronizationSecretsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSynchronizationSecretsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *ItemSynchronizationSecretsCountRequestBuilder when successful
func (m *ItemSynchronizationSecretsRequestBuilder) Count()(*ItemSynchronizationSecretsCountRequestBuilder) {
    return NewItemSynchronizationSecretsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Put provide credentials for establishing connectivity with the target system.
// Deprecated: This method is obsolete. Use PutAsSecretsPutResponse instead.
// returns a ItemSynchronizationSecretsResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-serviceprincipal-put-synchronization?view=graph-rest-1.0
func (m *ItemSynchronizationSecretsRequestBuilder) Put(ctx context.Context, body ItemSynchronizationSecretsPutRequestBodyable, requestConfiguration *ItemSynchronizationSecretsRequestBuilderPutRequestConfiguration)(ItemSynchronizationSecretsResponseable, error) {
    requestInfo, err := m.ToPutRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationSecretsResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationSecretsResponseable), nil
}
// PutAsSecretsPutResponse provide credentials for establishing connectivity with the target system.
// returns a ItemSynchronizationSecretsPutResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/synchronization-serviceprincipal-put-synchronization?view=graph-rest-1.0
func (m *ItemSynchronizationSecretsRequestBuilder) PutAsSecretsPutResponse(ctx context.Context, body ItemSynchronizationSecretsPutRequestBodyable, requestConfiguration *ItemSynchronizationSecretsRequestBuilderPutRequestConfiguration)(ItemSynchronizationSecretsPutResponseable, error) {
    requestInfo, err := m.ToPutRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSynchronizationSecretsPutResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSynchronizationSecretsPutResponseable), nil
}
// ToPutRequestInformation provide credentials for establishing connectivity with the target system.
// returns a *RequestInformation when successful
func (m *ItemSynchronizationSecretsRequestBuilder) ToPutRequestInformation(ctx context.Context, body ItemSynchronizationSecretsPutRequestBodyable, requestConfiguration *ItemSynchronizationSecretsRequestBuilderPutRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PUT, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *ItemSynchronizationSecretsRequestBuilder when successful
func (m *ItemSynchronizationSecretsRequestBuilder) WithUrl(rawUrl string)(*ItemSynchronizationSecretsRequestBuilder) {
    return NewItemSynchronizationSecretsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
