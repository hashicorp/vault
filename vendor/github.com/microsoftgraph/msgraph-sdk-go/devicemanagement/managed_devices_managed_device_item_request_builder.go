package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ManagedDevicesManagedDeviceItemRequestBuilder provides operations to manage the managedDevices property of the microsoft.graph.deviceManagement entity.
type ManagedDevicesManagedDeviceItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ManagedDevicesManagedDeviceItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesManagedDeviceItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ManagedDevicesManagedDeviceItemRequestBuilderGetQueryParameters the list of managed devices.
type ManagedDevicesManagedDeviceItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ManagedDevicesManagedDeviceItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesManagedDeviceItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ManagedDevicesManagedDeviceItemRequestBuilderGetQueryParameters
}
// ManagedDevicesManagedDeviceItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ManagedDevicesManagedDeviceItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BypassActivationLock provides operations to call the bypassActivationLock method.
// returns a *ManagedDevicesItemBypassActivationLockRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) BypassActivationLock()(*ManagedDevicesItemBypassActivationLockRequestBuilder) {
    return NewManagedDevicesItemBypassActivationLockRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CleanWindowsDevice provides operations to call the cleanWindowsDevice method.
// returns a *ManagedDevicesItemCleanWindowsDeviceRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) CleanWindowsDevice()(*ManagedDevicesItemCleanWindowsDeviceRequestBuilder) {
    return NewManagedDevicesItemCleanWindowsDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewManagedDevicesManagedDeviceItemRequestBuilderInternal instantiates a new ManagedDevicesManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDevicesManagedDeviceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDevicesManagedDeviceItemRequestBuilder) {
    m := &ManagedDevicesManagedDeviceItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/managedDevices/{managedDevice%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewManagedDevicesManagedDeviceItemRequestBuilder instantiates a new ManagedDevicesManagedDeviceItemRequestBuilder and sets the default values.
func NewManagedDevicesManagedDeviceItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ManagedDevicesManagedDeviceItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewManagedDevicesManagedDeviceItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete deletes a managedDevice.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-devices-manageddevice-delete?view=graph-rest-1.0
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DeleteUserFromSharedAppleDevice provides operations to call the deleteUserFromSharedAppleDevice method.
// returns a *ManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) DeleteUserFromSharedAppleDevice()(*ManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) {
    return NewManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCategory provides operations to manage the deviceCategory property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemDeviceCategoryRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) DeviceCategory()(*ManagedDevicesItemDeviceCategoryRequestBuilder) {
    return NewManagedDevicesItemDeviceCategoryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceCompliancePolicyStates provides operations to manage the deviceCompliancePolicyStates property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemDeviceCompliancePolicyStatesRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) DeviceCompliancePolicyStates()(*ManagedDevicesItemDeviceCompliancePolicyStatesRequestBuilder) {
    return NewManagedDevicesItemDeviceCompliancePolicyStatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceConfigurationStates provides operations to manage the deviceConfigurationStates property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemDeviceConfigurationStatesRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) DeviceConfigurationStates()(*ManagedDevicesItemDeviceConfigurationStatesRequestBuilder) {
    return NewManagedDevicesItemDeviceConfigurationStatesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DisableLostMode provides operations to call the disableLostMode method.
// returns a *ManagedDevicesItemDisableLostModeRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) DisableLostMode()(*ManagedDevicesItemDisableLostModeRequestBuilder) {
    return NewManagedDevicesItemDisableLostModeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the list of managed devices.
// returns a ManagedDeviceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable), nil
}
// LocateDevice provides operations to call the locateDevice method.
// returns a *ManagedDevicesItemLocateDeviceRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) LocateDevice()(*ManagedDevicesItemLocateDeviceRequestBuilder) {
    return NewManagedDevicesItemLocateDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LogCollectionRequests provides operations to manage the logCollectionRequests property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemLogCollectionRequestsRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) LogCollectionRequests()(*ManagedDevicesItemLogCollectionRequestsRequestBuilder) {
    return NewManagedDevicesItemLogCollectionRequestsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LogoutSharedAppleDeviceActiveUser provides operations to call the logoutSharedAppleDeviceActiveUser method.
// returns a *ManagedDevicesItemLogoutSharedAppleDeviceActiveUserRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) LogoutSharedAppleDeviceActiveUser()(*ManagedDevicesItemLogoutSharedAppleDeviceActiveUserRequestBuilder) {
    return NewManagedDevicesItemLogoutSharedAppleDeviceActiveUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the navigation property managedDevices in deviceManagement
// returns a ManagedDeviceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedDeviceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable), nil
}
// RebootNow provides operations to call the rebootNow method.
// returns a *ManagedDevicesItemRebootNowRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) RebootNow()(*ManagedDevicesItemRebootNowRequestBuilder) {
    return NewManagedDevicesItemRebootNowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RecoverPasscode provides operations to call the recoverPasscode method.
