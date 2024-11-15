package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SecurityGetAttackSimulationRepeatOffendersRequestBuilder provides operations to call the getAttackSimulationRepeatOffenders method.
type SecurityGetAttackSimulationRepeatOffendersRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetQueryParameters list the tenant users who have yielded to attacks more than once in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
type SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetQueryParameters
}
// NewSecurityGetAttackSimulationRepeatOffendersRequestBuilderInternal instantiates a new SecurityGetAttackSimulationRepeatOffendersRequestBuilder and sets the default values.
func NewSecurityGetAttackSimulationRepeatOffendersRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityGetAttackSimulationRepeatOffendersRequestBuilder) {
    m := &SecurityGetAttackSimulationRepeatOffendersRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/security/getAttackSimulationRepeatOffenders(){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewSecurityGetAttackSimulationRepeatOffendersRequestBuilder instantiates a new SecurityGetAttackSimulationRepeatOffendersRequestBuilder and sets the default values.
func NewSecurityGetAttackSimulationRepeatOffendersRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityGetAttackSimulationRepeatOffendersRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSecurityGetAttackSimulationRepeatOffendersRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list the tenant users who have yielded to attacks more than once in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// Deprecated: This method is obsolete. Use GetAsGetAttackSimulationRepeatOffendersGetResponse instead.
// returns a SecurityGetAttackSimulationRepeatOffendersResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securityreportsroot-getattacksimulationrepeatoffenders?view=graph-rest-1.0
func (m *SecurityGetAttackSimulationRepeatOffendersRequestBuilder) Get(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetRequestConfiguration)(SecurityGetAttackSimulationRepeatOffendersResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateSecurityGetAttackSimulationRepeatOffendersResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(SecurityGetAttackSimulationRepeatOffendersResponseable), nil
}
// GetAsGetAttackSimulationRepeatOffendersGetResponse list the tenant users who have yielded to attacks more than once in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// returns a SecurityGetAttackSimulationRepeatOffendersGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securityreportsroot-getattacksimulationrepeatoffenders?view=graph-rest-1.0
func (m *SecurityGetAttackSimulationRepeatOffendersRequestBuilder) GetAsGetAttackSimulationRepeatOffendersGetResponse(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetRequestConfiguration)(SecurityGetAttackSimulationRepeatOffendersGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateSecurityGetAttackSimulationRepeatOffendersGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(SecurityGetAttackSimulationRepeatOffendersGetResponseable), nil
}
// ToGetRequestInformation list the tenant users who have yielded to attacks more than once in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// returns a *RequestInformation when successful
func (m *SecurityGetAttackSimulationRepeatOffendersRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationRepeatOffendersRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SecurityGetAttackSimulationRepeatOffendersRequestBuilder when successful
func (m *SecurityGetAttackSimulationRepeatOffendersRequestBuilder) WithUrl(rawUrl string)(*SecurityGetAttackSimulationRepeatOffendersRequestBuilder) {
    return NewSecurityGetAttackSimulationRepeatOffendersRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
