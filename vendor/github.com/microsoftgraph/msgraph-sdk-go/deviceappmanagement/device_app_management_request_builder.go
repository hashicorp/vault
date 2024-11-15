package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeviceAppManagementRequestBuilder provides operations to manage the deviceAppManagement singleton.
type DeviceAppManagementRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeviceAppManagementRequestBuilderGetQueryParameters read properties and relationships of the deviceAppManagement object.
type DeviceAppManagementRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DeviceAppManagementRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceAppManagementRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeviceAppManagementRequestBuilderGetQueryParameters
}
// DeviceAppManagementRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceAppManagementRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AndroidManagedAppProtections provides operations to manage the androidManagedAppProtections property of the microsoft.graph.deviceAppManagement entity.
// returns a *AndroidManagedAppProtectionsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) AndroidManagedAppProtections()(*AndroidManagedAppProtectionsRequestBuilder) {
    return NewAndroidManagedAppProtectionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeviceAppManagementRequestBuilderInternal instantiates a new DeviceAppManagementRequestBuilder and sets the default values.
func NewDeviceAppManagementRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceAppManagementRequestBuilder) {
    m := &DeviceAppManagementRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDeviceAppManagementRequestBuilder instantiates a new DeviceAppManagementRequestBuilder and sets the default values.
func NewDeviceAppManagementRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceAppManagementRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeviceAppManagementRequestBuilderInternal(urlParams, requestAdapter)
}
// DefaultManagedAppProtections provides operations to manage the defaultManagedAppProtections property of the microsoft.graph.deviceAppManagement entity.
// returns a *DefaultManagedAppProtectionsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) DefaultManagedAppProtections()(*DefaultManagedAppProtectionsRequestBuilder) {
    return NewDefaultManagedAppProtectionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get read properties and relationships of the deviceAppManagement object.
// returns a DeviceAppManagementable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-apps-deviceappmanagement-get?view=graph-rest-1.0
func (m *DeviceAppManagementRequestBuilder) Get(ctx context.Context, requestConfiguration *DeviceAppManagementRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceAppManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable), nil
}
// IosManagedAppProtections provides operations to manage the iosManagedAppProtections property of the microsoft.graph.deviceAppManagement entity.
// returns a *IosManagedAppProtectionsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) IosManagedAppProtections()(*IosManagedAppProtectionsRequestBuilder) {
    return NewIosManagedAppProtectionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedAppPolicies provides operations to manage the managedAppPolicies property of the microsoft.graph.deviceAppManagement entity.
// returns a *ManagedAppPoliciesRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) ManagedAppPolicies()(*ManagedAppPoliciesRequestBuilder) {
    return NewManagedAppPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedAppRegistrations provides operations to manage the managedAppRegistrations property of the microsoft.graph.deviceAppManagement entity.
// returns a *ManagedAppRegistrationsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) ManagedAppRegistrations()(*ManagedAppRegistrationsRequestBuilder) {
    return NewManagedAppRegistrationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedAppStatuses provides operations to manage the managedAppStatuses property of the microsoft.graph.deviceAppManagement entity.
// returns a *ManagedAppStatusesRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) ManagedAppStatuses()(*ManagedAppStatusesRequestBuilder) {
    return NewManagedAppStatusesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedEBooks provides operations to manage the managedEBooks property of the microsoft.graph.deviceAppManagement entity.
// returns a *ManagedEBooksRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) ManagedEBooks()(*ManagedEBooksRequestBuilder) {
    return NewManagedEBooksRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MdmWindowsInformationProtectionPolicies provides operations to manage the mdmWindowsInformationProtectionPolicies property of the microsoft.graph.deviceAppManagement entity.
// returns a *MdmWindowsInformationProtectionPoliciesRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) MdmWindowsInformationProtectionPolicies()(*MdmWindowsInformationProtectionPoliciesRequestBuilder) {
    return NewMdmWindowsInformationProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MobileAppCategories provides operations to manage the mobileAppCategories property of the microsoft.graph.deviceAppManagement entity.
// returns a *MobileAppCategoriesRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) MobileAppCategories()(*MobileAppCategoriesRequestBuilder) {
    return NewMobileAppCategoriesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MobileAppConfigurations provides operations to manage the mobileAppConfigurations property of the microsoft.graph.deviceAppManagement entity.
// returns a *MobileAppConfigurationsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) MobileAppConfigurations()(*MobileAppConfigurationsRequestBuilder) {
    return NewMobileAppConfigurationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MobileApps provides operations to manage the mobileApps property of the microsoft.graph.deviceAppManagement entity.
// returns a *MobileAppsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) MobileApps()(*MobileAppsRequestBuilder) {
    return NewMobileAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the properties of a deviceAppManagement object.
// returns a DeviceAppManagementable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-books-deviceappmanagement-update?view=graph-rest-1.0
func (m *DeviceAppManagementRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable, requestConfiguration *DeviceAppManagementRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceAppManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable), nil
}
// SyncMicrosoftStoreForBusinessApps provides operations to call the syncMicrosoftStoreForBusinessApps method.
// returns a *SyncMicrosoftStoreForBusinessAppsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) SyncMicrosoftStoreForBusinessApps()(*SyncMicrosoftStoreForBusinessAppsRequestBuilder) {
    return NewSyncMicrosoftStoreForBusinessAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TargetedManagedAppConfigurations provides operations to manage the targetedManagedAppConfigurations property of the microsoft.graph.deviceAppManagement entity.
// returns a *TargetedManagedAppConfigurationsRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) TargetedManagedAppConfigurations()(*TargetedManagedAppConfigurationsRequestBuilder) {
    return NewTargetedManagedAppConfigurationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation read properties and relationships of the deviceAppManagement object.
// returns a *RequestInformation when successful
func (m *DeviceAppManagementRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeviceAppManagementRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a deviceAppManagement object.
// returns a *RequestInformation when successful
func (m *DeviceAppManagementRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceAppManagementable, requestConfiguration *DeviceAppManagementRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// VppTokens provides operations to manage the vppTokens property of the microsoft.graph.deviceAppManagement entity.
// returns a *VppTokensRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) VppTokens()(*VppTokensRequestBuilder) {
    return NewVppTokensRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsInformationProtectionPolicies provides operations to manage the windowsInformationProtectionPolicies property of the microsoft.graph.deviceAppManagement entity.
// returns a *WindowsInformationProtectionPoliciesRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) WindowsInformationProtectionPolicies()(*WindowsInformationProtectionPoliciesRequestBuilder) {
    return NewWindowsInformationProtectionPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *DeviceAppManagementRequestBuilder when successful
func (m *DeviceAppManagementRequestBuilder) WithUrl(rawUrl string)(*DeviceAppManagementRequestBuilder) {
    return NewDeviceAppManagementRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