// returns a *ManagedDevicesItemRecoverPasscodeRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) RecoverPasscode()(*ManagedDevicesItemRecoverPasscodeRequestBuilder) {
    return NewManagedDevicesItemRecoverPasscodeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RemoteLock provides operations to call the remoteLock method.
// returns a *ManagedDevicesItemRemoteLockRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) RemoteLock()(*ManagedDevicesItemRemoteLockRequestBuilder) {
    return NewManagedDevicesItemRemoteLockRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RequestRemoteAssistance provides operations to call the requestRemoteAssistance method.
// returns a *ManagedDevicesItemRequestRemoteAssistanceRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) RequestRemoteAssistance()(*ManagedDevicesItemRequestRemoteAssistanceRequestBuilder) {
    return NewManagedDevicesItemRequestRemoteAssistanceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ResetPasscode provides operations to call the resetPasscode method.
// returns a *ManagedDevicesItemResetPasscodeRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) ResetPasscode()(*ManagedDevicesItemResetPasscodeRequestBuilder) {
    return NewManagedDevicesItemResetPasscodeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Retire provides operations to call the retire method.
// returns a *ManagedDevicesItemRetireRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Retire()(*ManagedDevicesItemRetireRequestBuilder) {
    return NewManagedDevicesItemRetireRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ShutDown provides operations to call the shutDown method.
// returns a *ManagedDevicesItemShutDownRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) ShutDown()(*ManagedDevicesItemShutDownRequestBuilder) {
    return NewManagedDevicesItemShutDownRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SyncDevice provides operations to call the syncDevice method.
// returns a *ManagedDevicesItemSyncDeviceRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) SyncDevice()(*ManagedDevicesItemSyncDeviceRequestBuilder) {
    return NewManagedDevicesItemSyncDeviceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation deletes a managedDevice.
// returns a *RequestInformation when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation the list of managed devices.
// returns a *RequestInformation when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property managedDevices in deviceManagement
// returns a *RequestInformation when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedDeviceable, requestConfiguration *ManagedDevicesManagedDeviceItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// UpdateWindowsDeviceAccount provides operations to call the updateWindowsDeviceAccount method.
// returns a *ManagedDevicesItemUpdateWindowsDeviceAccountRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) UpdateWindowsDeviceAccount()(*ManagedDevicesItemUpdateWindowsDeviceAccountRequestBuilder) {
    return NewManagedDevicesItemUpdateWindowsDeviceAccountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Users provides operations to manage the users property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemUsersRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Users()(*ManagedDevicesItemUsersRequestBuilder) {
    return NewManagedDevicesItemUsersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsDefenderScan provides operations to call the windowsDefenderScan method.
// returns a *ManagedDevicesItemWindowsDefenderScanRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) WindowsDefenderScan()(*ManagedDevicesItemWindowsDefenderScanRequestBuilder) {
    return NewManagedDevicesItemWindowsDefenderScanRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsDefenderUpdateSignatures provides operations to call the windowsDefenderUpdateSignatures method.
// returns a *ManagedDevicesItemWindowsDefenderUpdateSignaturesRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) WindowsDefenderUpdateSignatures()(*ManagedDevicesItemWindowsDefenderUpdateSignaturesRequestBuilder) {
    return NewManagedDevicesItemWindowsDefenderUpdateSignaturesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WindowsProtectionState provides operations to manage the windowsProtectionState property of the microsoft.graph.managedDevice entity.
// returns a *ManagedDevicesItemWindowsProtectionStateRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) WindowsProtectionState()(*ManagedDevicesItemWindowsProtectionStateRequestBuilder) {
    return NewManagedDevicesItemWindowsProtectionStateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Wipe provides operations to call the wipe method.
// returns a *ManagedDevicesItemWipeRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) Wipe()(*ManagedDevicesItemWipeRequestBuilder) {
    return NewManagedDevicesItemWipeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ManagedDevicesManagedDeviceItemRequestBuilder when successful
func (m *ManagedDevicesManagedDeviceItemRequestBuilder) WithUrl(rawUrl string)(*ManagedDevicesManagedDeviceItemRequestBuilder) {
    return NewManagedDevicesManagedDeviceItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
