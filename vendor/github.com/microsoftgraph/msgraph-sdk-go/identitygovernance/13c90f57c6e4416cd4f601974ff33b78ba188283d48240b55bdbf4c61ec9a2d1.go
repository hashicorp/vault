package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder provides operations to manage the customExtensionStageSettings property of the microsoft.graph.accessPackageAssignmentPolicy entity.
type EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetQueryParameters the collection of stages when to execute one or more custom access package workflow extensions. Supports $expand.
type EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetQueryParameters struct {
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
// EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetQueryParameters
}
// EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByCustomExtensionStageSettingId provides operations to manage the customExtensionStageSettings property of the microsoft.graph.accessPackageAssignmentPolicy entity.
// returns a *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCustomExtensionStageSettingItemRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) ByCustomExtensionStageSettingId(customExtensionStageSettingId string)(*EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCustomExtensionStageSettingItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if customExtensionStageSettingId != "" {
        urlTplParams["customExtensionStageSetting%2Did"] = customExtensionStageSettingId
    }
    return NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCustomExtensionStageSettingItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderInternal instantiates a new EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder and sets the default values.
func NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) {
    m := &EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/entitlementManagement/accessPackages/{accessPackage%2Did}/assignmentPolicies/{accessPackageAssignmentPolicy%2Did}/customExtensionStageSettings{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder instantiates a new EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder and sets the default values.
func NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCountRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) Count()(*EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCountRequestBuilder) {
    return NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the collection of stages when to execute one or more custom access package workflow extensions. Supports $expand.
// returns a CustomExtensionStageSettingCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) Get(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomExtensionStageSettingCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingCollectionResponseable), nil
}
// Post create new navigation property to customExtensionStageSettings for identityGovernance
// returns a CustomExtensionStageSettingable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingable, requestConfiguration *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomExtensionStageSettingFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingable), nil
}
// ToGetRequestInformation the collection of stages when to execute one or more custom access package workflow extensions. Supports $expand.
// returns a *RequestInformation when successful
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to customExtensionStageSettings for identityGovernance
// returns a *RequestInformation when successful
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomExtensionStageSettingable, requestConfiguration *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder when successful
func (m *EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) WithUrl(rawUrl string)(*EntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder) {
    return NewEntitlementManagementAccessPackagesItemAssignmentPoliciesItemCustomExtensionStageSettingsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
