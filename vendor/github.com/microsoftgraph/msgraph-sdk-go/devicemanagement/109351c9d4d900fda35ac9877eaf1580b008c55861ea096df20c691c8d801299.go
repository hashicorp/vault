package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder provides operations to manage the scheduledActionConfigurations property of the microsoft.graph.deviceComplianceScheduledActionForRule entity.
type DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetQueryParameters list properties and relationships of the deviceComplianceActionItem objects.
type DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetQueryParameters struct {
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
// DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetQueryParameters
}
// DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByDeviceComplianceActionItemId provides operations to manage the scheduledActionConfigurations property of the microsoft.graph.deviceComplianceScheduledActionForRule entity.
// returns a *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsDeviceComplianceActionItemItemRequestBuilder when successful
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) ByDeviceComplianceActionItemId(deviceComplianceActionItemId string)(*DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsDeviceComplianceActionItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if deviceComplianceActionItemId != "" {
        urlTplParams["deviceComplianceActionItem%2Did"] = deviceComplianceActionItemId
    }
    return NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsDeviceComplianceActionItemItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderInternal instantiates a new DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder and sets the default values.
func NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) {
    m := &DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/deviceCompliancePolicies/{deviceCompliancePolicy%2Did}/scheduledActionsForRule/{deviceComplianceScheduledActionForRule%2Did}/scheduledActionConfigurations{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder instantiates a new DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder and sets the default values.
func NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsCountRequestBuilder when successful
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) Count()(*DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsCountRequestBuilder) {
    return NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get list properties and relationships of the deviceComplianceActionItem objects.
// returns a DeviceComplianceActionItemCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-deviceconfig-devicecomplianceactionitem-list?view=graph-rest-1.0
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) Get(ctx context.Context, requestConfiguration *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceComplianceActionItemCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemCollectionResponseable), nil
}
// Post create a new deviceComplianceActionItem object.
// returns a DeviceComplianceActionItemable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-deviceconfig-devicecomplianceactionitem-create?view=graph-rest-1.0
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemable, requestConfiguration *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceComplianceActionItemFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemable), nil
}
// ToGetRequestInformation list properties and relationships of the deviceComplianceActionItem objects.
// returns a *RequestInformation when successful
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create a new deviceComplianceActionItem object.
// returns a *RequestInformation when successful
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceComplianceActionItemable, requestConfiguration *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder when successful
func (m *DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) WithUrl(rawUrl string)(*DeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder) {
    return NewDeviceCompliancePoliciesItemScheduledActionsForRuleItemScheduledActionConfigurationsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
