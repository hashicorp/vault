package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AttackSimulationSimulationsItemLoginPageRequestBuilder provides operations to manage the loginPage property of the microsoft.graph.simulation entity.
type AttackSimulationSimulationsItemLoginPageRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AttackSimulationSimulationsItemLoginPageRequestBuilderGetQueryParameters the login page associated with a simulation during its creation.
type AttackSimulationSimulationsItemLoginPageRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// AttackSimulationSimulationsItemLoginPageRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationSimulationsItemLoginPageRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AttackSimulationSimulationsItemLoginPageRequestBuilderGetQueryParameters
}
// NewAttackSimulationSimulationsItemLoginPageRequestBuilderInternal instantiates a new AttackSimulationSimulationsItemLoginPageRequestBuilder and sets the default values.
func NewAttackSimulationSimulationsItemLoginPageRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationSimulationsItemLoginPageRequestBuilder) {
    m := &AttackSimulationSimulationsItemLoginPageRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/attackSimulation/simulations/{simulation%2Did}/loginPage{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewAttackSimulationSimulationsItemLoginPageRequestBuilder instantiates a new AttackSimulationSimulationsItemLoginPageRequestBuilder and sets the default values.
func NewAttackSimulationSimulationsItemLoginPageRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationSimulationsItemLoginPageRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAttackSimulationSimulationsItemLoginPageRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the login page associated with a simulation during its creation.
// returns a LoginPageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationSimulationsItemLoginPageRequestBuilder) Get(ctx context.Context, requestConfiguration *AttackSimulationSimulationsItemLoginPageRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LoginPageable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateLoginPageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LoginPageable), nil
}
// ToGetRequestInformation the login page associated with a simulation during its creation.
// returns a *RequestInformation when successful
func (m *AttackSimulationSimulationsItemLoginPageRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationSimulationsItemLoginPageRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *AttackSimulationSimulationsItemLoginPageRequestBuilder when successful
func (m *AttackSimulationSimulationsItemLoginPageRequestBuilder) WithUrl(rawUrl string)(*AttackSimulationSimulationsItemLoginPageRequestBuilder) {
    return NewAttackSimulationSimulationsItemLoginPageRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
