package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// UserItemRequestBuilder provides operations to manage the collection of user entities.
type UserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// UserItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// UserItemRequestBuilderGetQueryParameters read properties and relationships of the user object.
type UserItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// UserItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserItemRequestBuilderGetQueryParameters
}
// UserItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Activities provides operations to manage the activities property of the microsoft.graph.user entity.
// returns a *ItemActivitiesRequestBuilder when successful
func (m *UserItemRequestBuilder) Activities()(*ItemActivitiesRequestBuilder) {
    return NewItemActivitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AgreementAcceptances provides operations to manage the agreementAcceptances property of the microsoft.graph.user entity.
// returns a *ItemAgreementAcceptancesRequestBuilder when successful
func (m *UserItemRequestBuilder) AgreementAcceptances()(*ItemAgreementAcceptancesRequestBuilder) {
    return NewItemAgreementAcceptancesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AppRoleAssignments provides operations to manage the appRoleAssignments property of the microsoft.graph.user entity.
// returns a *ItemAppRoleAssignmentsRequestBuilder when successful
func (m *UserItemRequestBuilder) AppRoleAssignments()(*ItemAppRoleAssignmentsRequestBuilder) {
    return NewItemAppRoleAssignmentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AssignLicense provides operations to call the assignLicense method.
// returns a *ItemAssignLicenseRequestBuilder when successful
func (m *UserItemRequestBuilder) AssignLicense()(*ItemAssignLicenseRequestBuilder) {
    return NewItemAssignLicenseRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Authentication provides operations to manage the authentication property of the microsoft.graph.user entity.
// returns a *ItemAuthenticationRequestBuilder when successful
func (m *UserItemRequestBuilder) Authentication()(*ItemAuthenticationRequestBuilder) {
    return NewItemAuthenticationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.user entity.
// returns a *ItemCalendarRequestBuilder when successful
func (m *UserItemRequestBuilder) Calendar()(*ItemCalendarRequestBuilder) {
    return NewItemCalendarRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CalendarGroups provides operations to manage the calendarGroups property of the microsoft.graph.user entity.
// returns a *ItemCalendarGroupsRequestBuilder when successful
func (m *UserItemRequestBuilder) CalendarGroups()(*ItemCalendarGroupsRequestBuilder) {
    return NewItemCalendarGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Calendars provides operations to manage the calendars property of the microsoft.graph.user entity.
// returns a *ItemCalendarsRequestBuilder when successful
func (m *UserItemRequestBuilder) Calendars()(*ItemCalendarsRequestBuilder) {
    return NewItemCalendarsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CalendarView provides operations to manage the calendarView property of the microsoft.graph.user entity.
// returns a *ItemCalendarViewRequestBuilder when successful
func (m *UserItemRequestBuilder) CalendarView()(*ItemCalendarViewRequestBuilder) {
    return NewItemCalendarViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ChangePassword provides operations to call the changePassword method.
// returns a *ItemChangePasswordRequestBuilder when successful
func (m *UserItemRequestBuilder) ChangePassword()(*ItemChangePasswordRequestBuilder) {
    return NewItemChangePasswordRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Chats provides operations to manage the chats property of the microsoft.graph.user entity.
// returns a *ItemChatsRequestBuilder when successful
func (m *UserItemRequestBuilder) Chats()(*ItemChatsRequestBuilder) {
    return NewItemChatsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
// returns a *ItemCheckMemberGroupsRequestBuilder when successful
func (m *UserItemRequestBuilder) CheckMemberGroups()(*ItemCheckMemberGroupsRequestBuilder) {
    return NewItemCheckMemberGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
// returns a *ItemCheckMemberObjectsRequestBuilder when successful
func (m *UserItemRequestBuilder) CheckMemberObjects()(*ItemCheckMemberObjectsRequestBuilder) {
    return NewItemCheckMemberObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CloudClipboard provides operations to manage the cloudClipboard property of the microsoft.graph.user entity.
// returns a *ItemCloudClipboardRequestBuilder when successful
func (m *UserItemRequestBuilder) CloudClipboard()(*ItemCloudClipboardRequestBuilder) {
    return NewItemCloudClipboardRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewUserItemRequestBuilderInternal instantiates a new UserItemRequestBuilder and sets the default values.
func NewUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserItemRequestBuilder) {
    m := &UserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewUserItemRequestBuilder instantiates a new UserItemRequestBuilder and sets the default values.
func NewUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// ContactFolders provides operations to manage the contactFolders property of the microsoft.graph.user entity.
// returns a *ItemContactFoldersRequestBuilder when successful
func (m *UserItemRequestBuilder) ContactFolders()(*ItemContactFoldersRequestBuilder) {
    return NewItemContactFoldersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Contacts provides operations to manage the contacts property of the microsoft.graph.user entity.
// returns a *ItemContactsRequestBuilder when successful
func (m *UserItemRequestBuilder) Contacts()(*ItemContactsRequestBuilder) {
    return NewItemContactsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// CreatedObjects provides operations to manage the createdObjects property of the microsoft.graph.user entity.
// returns a *ItemCreatedObjectsRequestBuilder when successful
func (m *UserItemRequestBuilder) CreatedObjects()(*ItemCreatedObjectsRequestBuilder) {
    return NewItemCreatedObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Delete deletes a user.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-mam-user-delete?view=graph-rest-1.0
func (m *UserItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *UserItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// DeviceManagementTroubleshootingEvents provides operations to manage the deviceManagementTroubleshootingEvents property of the microsoft.graph.user entity.
// returns a *ItemDeviceManagementTroubleshootingEventsRequestBuilder when successful
func (m *UserItemRequestBuilder) DeviceManagementTroubleshootingEvents()(*ItemDeviceManagementTroubleshootingEventsRequestBuilder) {
    return NewItemDeviceManagementTroubleshootingEventsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DirectReports provides operations to manage the directReports property of the microsoft.graph.user entity.
// returns a *ItemDirectReportsRequestBuilder when successful
func (m *UserItemRequestBuilder) DirectReports()(*ItemDirectReportsRequestBuilder) {
    return NewItemDirectReportsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drive provides operations to manage the drive property of the microsoft.graph.user entity.
// returns a *ItemDriveRequestBuilder when successful
func (m *UserItemRequestBuilder) Drive()(*ItemDriveRequestBuilder) {
    return NewItemDriveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Drives provides operations to manage the drives property of the microsoft.graph.user entity.
// returns a *ItemDrivesRequestBuilder when successful
func (m *UserItemRequestBuilder) Drives()(*ItemDrivesRequestBuilder) {
    return NewItemDrivesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EmployeeExperience provides operations to manage the employeeExperience property of the microsoft.graph.user entity.
// returns a *ItemEmployeeExperienceRequestBuilder when successful
func (m *UserItemRequestBuilder) EmployeeExperience()(*ItemEmployeeExperienceRequestBuilder) {
    return NewItemEmployeeExperienceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Events provides operations to manage the events property of the microsoft.graph.user entity.
// returns a *ItemEventsRequestBuilder when successful
func (m *UserItemRequestBuilder) Events()(*ItemEventsRequestBuilder) {
    return NewItemEventsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExportDeviceAndAppManagementData provides operations to call the exportDeviceAndAppManagementData method.
// returns a *ItemExportDeviceAndAppManagementDataRequestBuilder when successful
func (m *UserItemRequestBuilder) ExportDeviceAndAppManagementData()(*ItemExportDeviceAndAppManagementDataRequestBuilder) {
    return NewItemExportDeviceAndAppManagementDataRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ExportDeviceAndAppManagementDataWithSkipWithTop provides operations to call the exportDeviceAndAppManagementData method.
// returns a *ItemExportDeviceAndAppManagementDataWithSkipWithTopRequestBuilder when successful
func (m *UserItemRequestBuilder) ExportDeviceAndAppManagementDataWithSkipWithTop(skip *int32, top *int32)(*ItemExportDeviceAndAppManagementDataWithSkipWithTopRequestBuilder) {
    return NewItemExportDeviceAndAppManagementDataWithSkipWithTopRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, skip, top)
}
// ExportPersonalData provides operations to call the exportPersonalData method.
// returns a *ItemExportPersonalDataRequestBuilder when successful
func (m *UserItemRequestBuilder) ExportPersonalData()(*ItemExportPersonalDataRequestBuilder) {
    return NewItemExportPersonalDataRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.user entity.
// returns a *ItemExtensionsRequestBuilder when successful
func (m *UserItemRequestBuilder) Extensions()(*ItemExtensionsRequestBuilder) {
    return NewItemExtensionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FindMeetingTimes provides operations to call the findMeetingTimes method.
// returns a *ItemFindMeetingTimesRequestBuilder when successful
func (m *UserItemRequestBuilder) FindMeetingTimes()(*ItemFindMeetingTimesRequestBuilder) {
    return NewItemFindMeetingTimesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FollowedSites provides operations to manage the followedSites property of the microsoft.graph.user entity.
// returns a *ItemFollowedSitesRequestBuilder when successful
func (m *UserItemRequestBuilder) FollowedSites()(*ItemFollowedSitesRequestBuilder) {
    return NewItemFollowedSitesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get read properties and relationships of the user object.
// returns a Userable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-onboarding-user-get?view=graph-rest-1.0
func (m *UserItemRequestBuilder) Get(ctx context.Context, requestConfiguration *UserItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable), nil
}
// GetMailTips provides operations to call the getMailTips method.
// returns a *ItemGetMailTipsRequestBuilder when successful
func (m *UserItemRequestBuilder) GetMailTips()(*ItemGetMailTipsRequestBuilder) {
    return NewItemGetMailTipsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetManagedAppDiagnosticStatuses provides operations to call the getManagedAppDiagnosticStatuses method.
// returns a *ItemGetManagedAppDiagnosticStatusesRequestBuilder when successful
func (m *UserItemRequestBuilder) GetManagedAppDiagnosticStatuses()(*ItemGetManagedAppDiagnosticStatusesRequestBuilder) {
    return NewItemGetManagedAppDiagnosticStatusesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetManagedAppPolicies provides operations to call the getManagedAppPolicies method.
// returns a *ItemGetManagedAppPoliciesRequestBuilder when successful
func (m *UserItemRequestBuilder) GetManagedAppPolicies()(*ItemGetManagedAppPoliciesRequestBuilder) {
    return NewItemGetManagedAppPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetManagedDevicesWithAppFailures provides operations to call the getManagedDevicesWithAppFailures method.
// returns a *ItemGetManagedDevicesWithAppFailuresRequestBuilder when successful
func (m *UserItemRequestBuilder) GetManagedDevicesWithAppFailures()(*ItemGetManagedDevicesWithAppFailuresRequestBuilder) {
    return NewItemGetManagedDevicesWithAppFailuresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetMemberGroups provides operations to call the getMemberGroups method.
// returns a *ItemGetMemberGroupsRequestBuilder when successful
func (m *UserItemRequestBuilder) GetMemberGroups()(*ItemGetMemberGroupsRequestBuilder) {
    return NewItemGetMemberGroupsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetMemberObjects provides operations to call the getMemberObjects method.
// returns a *ItemGetMemberObjectsRequestBuilder when successful
func (m *UserItemRequestBuilder) GetMemberObjects()(*ItemGetMemberObjectsRequestBuilder) {
    return NewItemGetMemberObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// InferenceClassification provides operations to manage the inferenceClassification property of the microsoft.graph.user entity.
// returns a *ItemInferenceClassificationRequestBuilder when successful
func (m *UserItemRequestBuilder) InferenceClassification()(*ItemInferenceClassificationRequestBuilder) {
    return NewItemInferenceClassificationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Insights provides operations to manage the insights property of the microsoft.graph.user entity.
// returns a *ItemInsightsRequestBuilder when successful
func (m *UserItemRequestBuilder) Insights()(*ItemInsightsRequestBuilder) {
    return NewItemInsightsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// JoinedTeams provides operations to manage the joinedTeams property of the microsoft.graph.user entity.
// returns a *ItemJoinedTeamsRequestBuilder when successful
func (m *UserItemRequestBuilder) JoinedTeams()(*ItemJoinedTeamsRequestBuilder) {
    return NewItemJoinedTeamsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LicenseDetails provides operations to manage the licenseDetails property of the microsoft.graph.user entity.
// returns a *ItemLicenseDetailsRequestBuilder when successful
func (m *UserItemRequestBuilder) LicenseDetails()(*ItemLicenseDetailsRequestBuilder) {
    return NewItemLicenseDetailsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MailboxSettings the mailboxSettings property
// returns a *ItemMailboxSettingsRequestBuilder when successful
func (m *UserItemRequestBuilder) MailboxSettings()(*ItemMailboxSettingsRequestBuilder) {
    return NewItemMailboxSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MailFolders provides operations to manage the mailFolders property of the microsoft.graph.user entity.
// returns a *ItemMailFoldersRequestBuilder when successful
func (m *UserItemRequestBuilder) MailFolders()(*ItemMailFoldersRequestBuilder) {
    return NewItemMailFoldersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedAppRegistrations provides operations to manage the managedAppRegistrations property of the microsoft.graph.user entity.
// returns a *ItemManagedAppRegistrationsRequestBuilder when successful
func (m *UserItemRequestBuilder) ManagedAppRegistrations()(*ItemManagedAppRegistrationsRequestBuilder) {
    return NewItemManagedAppRegistrationsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedDevices provides operations to manage the managedDevices property of the microsoft.graph.user entity.
// returns a *ItemManagedDevicesRequestBuilder when successful
func (m *UserItemRequestBuilder) ManagedDevices()(*ItemManagedDevicesRequestBuilder) {
    return NewItemManagedDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Manager provides operations to manage the manager property of the microsoft.graph.user entity.
// returns a *ItemManagerRequestBuilder when successful
func (m *UserItemRequestBuilder) Manager()(*ItemManagerRequestBuilder) {
    return NewItemManagerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MemberOf provides operations to manage the memberOf property of the microsoft.graph.user entity.
// returns a *ItemMemberOfRequestBuilder when successful
func (m *UserItemRequestBuilder) MemberOf()(*ItemMemberOfRequestBuilder) {
    return NewItemMemberOfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Messages provides operations to manage the messages property of the microsoft.graph.user entity.
// returns a *ItemMessagesRequestBuilder when successful
func (m *UserItemRequestBuilder) Messages()(*ItemMessagesRequestBuilder) {
    return NewItemMessagesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Oauth2PermissionGrants provides operations to manage the oauth2PermissionGrants property of the microsoft.graph.user entity.
// returns a *ItemOauth2PermissionGrantsRequestBuilder when successful
func (m *UserItemRequestBuilder) Oauth2PermissionGrants()(*ItemOauth2PermissionGrantsRequestBuilder) {
    return NewItemOauth2PermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.user entity.
// returns a *ItemOnenoteRequestBuilder when successful
func (m *UserItemRequestBuilder) Onenote()(*ItemOnenoteRequestBuilder) {
    return NewItemOnenoteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OnlineMeetings provides operations to manage the onlineMeetings property of the microsoft.graph.user entity.
// returns a *ItemOnlineMeetingsRequestBuilder when successful
func (m *UserItemRequestBuilder) OnlineMeetings()(*ItemOnlineMeetingsRequestBuilder) {
    return NewItemOnlineMeetingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Outlook provides operations to manage the outlook property of the microsoft.graph.user entity.
// returns a *ItemOutlookRequestBuilder when successful
func (m *UserItemRequestBuilder) Outlook()(*ItemOutlookRequestBuilder) {
    return NewItemOutlookRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OwnedDevices provides operations to manage the ownedDevices property of the microsoft.graph.user entity.
// returns a *ItemOwnedDevicesRequestBuilder when successful
func (m *UserItemRequestBuilder) OwnedDevices()(*ItemOwnedDevicesRequestBuilder) {
    return NewItemOwnedDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OwnedObjects provides operations to manage the ownedObjects property of the microsoft.graph.user entity.
// returns a *ItemOwnedObjectsRequestBuilder when successful
func (m *UserItemRequestBuilder) OwnedObjects()(*ItemOwnedObjectsRequestBuilder) {
    return NewItemOwnedObjectsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the properties of a user object.
// returns a Userable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-mam-user-update?view=graph-rest-1.0
func (m *UserItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, requestConfiguration *UserItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable), nil
}
// People provides operations to manage the people property of the microsoft.graph.user entity.
// returns a *ItemPeopleRequestBuilder when successful
func (m *UserItemRequestBuilder) People()(*ItemPeopleRequestBuilder) {
    return NewItemPeopleRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// PermissionGrants provides operations to manage the permissionGrants property of the microsoft.graph.user entity.
// returns a *ItemPermissionGrantsRequestBuilder when successful
func (m *UserItemRequestBuilder) PermissionGrants()(*ItemPermissionGrantsRequestBuilder) {
    return NewItemPermissionGrantsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Photo provides operations to manage the photo property of the microsoft.graph.user entity.
// returns a *ItemPhotoRequestBuilder when successful
func (m *UserItemRequestBuilder) Photo()(*ItemPhotoRequestBuilder) {
    return NewItemPhotoRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Photos provides operations to manage the photos property of the microsoft.graph.user entity.
// returns a *ItemPhotosRequestBuilder when successful
func (m *UserItemRequestBuilder) Photos()(*ItemPhotosRequestBuilder) {
    return NewItemPhotosRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Planner provides operations to manage the planner property of the microsoft.graph.user entity.
// returns a *ItemPlannerRequestBuilder when successful
func (m *UserItemRequestBuilder) Planner()(*ItemPlannerRequestBuilder) {
    return NewItemPlannerRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Presence provides operations to manage the presence property of the microsoft.graph.user entity.
// returns a *ItemPresenceRequestBuilder when successful
func (m *UserItemRequestBuilder) Presence()(*ItemPresenceRequestBuilder) {
    return NewItemPresenceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RegisteredDevices provides operations to manage the registeredDevices property of the microsoft.graph.user entity.
// returns a *ItemRegisteredDevicesRequestBuilder when successful
func (m *UserItemRequestBuilder) RegisteredDevices()(*ItemRegisteredDevicesRequestBuilder) {
    return NewItemRegisteredDevicesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ReminderViewWithStartDateTimeWithEndDateTime provides operations to call the reminderView method.
// returns a *ItemReminderViewWithStartDateTimeWithEndDateTimeRequestBuilder when successful
func (m *UserItemRequestBuilder) ReminderViewWithStartDateTimeWithEndDateTime(endDateTime *string, startDateTime *string)(*ItemReminderViewWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return NewItemReminderViewWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, startDateTime)
}
// RemoveAllDevicesFromManagement provides operations to call the removeAllDevicesFromManagement method.
// returns a *ItemRemoveAllDevicesFromManagementRequestBuilder when successful
func (m *UserItemRequestBuilder) RemoveAllDevicesFromManagement()(*ItemRemoveAllDevicesFromManagementRequestBuilder) {
    return NewItemRemoveAllDevicesFromManagementRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ReprocessLicenseAssignment provides operations to call the reprocessLicenseAssignment method.
// returns a *ItemReprocessLicenseAssignmentRequestBuilder when successful
func (m *UserItemRequestBuilder) ReprocessLicenseAssignment()(*ItemReprocessLicenseAssignmentRequestBuilder) {
    return NewItemReprocessLicenseAssignmentRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Restore provides operations to call the restore method.
// returns a *ItemRestoreRequestBuilder when successful
func (m *UserItemRequestBuilder) Restore()(*ItemRestoreRequestBuilder) {
    return NewItemRestoreRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RetryServiceProvisioning provides operations to call the retryServiceProvisioning method.
// returns a *ItemRetryServiceProvisioningRequestBuilder when successful
func (m *UserItemRequestBuilder) RetryServiceProvisioning()(*ItemRetryServiceProvisioningRequestBuilder) {
    return NewItemRetryServiceProvisioningRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RevokeSignInSessions provides operations to call the revokeSignInSessions method.
// returns a *ItemRevokeSignInSessionsRequestBuilder when successful
func (m *UserItemRequestBuilder) RevokeSignInSessions()(*ItemRevokeSignInSessionsRequestBuilder) {
    return NewItemRevokeSignInSessionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ScopedRoleMemberOf provides operations to manage the scopedRoleMemberOf property of the microsoft.graph.user entity.
// returns a *ItemScopedRoleMemberOfRequestBuilder when successful
func (m *UserItemRequestBuilder) ScopedRoleMemberOf()(*ItemScopedRoleMemberOfRequestBuilder) {
    return NewItemScopedRoleMemberOfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SendMail provides operations to call the sendMail method.
// returns a *ItemSendMailRequestBuilder when successful
func (m *UserItemRequestBuilder) SendMail()(*ItemSendMailRequestBuilder) {
    return NewItemSendMailRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceProvisioningErrors the serviceProvisioningErrors property
// returns a *ItemServiceProvisioningErrorsRequestBuilder when successful
func (m *UserItemRequestBuilder) ServiceProvisioningErrors()(*ItemServiceProvisioningErrorsRequestBuilder) {
    return NewItemServiceProvisioningErrorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Settings provides operations to manage the settings property of the microsoft.graph.user entity.
// returns a *ItemSettingsRequestBuilder when successful
func (m *UserItemRequestBuilder) Settings()(*ItemSettingsRequestBuilder) {
    return NewItemSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Solutions provides operations to manage the solutions property of the microsoft.graph.user entity.
// returns a *ItemSolutionsRequestBuilder when successful
func (m *UserItemRequestBuilder) Solutions()(*ItemSolutionsRequestBuilder) {
    return NewItemSolutionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sponsors provides operations to manage the sponsors property of the microsoft.graph.user entity.
// returns a *ItemSponsorsRequestBuilder when successful
func (m *UserItemRequestBuilder) Sponsors()(*ItemSponsorsRequestBuilder) {
    return NewItemSponsorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Teamwork provides operations to manage the teamwork property of the microsoft.graph.user entity.
// returns a *ItemTeamworkRequestBuilder when successful
func (m *UserItemRequestBuilder) Teamwork()(*ItemTeamworkRequestBuilder) {
    return NewItemTeamworkRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation deletes a user.
// returns a *RequestInformation when successful
func (m *UserItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *UserItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Todo provides operations to manage the todo property of the microsoft.graph.user entity.
// returns a *ItemTodoRequestBuilder when successful
func (m *UserItemRequestBuilder) Todo()(*ItemTodoRequestBuilder) {
    return NewItemTodoRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation read properties and relationships of the user object.
// returns a *RequestInformation when successful
func (m *UserItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *UserItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a user object.
// returns a *RequestInformation when successful
func (m *UserItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, requestConfiguration *UserItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// TransitiveMemberOf provides operations to manage the transitiveMemberOf property of the microsoft.graph.user entity.
// returns a *ItemTransitiveMemberOfRequestBuilder when successful
func (m *UserItemRequestBuilder) TransitiveMemberOf()(*ItemTransitiveMemberOfRequestBuilder) {
    return NewItemTransitiveMemberOfRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TranslateExchangeIds provides operations to call the translateExchangeIds method.
// returns a *ItemTranslateExchangeIdsRequestBuilder when successful
func (m *UserItemRequestBuilder) TranslateExchangeIds()(*ItemTranslateExchangeIdsRequestBuilder) {
    return NewItemTranslateExchangeIdsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WipeManagedAppRegistrationsByDeviceTag provides operations to call the wipeManagedAppRegistrationsByDeviceTag method.
// returns a *ItemWipeManagedAppRegistrationsByDeviceTagRequestBuilder when successful
func (m *UserItemRequestBuilder) WipeManagedAppRegistrationsByDeviceTag()(*ItemWipeManagedAppRegistrationsByDeviceTagRequestBuilder) {
    return NewItemWipeManagedAppRegistrationsByDeviceTagRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *UserItemRequestBuilder when successful
func (m *UserItemRequestBuilder) WithUrl(rawUrl string)(*UserItemRequestBuilder) {
    return NewUserItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
