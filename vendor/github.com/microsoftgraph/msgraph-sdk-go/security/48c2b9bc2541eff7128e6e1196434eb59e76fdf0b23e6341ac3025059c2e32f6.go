package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder provides operations to manage the details property of the microsoft.graph.endUserNotification entity.
type AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetQueryParameters get details from security
type AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetQueryParameters
}
// AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderInternal instantiates a new AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder and sets the default values.
func NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) {
    m := &AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/attackSimulation/endUserNotifications/{endUserNotification%2Did}/details/{endUserNotificationDetail%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder instantiates a new AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder and sets the default values.
func NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property details for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get details from security
// returns a EndUserNotificationDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) Get(ctx context.Context, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEndUserNotificationDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable), nil
}
// Patch update the navigation property details in security
// returns a EndUserNotificationDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateEndUserNotificationDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable), nil
}
// ToDeleteRequestInformation delete navigation property details for security
// returns a *RequestInformation when successful
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get details from security
// returns a *RequestInformation when successful
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property details in security
// returns a *RequestInformation when successful
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.EndUserNotificationDetailable, requestConfiguration *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder when successful
func (m *AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) WithUrl(rawUrl string)(*AttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder) {
    return NewAttackSimulationEndUserNotificationsItemDetailsEndUserNotificationDetailItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
