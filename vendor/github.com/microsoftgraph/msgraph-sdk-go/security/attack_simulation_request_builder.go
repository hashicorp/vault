package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AttackSimulationRequestBuilder provides operations to manage the attackSimulation property of the microsoft.graph.security entity.
type AttackSimulationRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AttackSimulationRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AttackSimulationRequestBuilderGetQueryParameters get attackSimulation from security
type AttackSimulationRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// AttackSimulationRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AttackSimulationRequestBuilderGetQueryParameters
}
// AttackSimulationRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAttackSimulationRequestBuilderInternal instantiates a new AttackSimulationRequestBuilder and sets the default values.
func NewAttackSimulationRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationRequestBuilder) {
    m := &AttackSimulationRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/attackSimulation{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewAttackSimulationRequestBuilder instantiates a new AttackSimulationRequestBuilder and sets the default values.
func NewAttackSimulationRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAttackSimulationRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property attackSimulation for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationRequestBuilder) Delete(ctx context.Context, requestConfiguration *AttackSimulationRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// EndUserNotifications provides operations to manage the endUserNotifications property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationEndUserNotificationsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) EndUserNotifications()(*AttackSimulationEndUserNotificationsRequestBuilder) {
    return NewAttackSimulationEndUserNotificationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get attackSimulation from security
// returns a AttackSimulationRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationRequestBuilder) Get(ctx context.Context, requestConfiguration *AttackSimulationRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAttackSimulationRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable), nil
}
// LandingPages provides operations to manage the landingPages property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationLandingPagesRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) LandingPages()(*AttackSimulationLandingPagesRequestBuilder) {
    return NewAttackSimulationLandingPagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LoginPages provides operations to manage the loginPages property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationLoginPagesRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) LoginPages()(*AttackSimulationLoginPagesRequestBuilder) {
    return NewAttackSimulationLoginPagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Operations provides operations to manage the operations property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationOperationsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) Operations()(*AttackSimulationOperationsRequestBuilder) {
    return NewAttackSimulationOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property attackSimulation in security
// returns a AttackSimulationRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable, requestConfiguration *AttackSimulationRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAttackSimulationRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable), nil
}
// Payloads provides operations to manage the payloads property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationPayloadsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) Payloads()(*AttackSimulationPayloadsRequestBuilder) {
    return NewAttackSimulationPayloadsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SimulationAutomations provides operations to manage the simulationAutomations property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationSimulationAutomationsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) SimulationAutomations()(*AttackSimulationSimulationAutomationsRequestBuilder) {
    return NewAttackSimulationSimulationAutomationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Simulations provides operations to manage the simulations property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationSimulationsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) Simulations()(*AttackSimulationSimulationsRequestBuilder) {
    return NewAttackSimulationSimulationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation delete navigation property attackSimulation for security
// returns a *RequestInformation when successful
func (m *AttackSimulationRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get attackSimulation from security
// returns a *RequestInformation when successful
func (m *AttackSimulationRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property attackSimulation in security
// returns a *RequestInformation when successful
func (m *AttackSimulationRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttackSimulationRootable, requestConfiguration *AttackSimulationRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// Trainings provides operations to manage the trainings property of the microsoft.graph.attackSimulationRoot entity.
// returns a *AttackSimulationTrainingsRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) Trainings()(*AttackSimulationTrainingsRequestBuilder) {
    return NewAttackSimulationTrainingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AttackSimulationRequestBuilder when successful
func (m *AttackSimulationRequestBuilder) WithUrl(rawUrl string)(*AttackSimulationRequestBuilder) {
    return NewAttackSimulationRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
