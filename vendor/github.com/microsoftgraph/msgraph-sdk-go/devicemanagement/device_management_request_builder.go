package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeviceManagementRequestBuilder provides operations to manage the deviceManagement singleton.
type DeviceManagementRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeviceManagementRequestBuilderGetQueryParameters read properties and relationships of the deviceManagement object.
type DeviceManagementRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DeviceManagementRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceManagementRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeviceManagementRequestBuilderGetQueryParameters
}
// DeviceManagementRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceManagementRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ApplePushNotificationCertificate provides operations to manage the applePushNotificationCertificate property of the microsoft.graph.deviceManagement entity.
// returns a *ApplePushNotificationCertificateRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ApplePushNotificationCertificate()(*ApplePushNotificationCertificateRequestBuilder) {
    return NewApplePushNotificationCertificateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuditEvents provides operations to manage the auditEvents property of the microsoft.graph.deviceManagement entity.
// returns a *AuditEventsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) AuditEvents()(*AuditEventsRequestBuilder) {
    return NewAuditEventsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ComplianceManagementPartners provides operations to manage the complianceManagementPartners property of the microsoft.graph.deviceManagement entity.
// returns a *ComplianceManagementPartnersRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ComplianceManagementPartners()(*ComplianceManagementPartnersRequestBuilder) {
    return NewComplianceManagementPartnersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ConditionalAccessSettings provides operations to manage the conditionalAccessSettings property of the microsoft.graph.deviceManagement entity.
// returns a *ConditionalAccessSettingsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ConditionalAccessSettings()(*ConditionalAccessSettingsRequestBuilder) {
    return NewConditionalAccessSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewDeviceManagementRequestBuilderInternal instantiates a new DeviceManagementRequestBuilder and sets the default values.
func NewDeviceManagementRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceManagementRequestBuilder) {
    m := &DeviceManagementRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDeviceManagementRequestBuilder instantiates a new DeviceManagementRequestBuilder and sets the default values.
func NewDeviceManagementRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceManagementRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeviceManagementRequestBuilderInternal(urlParams, requestAdapter)
}
// DetectedApps provides operations to manage the detectedApps property of the microsoft.graph.deviceManagement entity.
// returns a *DetectedAppsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DetectedApps()(*DetectedAppsRequestBuilder) {
    return NewDetectedAppsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCategories provides operations to manage the deviceCategories property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceCategoriesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceCategories()(*DeviceCategoriesRequestBuilder) {
    return NewDeviceCategoriesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCompliancePolicies provides operations to manage the deviceCompliancePolicies property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceCompliancePoliciesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceCompliancePolicies()(*DeviceCompliancePoliciesRequestBuilder) {
    return NewDeviceCompliancePoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCompliancePolicyDeviceStateSummary provides operations to manage the deviceCompliancePolicyDeviceStateSummary property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceCompliancePolicyDeviceStateSummaryRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceCompliancePolicyDeviceStateSummary()(*DeviceCompliancePolicyDeviceStateSummaryRequestBuilder) {
    return NewDeviceCompliancePolicyDeviceStateSummaryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCompliancePolicySettingStateSummaries provides operations to manage the deviceCompliancePolicySettingStateSummaries property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceCompliancePolicySettingStateSummariesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceCompliancePolicySettingStateSummaries()(*DeviceCompliancePolicySettingStateSummariesRequestBuilder) {
    return NewDeviceCompliancePolicySettingStateSummariesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceConfigurationDeviceStateSummaries provides operations to manage the deviceConfigurationDeviceStateSummaries property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceConfigurationDeviceStateSummariesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceConfigurationDeviceStateSummaries()(*DeviceConfigurationDeviceStateSummariesRequestBuilder) {
    return NewDeviceConfigurationDeviceStateSummariesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceConfigurations provides operations to manage the deviceConfigurations property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceConfigurationsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceConfigurations()(*DeviceConfigurationsRequestBuilder) {
    return NewDeviceConfigurationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceEnrollmentConfigurations provides operations to manage the deviceEnrollmentConfigurations property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceEnrollmentConfigurationsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceEnrollmentConfigurations()(*DeviceEnrollmentConfigurationsRequestBuilder) {
    return NewDeviceEnrollmentConfigurationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceManagementPartners provides operations to manage the deviceManagementPartners property of the microsoft.graph.deviceManagement entity.
// returns a *DeviceManagementPartnersRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) DeviceManagementPartners()(*DeviceManagementPartnersRequestBuilder) {
    return NewDeviceManagementPartnersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExchangeConnectors provides operations to manage the exchangeConnectors property of the microsoft.graph.deviceManagement entity.
// returns a *ExchangeConnectorsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ExchangeConnectors()(*ExchangeConnectorsRequestBuilder) {
    return NewExchangeConnectorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get read properties and relationships of the deviceManagement object.
// returns a DeviceManagementable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-tem-devicemanagement-get?view=graph-rest-1.0
func (m *DeviceManagementRequestBuilder) Get(ctx context.Context, requestConfiguration *DeviceManagementRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable), nil
}
// GetEffectivePermissionsWithScope provides operations to call the getEffectivePermissions method.
// returns a *GetEffectivePermissionsWithScopeRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) GetEffectivePermissionsWithScope(scope *string)(*GetEffectivePermissionsWithScopeRequestBuilder) {
    return NewGetEffectivePermissionsWithScopeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, scope)
}
// ImportedWindowsAutopilotDeviceIdentities provides operations to manage the importedWindowsAutopilotDeviceIdentities property of the microsoft.graph.deviceManagement entity.
// returns a *ImportedWindowsAutopilotDeviceIdentitiesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ImportedWindowsAutopilotDeviceIdentities()(*ImportedWindowsAutopilotDeviceIdentitiesRequestBuilder) {
    return NewImportedWindowsAutopilotDeviceIdentitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IosUpdateStatuses provides operations to manage the iosUpdateStatuses property of the microsoft.graph.deviceManagement entity.
// returns a *IosUpdateStatusesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) IosUpdateStatuses()(*IosUpdateStatusesRequestBuilder) {
    return NewIosUpdateStatusesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedDeviceOverview provides operations to manage the managedDeviceOverview property of the microsoft.graph.deviceManagement entity.
// returns a *ManagedDeviceOverviewRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ManagedDeviceOverview()(*ManagedDeviceOverviewRequestBuilder) {
    return NewManagedDeviceOverviewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedDevices provides operations to manage the managedDevices property of the microsoft.graph.deviceManagement entity.
// returns a *ManagedDevicesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ManagedDevices()(*ManagedDevicesRequestBuilder) {
    return NewManagedDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MobileAppTroubleshootingEvents provides operations to manage the mobileAppTroubleshootingEvents property of the microsoft.graph.deviceManagement entity.
// returns a *MobileAppTroubleshootingEventsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) MobileAppTroubleshootingEvents()(*MobileAppTroubleshootingEventsRequestBuilder) {
    return NewMobileAppTroubleshootingEventsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MobileThreatDefenseConnectors provides operations to manage the mobileThreatDefenseConnectors property of the microsoft.graph.deviceManagement entity.
// returns a *MobileThreatDefenseConnectorsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) MobileThreatDefenseConnectors()(*MobileThreatDefenseConnectorsRequestBuilder) {
    return NewMobileThreatDefenseConnectorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NotificationMessageTemplates provides operations to manage the notificationMessageTemplates property of the microsoft.graph.deviceManagement entity.
// returns a *NotificationMessageTemplatesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) NotificationMessageTemplates()(*NotificationMessageTemplatesRequestBuilder) {
    return NewNotificationMessageTemplatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the properties of a deviceManagement object.
// returns a DeviceManagementable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-tem-devicemanagement-update?view=graph-rest-1.0
func (m *DeviceManagementRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable, requestConfiguration *DeviceManagementRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceManagementFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable), nil
}
// RemoteAssistancePartners provides operations to manage the remoteAssistancePartners property of the microsoft.graph.deviceManagement entity.
// returns a *RemoteAssistancePartnersRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) RemoteAssistancePartners()(*RemoteAssistancePartnersRequestBuilder) {
    return NewRemoteAssistancePartnersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Reports provides operations to manage the reports property of the microsoft.graph.deviceManagement entity.
// returns a *ReportsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) Reports()(*ReportsRequestBuilder) {
    return NewReportsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ResourceOperations provides operations to manage the resourceOperations property of the microsoft.graph.deviceManagement entity.
// returns a *ResourceOperationsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) ResourceOperations()(*ResourceOperationsRequestBuilder) {
    return NewResourceOperationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoleAssignments provides operations to manage the roleAssignments property of the microsoft.graph.deviceManagement entity.
// returns a *RoleAssignmentsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) RoleAssignments()(*RoleAssignmentsRequestBuilder) {
    return NewRoleAssignmentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoleDefinitions provides operations to manage the roleDefinitions property of the microsoft.graph.deviceManagement entity.
// returns a *RoleDefinitionsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) RoleDefinitions()(*RoleDefinitionsRequestBuilder) {
    return NewRoleDefinitionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SoftwareUpdateStatusSummary provides operations to manage the softwareUpdateStatusSummary property of the microsoft.graph.deviceManagement entity.
// returns a *SoftwareUpdateStatusSummaryRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) SoftwareUpdateStatusSummary()(*SoftwareUpdateStatusSummaryRequestBuilder) {
    return NewSoftwareUpdateStatusSummaryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TelecomExpenseManagementPartners provides operations to manage the telecomExpenseManagementPartners property of the microsoft.graph.deviceManagement entity.
// returns a *TelecomExpenseManagementPartnersRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) TelecomExpenseManagementPartners()(*TelecomExpenseManagementPartnersRequestBuilder) {
    return NewTelecomExpenseManagementPartnersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TermsAndConditions provides operations to manage the termsAndConditions property of the microsoft.graph.deviceManagement entity.
// returns a *TermsAndConditionsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) TermsAndConditions()(*TermsAndConditionsRequestBuilder) {
    return NewTermsAndConditionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation read properties and relationships of the deviceManagement object.
// returns a *RequestInformation when successful
func (m *DeviceManagementRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeviceManagementRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a deviceManagement object.
// returns a *RequestInformation when successful
func (m *DeviceManagementRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceManagementable, requestConfiguration *DeviceManagementRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// TroubleshootingEvents provides operations to manage the troubleshootingEvents property of the microsoft.graph.deviceManagement entity.
// returns a *TroubleshootingEventsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) TroubleshootingEvents()(*TroubleshootingEventsRequestBuilder) {
    return NewTroubleshootingEventsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthApplicationPerformance provides operations to manage the userExperienceAnalyticsAppHealthApplicationPerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthApplicationPerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthApplicationPerformance()(*UserExperienceAnalyticsAppHealthApplicationPerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthApplicationPerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetails provides operations to manage the userExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetails property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetailsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetails()(*UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetailsRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDetailsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceId provides operations to manage the userExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceId property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceIdRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceId()(*UserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceIdRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthApplicationPerformanceByAppVersionDeviceIdRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersion provides operations to manage the userExperienceAnalyticsAppHealthApplicationPerformanceByOSVersion property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersionRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersion()(*UserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersionRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthApplicationPerformanceByOSVersionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthDeviceModelPerformance provides operations to manage the userExperienceAnalyticsAppHealthDeviceModelPerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthDeviceModelPerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthDeviceModelPerformance()(*UserExperienceAnalyticsAppHealthDeviceModelPerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthDeviceModelPerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthDevicePerformance provides operations to manage the userExperienceAnalyticsAppHealthDevicePerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthDevicePerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthDevicePerformance()(*UserExperienceAnalyticsAppHealthDevicePerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthDevicePerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthDevicePerformanceDetails provides operations to manage the userExperienceAnalyticsAppHealthDevicePerformanceDetails property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthDevicePerformanceDetailsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthDevicePerformanceDetails()(*UserExperienceAnalyticsAppHealthDevicePerformanceDetailsRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthDevicePerformanceDetailsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthOSVersionPerformance provides operations to manage the userExperienceAnalyticsAppHealthOSVersionPerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthOSVersionPerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthOSVersionPerformance()(*UserExperienceAnalyticsAppHealthOSVersionPerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthOSVersionPerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsAppHealthOverview provides operations to manage the userExperienceAnalyticsAppHealthOverview property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsAppHealthOverviewRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsAppHealthOverview()(*UserExperienceAnalyticsAppHealthOverviewRequestBuilder) {
    return NewUserExperienceAnalyticsAppHealthOverviewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsBaselines provides operations to manage the userExperienceAnalyticsBaselines property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsBaselinesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsBaselines()(*UserExperienceAnalyticsBaselinesRequestBuilder) {
    return NewUserExperienceAnalyticsBaselinesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsCategories provides operations to manage the userExperienceAnalyticsCategories property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsCategoriesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsCategories()(*UserExperienceAnalyticsCategoriesRequestBuilder) {
    return NewUserExperienceAnalyticsCategoriesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsDevicePerformance provides operations to manage the userExperienceAnalyticsDevicePerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsDevicePerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsDevicePerformance()(*UserExperienceAnalyticsDevicePerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsDevicePerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsDeviceScores provides operations to manage the userExperienceAnalyticsDeviceScores property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsDeviceScoresRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsDeviceScores()(*UserExperienceAnalyticsDeviceScoresRequestBuilder) {
    return NewUserExperienceAnalyticsDeviceScoresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsDeviceStartupHistory provides operations to manage the userExperienceAnalyticsDeviceStartupHistory property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsDeviceStartupHistoryRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsDeviceStartupHistory()(*UserExperienceAnalyticsDeviceStartupHistoryRequestBuilder) {
    return NewUserExperienceAnalyticsDeviceStartupHistoryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsDeviceStartupProcesses provides operations to manage the userExperienceAnalyticsDeviceStartupProcesses property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsDeviceStartupProcessesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsDeviceStartupProcesses()(*UserExperienceAnalyticsDeviceStartupProcessesRequestBuilder) {
    return NewUserExperienceAnalyticsDeviceStartupProcessesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsDeviceStartupProcessPerformance provides operations to manage the userExperienceAnalyticsDeviceStartupProcessPerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsDeviceStartupProcessPerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsDeviceStartupProcessPerformance()(*UserExperienceAnalyticsDeviceStartupProcessPerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsDeviceStartupProcessPerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsMetricHistory provides operations to manage the userExperienceAnalyticsMetricHistory property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsMetricHistoryRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsMetricHistory()(*UserExperienceAnalyticsMetricHistoryRequestBuilder) {
    return NewUserExperienceAnalyticsMetricHistoryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsModelScores provides operations to manage the userExperienceAnalyticsModelScores property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsModelScoresRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsModelScores()(*UserExperienceAnalyticsModelScoresRequestBuilder) {
    return NewUserExperienceAnalyticsModelScoresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsOverview provides operations to manage the userExperienceAnalyticsOverview property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsOverviewRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsOverview()(*UserExperienceAnalyticsOverviewRequestBuilder) {
    return NewUserExperienceAnalyticsOverviewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsScoreHistory provides operations to manage the userExperienceAnalyticsScoreHistory property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsScoreHistoryRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsScoreHistory()(*UserExperienceAnalyticsScoreHistoryRequestBuilder) {
    return NewUserExperienceAnalyticsScoreHistoryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsSummarizeWorkFromAnywhereDevices provides operations to call the userExperienceAnalyticsSummarizeWorkFromAnywhereDevices method.
// returns a *UserExperienceAnalyticsSummarizeWorkFromAnywhereDevicesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsSummarizeWorkFromAnywhereDevices()(*UserExperienceAnalyticsSummarizeWorkFromAnywhereDevicesRequestBuilder) {
    return NewUserExperienceAnalyticsSummarizeWorkFromAnywhereDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric provides operations to manage the userExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetric()(*UserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricRequestBuilder) {
    return NewUserExperienceAnalyticsWorkFromAnywhereHardwareReadinessMetricRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsWorkFromAnywhereMetrics provides operations to manage the userExperienceAnalyticsWorkFromAnywhereMetrics property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsWorkFromAnywhereMetricsRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsWorkFromAnywhereMetrics()(*UserExperienceAnalyticsWorkFromAnywhereMetricsRequestBuilder) {
    return NewUserExperienceAnalyticsWorkFromAnywhereMetricsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UserExperienceAnalyticsWorkFromAnywhereModelPerformance provides operations to manage the userExperienceAnalyticsWorkFromAnywhereModelPerformance property of the microsoft.graph.deviceManagement entity.
// returns a *UserExperienceAnalyticsWorkFromAnywhereModelPerformanceRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) UserExperienceAnalyticsWorkFromAnywhereModelPerformance()(*UserExperienceAnalyticsWorkFromAnywhereModelPerformanceRequestBuilder) {
    return NewUserExperienceAnalyticsWorkFromAnywhereModelPerformanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// VerifyWindowsEnrollmentAutoDiscoveryWithDomainName provides operations to call the verifyWindowsEnrollmentAutoDiscovery method.
// returns a *VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) VerifyWindowsEnrollmentAutoDiscoveryWithDomainName(domainName *string)(*VerifyWindowsEnrollmentAutoDiscoveryWithDomainNameRequestBuilder) {
    return NewVerifyWindowsEnrollmentAutoDiscoveryWithDomainNameRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, domainName)
}
// VirtualEndpoint provides operations to manage the virtualEndpoint property of the microsoft.graph.deviceManagement entity.
// returns a *VirtualEndpointRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) VirtualEndpoint()(*VirtualEndpointRequestBuilder) {
    return NewVirtualEndpointRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsAutopilotDeviceIdentities provides operations to manage the windowsAutopilotDeviceIdentities property of the microsoft.graph.deviceManagement entity.
// returns a *WindowsAutopilotDeviceIdentitiesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) WindowsAutopilotDeviceIdentities()(*WindowsAutopilotDeviceIdentitiesRequestBuilder) {
    return NewWindowsAutopilotDeviceIdentitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsInformationProtectionAppLearningSummaries provides operations to manage the windowsInformationProtectionAppLearningSummaries property of the microsoft.graph.deviceManagement entity.
// returns a *WindowsInformationProtectionAppLearningSummariesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) WindowsInformationProtectionAppLearningSummaries()(*WindowsInformationProtectionAppLearningSummariesRequestBuilder) {
    return NewWindowsInformationProtectionAppLearningSummariesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsInformationProtectionNetworkLearningSummaries provides operations to manage the windowsInformationProtectionNetworkLearningSummaries property of the microsoft.graph.deviceManagement entity.
// returns a *WindowsInformationProtectionNetworkLearningSummariesRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) WindowsInformationProtectionNetworkLearningSummaries()(*WindowsInformationProtectionNetworkLearningSummariesRequestBuilder) {
    return NewWindowsInformationProtectionNetworkLearningSummariesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsMalwareInformation provides operations to manage the windowsMalwareInformation property of the microsoft.graph.deviceManagement entity.
// returns a *WindowsMalwareInformationRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) WindowsMalwareInformation()(*WindowsMalwareInformationRequestBuilder) {
    return NewWindowsMalwareInformationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *DeviceManagementRequestBuilder when successful
func (m *DeviceManagementRequestBuilder) WithUrl(rawUrl string)(*DeviceManagementRequestBuilder) {
    return NewDeviceManagementRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
