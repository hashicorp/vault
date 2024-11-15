package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder provides operations to call the getAttackSimulationTrainingUserCoverage method.
type SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetQueryParameters list training coverage for tenant users in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
type SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetQueryParameters struct {
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
// SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetQueryParameters
}
// NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilderInternal instantiates a new SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder and sets the default values.
func NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) {
    m := &SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/security/getAttackSimulationTrainingUserCoverage(){?%24count,%24filter,%24search,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilder instantiates a new SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder and sets the default values.
func NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilderInternal(urlParams, requestAdapter)
}
// Get list training coverage for tenant users in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// Deprecated: This method is obsolete. Use GetAsGetAttackSimulationTrainingUserCoverageGetResponse instead.
// returns a SecurityGetAttackSimulationTrainingUserCoverageResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securityreportsroot-getattacksimulationtrainingusercoverage?view=graph-rest-1.0
func (m *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) Get(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetRequestConfiguration)(SecurityGetAttackSimulationTrainingUserCoverageResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateSecurityGetAttackSimulationTrainingUserCoverageResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(SecurityGetAttackSimulationTrainingUserCoverageResponseable), nil
}
// GetAsGetAttackSimulationTrainingUserCoverageGetResponse list training coverage for tenant users in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// returns a SecurityGetAttackSimulationTrainingUserCoverageGetResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securityreportsroot-getattacksimulationtrainingusercoverage?view=graph-rest-1.0
func (m *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) GetAsGetAttackSimulationTrainingUserCoverageGetResponse(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetRequestConfiguration)(SecurityGetAttackSimulationTrainingUserCoverageGetResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateSecurityGetAttackSimulationTrainingUserCoverageGetResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(SecurityGetAttackSimulationTrainingUserCoverageGetResponseable), nil
}
// ToGetRequestInformation list training coverage for tenant users in attack simulation and training campaigns. This function supports @odata.nextLink for pagination.
// returns a *RequestInformation when successful
func (m *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder when successful
func (m *SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) WithUrl(rawUrl string)(*SecurityGetAttackSimulationTrainingUserCoverageRequestBuilder) {
    return NewSecurityGetAttackSimulationTrainingUserCoverageRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
