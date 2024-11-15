package reports

import (
    "context"
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ReportsRequestBuilder provides operations to manage the reportRoot singleton.
type ReportsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ReportsRequestBuilderGetQueryParameters read properties and relationships of the reportRoot object.
type ReportsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ReportsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ReportsRequestBuilderGetQueryParameters
}
// ReportsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AuthenticationMethods provides operations to manage the authenticationMethods property of the microsoft.graph.reportRoot entity.
// returns a *AuthenticationMethodsRequestBuilder when successful
func (m *ReportsRequestBuilder) AuthenticationMethods()(*AuthenticationMethodsRequestBuilder) {
    return NewAuthenticationMethodsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewReportsRequestBuilderInternal instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    m := &ReportsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewReportsRequestBuilder instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewReportsRequestBuilderInternal(urlParams, requestAdapter)
}
// DailyPrintUsageByPrinter provides operations to manage the dailyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
// returns a *DailyPrintUsageByPrinterRequestBuilder when successful
func (m *ReportsRequestBuilder) DailyPrintUsageByPrinter()(*DailyPrintUsageByPrinterRequestBuilder) {
    return NewDailyPrintUsageByPrinterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DailyPrintUsageByUser provides operations to manage the dailyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
// returns a *DailyPrintUsageByUserRequestBuilder when successful
func (m *ReportsRequestBuilder) DailyPrintUsageByUser()(*DailyPrintUsageByUserRequestBuilder) {
    return NewDailyPrintUsageByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceConfigurationDeviceActivity provides operations to call the deviceConfigurationDeviceActivity method.
// returns a *DeviceConfigurationDeviceActivityRequestBuilder when successful
func (m *ReportsRequestBuilder) DeviceConfigurationDeviceActivity()(*DeviceConfigurationDeviceActivityRequestBuilder) {
    return NewDeviceConfigurationDeviceActivityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceConfigurationUserActivity provides operations to call the deviceConfigurationUserActivity method.
// returns a *DeviceConfigurationUserActivityRequestBuilder when successful
func (m *ReportsRequestBuilder) DeviceConfigurationUserActivity()(*DeviceConfigurationUserActivityRequestBuilder) {
    return NewDeviceConfigurationUserActivityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get read properties and relationships of the reportRoot object.
// returns a ReportRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-deviceconfig-reportroot-get?view=graph-rest-1.0
func (m *ReportsRequestBuilder) Get(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateReportRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable), nil
}
// GetEmailActivityCountsWithPeriod provides operations to call the getEmailActivityCounts method.
// returns a *GetEmailActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailActivityCountsWithPeriod(period *string)(*GetEmailActivityCountsWithPeriodRequestBuilder) {
    return NewGetEmailActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailActivityUserCountsWithPeriod provides operations to call the getEmailActivityUserCounts method.
// returns a *GetEmailActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailActivityUserCountsWithPeriod(period *string)(*GetEmailActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetEmailActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailActivityUserDetailWithDate provides operations to call the getEmailActivityUserDetail method.
// returns a *GetEmailActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetEmailActivityUserDetailWithDateRequestBuilder) {
    return NewGetEmailActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetEmailActivityUserDetailWithPeriod provides operations to call the getEmailActivityUserDetail method.
// returns a *GetEmailActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailActivityUserDetailWithPeriod(period *string)(*GetEmailActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetEmailActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailAppUsageAppsUserCountsWithPeriod provides operations to call the getEmailAppUsageAppsUserCounts method.
// returns a *GetEmailAppUsageAppsUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailAppUsageAppsUserCountsWithPeriod(period *string)(*GetEmailAppUsageAppsUserCountsWithPeriodRequestBuilder) {
    return NewGetEmailAppUsageAppsUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailAppUsageUserCountsWithPeriod provides operations to call the getEmailAppUsageUserCounts method.
// returns a *GetEmailAppUsageUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailAppUsageUserCountsWithPeriod(period *string)(*GetEmailAppUsageUserCountsWithPeriodRequestBuilder) {
    return NewGetEmailAppUsageUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailAppUsageUserDetailWithDate provides operations to call the getEmailAppUsageUserDetail method.
// returns a *GetEmailAppUsageUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailAppUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetEmailAppUsageUserDetailWithDateRequestBuilder) {
    return NewGetEmailAppUsageUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetEmailAppUsageUserDetailWithPeriod provides operations to call the getEmailAppUsageUserDetail method.
// returns a *GetEmailAppUsageUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailAppUsageUserDetailWithPeriod(period *string)(*GetEmailAppUsageUserDetailWithPeriodRequestBuilder) {
    return NewGetEmailAppUsageUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetEmailAppUsageVersionsUserCountsWithPeriod provides operations to call the getEmailAppUsageVersionsUserCounts method.
// returns a *GetEmailAppUsageVersionsUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetEmailAppUsageVersionsUserCountsWithPeriod(period *string)(*GetEmailAppUsageVersionsUserCountsWithPeriodRequestBuilder) {
    return NewGetEmailAppUsageVersionsUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTime provides operations to call the getGroupArchivedPrintJobs method.
// returns a *GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTimeRequestBuilder when successful
func (m *ReportsRequestBuilder) GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, groupId *string, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)(*GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return NewGetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, groupId, startDateTime)
}
// GetM365AppPlatformUserCountsWithPeriod provides operations to call the getM365AppPlatformUserCounts method.
// returns a *GetM365AppPlatformUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetM365AppPlatformUserCountsWithPeriod(period *string)(*GetM365AppPlatformUserCountsWithPeriodRequestBuilder) {
    return NewGetM365AppPlatformUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetM365AppUserCountsWithPeriod provides operations to call the getM365AppUserCounts method.
// returns a *GetM365AppUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetM365AppUserCountsWithPeriod(period *string)(*GetM365AppUserCountsWithPeriodRequestBuilder) {
    return NewGetM365AppUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetM365AppUserDetailWithDate provides operations to call the getM365AppUserDetail method.
// returns a *GetM365AppUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetM365AppUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetM365AppUserDetailWithDateRequestBuilder) {
    return NewGetM365AppUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetM365AppUserDetailWithPeriod provides operations to call the getM365AppUserDetail method.
// returns a *GetM365AppUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetM365AppUserDetailWithPeriod(period *string)(*GetM365AppUserDetailWithPeriodRequestBuilder) {
    return NewGetM365AppUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetMailboxUsageDetailWithPeriod provides operations to call the getMailboxUsageDetail method.
// returns a *GetMailboxUsageDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetMailboxUsageDetailWithPeriod(period *string)(*GetMailboxUsageDetailWithPeriodRequestBuilder) {
    return NewGetMailboxUsageDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetMailboxUsageMailboxCountsWithPeriod provides operations to call the getMailboxUsageMailboxCounts method.
// returns a *GetMailboxUsageMailboxCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetMailboxUsageMailboxCountsWithPeriod(period *string)(*GetMailboxUsageMailboxCountsWithPeriodRequestBuilder) {
    return NewGetMailboxUsageMailboxCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetMailboxUsageQuotaStatusMailboxCountsWithPeriod provides operations to call the getMailboxUsageQuotaStatusMailboxCounts method.
// returns a *GetMailboxUsageQuotaStatusMailboxCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetMailboxUsageQuotaStatusMailboxCountsWithPeriod(period *string)(*GetMailboxUsageQuotaStatusMailboxCountsWithPeriodRequestBuilder) {
    return NewGetMailboxUsageQuotaStatusMailboxCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetMailboxUsageStorageWithPeriod provides operations to call the getMailboxUsageStorage method.
// returns a *GetMailboxUsageStorageWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetMailboxUsageStorageWithPeriod(period *string)(*GetMailboxUsageStorageWithPeriodRequestBuilder) {
    return NewGetMailboxUsageStorageWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365ActivationCounts provides operations to call the getOffice365ActivationCounts method.
// returns a *GetOffice365ActivationCountsRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActivationCounts()(*GetOffice365ActivationCountsRequestBuilder) {
    return NewGetOffice365ActivationCountsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetOffice365ActivationsUserCounts provides operations to call the getOffice365ActivationsUserCounts method.
// returns a *GetOffice365ActivationsUserCountsRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActivationsUserCounts()(*GetOffice365ActivationsUserCountsRequestBuilder) {
    return NewGetOffice365ActivationsUserCountsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetOffice365ActivationsUserDetail provides operations to call the getOffice365ActivationsUserDetail method.
// returns a *GetOffice365ActivationsUserDetailRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActivationsUserDetail()(*GetOffice365ActivationsUserDetailRequestBuilder) {
    return NewGetOffice365ActivationsUserDetailRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GetOffice365ActiveUserCountsWithPeriod provides operations to call the getOffice365ActiveUserCounts method.
// returns a *GetOffice365ActiveUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActiveUserCountsWithPeriod(period *string)(*GetOffice365ActiveUserCountsWithPeriodRequestBuilder) {
    return NewGetOffice365ActiveUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365ActiveUserDetailWithDate provides operations to call the getOffice365ActiveUserDetail method.
// returns a *GetOffice365ActiveUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActiveUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetOffice365ActiveUserDetailWithDateRequestBuilder) {
    return NewGetOffice365ActiveUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetOffice365ActiveUserDetailWithPeriod provides operations to call the getOffice365ActiveUserDetail method.
// returns a *GetOffice365ActiveUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ActiveUserDetailWithPeriod(period *string)(*GetOffice365ActiveUserDetailWithPeriodRequestBuilder) {
    return NewGetOffice365ActiveUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365GroupsActivityCountsWithPeriod provides operations to call the getOffice365GroupsActivityCounts method.
// returns a *GetOffice365GroupsActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityCountsWithPeriod(period *string)(*GetOffice365GroupsActivityCountsWithPeriodRequestBuilder) {
    return NewGetOffice365GroupsActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365GroupsActivityDetailWithDate provides operations to call the getOffice365GroupsActivityDetail method.
// returns a *GetOffice365GroupsActivityDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetOffice365GroupsActivityDetailWithDateRequestBuilder) {
    return NewGetOffice365GroupsActivityDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetOffice365GroupsActivityDetailWithPeriod provides operations to call the getOffice365GroupsActivityDetail method.
// returns a *GetOffice365GroupsActivityDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityDetailWithPeriod(period *string)(*GetOffice365GroupsActivityDetailWithPeriodRequestBuilder) {
    return NewGetOffice365GroupsActivityDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365GroupsActivityFileCountsWithPeriod provides operations to call the getOffice365GroupsActivityFileCounts method.
// returns a *GetOffice365GroupsActivityFileCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityFileCountsWithPeriod(period *string)(*GetOffice365GroupsActivityFileCountsWithPeriodRequestBuilder) {
    return NewGetOffice365GroupsActivityFileCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365GroupsActivityGroupCountsWithPeriod provides operations to call the getOffice365GroupsActivityGroupCounts method.
// returns a *GetOffice365GroupsActivityGroupCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityGroupCountsWithPeriod(period *string)(*GetOffice365GroupsActivityGroupCountsWithPeriodRequestBuilder) {
    return NewGetOffice365GroupsActivityGroupCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365GroupsActivityStorageWithPeriod provides operations to call the getOffice365GroupsActivityStorage method.
// returns a *GetOffice365GroupsActivityStorageWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityStorageWithPeriod(period *string)(*GetOffice365GroupsActivityStorageWithPeriodRequestBuilder) {
    return NewGetOffice365GroupsActivityStorageWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOffice365ServicesUserCountsWithPeriod provides operations to call the getOffice365ServicesUserCounts method.
// returns a *GetOffice365ServicesUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOffice365ServicesUserCountsWithPeriod(period *string)(*GetOffice365ServicesUserCountsWithPeriodRequestBuilder) {
    return NewGetOffice365ServicesUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveActivityFileCountsWithPeriod provides operations to call the getOneDriveActivityFileCounts method.
// returns a *GetOneDriveActivityFileCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveActivityFileCountsWithPeriod(period *string)(*GetOneDriveActivityFileCountsWithPeriodRequestBuilder) {
    return NewGetOneDriveActivityFileCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveActivityUserCountsWithPeriod provides operations to call the getOneDriveActivityUserCounts method.
// returns a *GetOneDriveActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveActivityUserCountsWithPeriod(period *string)(*GetOneDriveActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetOneDriveActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveActivityUserDetailWithDate provides operations to call the getOneDriveActivityUserDetail method.
// returns a *GetOneDriveActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetOneDriveActivityUserDetailWithDateRequestBuilder) {
    return NewGetOneDriveActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetOneDriveActivityUserDetailWithPeriod provides operations to call the getOneDriveActivityUserDetail method.
// returns a *GetOneDriveActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveActivityUserDetailWithPeriod(period *string)(*GetOneDriveActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetOneDriveActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveUsageAccountCountsWithPeriod provides operations to call the getOneDriveUsageAccountCounts method.
// returns a *GetOneDriveUsageAccountCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountCountsWithPeriod(period *string)(*GetOneDriveUsageAccountCountsWithPeriodRequestBuilder) {
    return NewGetOneDriveUsageAccountCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveUsageAccountDetailWithDate provides operations to call the getOneDriveUsageAccountDetail method.
// returns a *GetOneDriveUsageAccountDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetOneDriveUsageAccountDetailWithDateRequestBuilder) {
    return NewGetOneDriveUsageAccountDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetOneDriveUsageAccountDetailWithPeriod provides operations to call the getOneDriveUsageAccountDetail method.
// returns a *GetOneDriveUsageAccountDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountDetailWithPeriod(period *string)(*GetOneDriveUsageAccountDetailWithPeriodRequestBuilder) {
    return NewGetOneDriveUsageAccountDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveUsageFileCountsWithPeriod provides operations to call the getOneDriveUsageFileCounts method.
// returns a *GetOneDriveUsageFileCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveUsageFileCountsWithPeriod(period *string)(*GetOneDriveUsageFileCountsWithPeriodRequestBuilder) {
    return NewGetOneDriveUsageFileCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetOneDriveUsageStorageWithPeriod provides operations to call the getOneDriveUsageStorage method.
// returns a *GetOneDriveUsageStorageWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetOneDriveUsageStorageWithPeriod(period *string)(*GetOneDriveUsageStorageWithPeriodRequestBuilder) {
    return NewGetOneDriveUsageStorageWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTime provides operations to call the getPrinterArchivedPrintJobs method.
// returns a *GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTimeRequestBuilder when successful
func (m *ReportsRequestBuilder) GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, printerId *string, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)(*GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return NewGetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, printerId, startDateTime)
}
// GetRelyingPartyDetailedSummaryWithPeriod provides operations to call the getRelyingPartyDetailedSummary method.
// returns a *GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetRelyingPartyDetailedSummaryWithPeriod(period *string)(*GetRelyingPartyDetailedSummaryWithPeriodRequestBuilder) {
    return NewGetRelyingPartyDetailedSummaryWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointActivityFileCountsWithPeriod provides operations to call the getSharePointActivityFileCounts method.
// returns a *GetSharePointActivityFileCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointActivityFileCountsWithPeriod(period *string)(*GetSharePointActivityFileCountsWithPeriodRequestBuilder) {
    return NewGetSharePointActivityFileCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointActivityPagesWithPeriod provides operations to call the getSharePointActivityPages method.
// returns a *GetSharePointActivityPagesWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointActivityPagesWithPeriod(period *string)(*GetSharePointActivityPagesWithPeriodRequestBuilder) {
    return NewGetSharePointActivityPagesWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointActivityUserCountsWithPeriod provides operations to call the getSharePointActivityUserCounts method.
// returns a *GetSharePointActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointActivityUserCountsWithPeriod(period *string)(*GetSharePointActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetSharePointActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointActivityUserDetailWithDate provides operations to call the getSharePointActivityUserDetail method.
// returns a *GetSharePointActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetSharePointActivityUserDetailWithDateRequestBuilder) {
    return NewGetSharePointActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetSharePointActivityUserDetailWithPeriod provides operations to call the getSharePointActivityUserDetail method.
// returns a *GetSharePointActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointActivityUserDetailWithPeriod(period *string)(*GetSharePointActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetSharePointActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointSiteUsageDetailWithDate provides operations to call the getSharePointSiteUsageDetail method.
// returns a *GetSharePointSiteUsageDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsageDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetSharePointSiteUsageDetailWithDateRequestBuilder) {
    return NewGetSharePointSiteUsageDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetSharePointSiteUsageDetailWithPeriod provides operations to call the getSharePointSiteUsageDetail method.
// returns a *GetSharePointSiteUsageDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsageDetailWithPeriod(period *string)(*GetSharePointSiteUsageDetailWithPeriodRequestBuilder) {
    return NewGetSharePointSiteUsageDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointSiteUsageFileCountsWithPeriod provides operations to call the getSharePointSiteUsageFileCounts method.
// returns a *GetSharePointSiteUsageFileCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsageFileCountsWithPeriod(period *string)(*GetSharePointSiteUsageFileCountsWithPeriodRequestBuilder) {
    return NewGetSharePointSiteUsageFileCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointSiteUsagePagesWithPeriod provides operations to call the getSharePointSiteUsagePages method.
// returns a *GetSharePointSiteUsagePagesWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsagePagesWithPeriod(period *string)(*GetSharePointSiteUsagePagesWithPeriodRequestBuilder) {
    return NewGetSharePointSiteUsagePagesWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointSiteUsageSiteCountsWithPeriod provides operations to call the getSharePointSiteUsageSiteCounts method.
// returns a *GetSharePointSiteUsageSiteCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsageSiteCountsWithPeriod(period *string)(*GetSharePointSiteUsageSiteCountsWithPeriodRequestBuilder) {
    return NewGetSharePointSiteUsageSiteCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSharePointSiteUsageStorageWithPeriod provides operations to call the getSharePointSiteUsageStorage method.
// returns a *GetSharePointSiteUsageStorageWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSharePointSiteUsageStorageWithPeriod(period *string)(*GetSharePointSiteUsageStorageWithPeriodRequestBuilder) {
    return NewGetSharePointSiteUsageStorageWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessActivityCountsWithPeriod provides operations to call the getSkypeForBusinessActivityCounts method.
// returns a *GetSkypeForBusinessActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityCountsWithPeriod(period *string)(*GetSkypeForBusinessActivityCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessActivityUserCounts method.
// returns a *GetSkypeForBusinessActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserCountsWithPeriod(period *string)(*GetSkypeForBusinessActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessActivityUserDetailWithDate provides operations to call the getSkypeForBusinessActivityUserDetail method.
// returns a *GetSkypeForBusinessActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetSkypeForBusinessActivityUserDetailWithDateRequestBuilder) {
    return NewGetSkypeForBusinessActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetSkypeForBusinessActivityUserDetailWithPeriod provides operations to call the getSkypeForBusinessActivityUserDetail method.
// returns a *GetSkypeForBusinessActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserDetailWithPeriod(period *string)(*GetSkypeForBusinessActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageDistributionUserCounts method.
// returns a *GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriod(period *string)(*GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessDeviceUsageUserCountsWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageUserCounts method.
// returns a *GetSkypeForBusinessDeviceUsageUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserCountsWithPeriod(period *string)(*GetSkypeForBusinessDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessDeviceUsageUserDetailWithDate provides operations to call the getSkypeForBusinessDeviceUsageUserDetail method.
// returns a *GetSkypeForBusinessDeviceUsageUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetSkypeForBusinessDeviceUsageUserDetailWithDateRequestBuilder) {
    return NewGetSkypeForBusinessDeviceUsageUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetSkypeForBusinessDeviceUsageUserDetailWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageUserDetail method.
// returns a *GetSkypeForBusinessDeviceUsageUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserDetailWithPeriod(period *string)(*GetSkypeForBusinessDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessOrganizerActivityCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityCounts method.
// returns a *GetSkypeForBusinessOrganizerActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityCountsWithPeriod(period *string)(*GetSkypeForBusinessOrganizerActivityCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessOrganizerActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityMinuteCounts method.
// returns a *GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriod(period *string)(*GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessOrganizerActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityUserCounts method.
// returns a *GetSkypeForBusinessOrganizerActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityUserCountsWithPeriod(period *string)(*GetSkypeForBusinessOrganizerActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessOrganizerActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessParticipantActivityCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityCounts method.
// returns a *GetSkypeForBusinessParticipantActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityCountsWithPeriod(period *string)(*GetSkypeForBusinessParticipantActivityCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessParticipantActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityMinuteCounts method.
// returns a *GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriod(period *string)(*GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessParticipantActivityMinuteCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessParticipantActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityUserCounts method.
// returns a *GetSkypeForBusinessParticipantActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityUserCountsWithPeriod(period *string)(*GetSkypeForBusinessParticipantActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessParticipantActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessPeerToPeerActivityCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityCounts method.
// returns a *GetSkypeForBusinessPeerToPeerActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityCountsWithPeriod(period *string)(*GetSkypeForBusinessPeerToPeerActivityCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessPeerToPeerActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityMinuteCounts method.
// returns a *GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriod(period *string)(*GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityUserCounts method.
// returns a *GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriod(period *string)(*GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getTeamsDeviceUsageDistributionUserCounts method.
// returns a *GetTeamsDeviceUsageDistributionUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageDistributionUserCountsWithPeriod(period *string)(*GetTeamsDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return NewGetTeamsDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsDeviceUsageUserCountsWithPeriod provides operations to call the getTeamsDeviceUsageUserCounts method.
// returns a *GetTeamsDeviceUsageUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserCountsWithPeriod(period *string)(*GetTeamsDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return NewGetTeamsDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsDeviceUsageUserDetailWithDate provides operations to call the getTeamsDeviceUsageUserDetail method.
// returns a *GetTeamsDeviceUsageUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetTeamsDeviceUsageUserDetailWithDateRequestBuilder) {
    return NewGetTeamsDeviceUsageUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetTeamsDeviceUsageUserDetailWithPeriod provides operations to call the getTeamsDeviceUsageUserDetail method.
// returns a *GetTeamsDeviceUsageUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserDetailWithPeriod(period *string)(*GetTeamsDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return NewGetTeamsDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsTeamActivityCountsWithPeriod provides operations to call the getTeamsTeamActivityCounts method.
// returns a *GetTeamsTeamActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsTeamActivityCountsWithPeriod(period *string)(*GetTeamsTeamActivityCountsWithPeriodRequestBuilder) {
    return NewGetTeamsTeamActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsTeamActivityDetailWithDate provides operations to call the getTeamsTeamActivityDetail method.
// returns a *GetTeamsTeamActivityDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsTeamActivityDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetTeamsTeamActivityDetailWithDateRequestBuilder) {
    return NewGetTeamsTeamActivityDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetTeamsTeamActivityDetailWithPeriod provides operations to call the getTeamsTeamActivityDetail method.
// returns a *GetTeamsTeamActivityDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsTeamActivityDetailWithPeriod(period *string)(*GetTeamsTeamActivityDetailWithPeriodRequestBuilder) {
    return NewGetTeamsTeamActivityDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsTeamActivityDistributionCountsWithPeriod provides operations to call the getTeamsTeamActivityDistributionCounts method.
// returns a *GetTeamsTeamActivityDistributionCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsTeamActivityDistributionCountsWithPeriod(period *string)(*GetTeamsTeamActivityDistributionCountsWithPeriodRequestBuilder) {
    return NewGetTeamsTeamActivityDistributionCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsTeamCountsWithPeriod provides operations to call the getTeamsTeamCounts method.
// returns a *GetTeamsTeamCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsTeamCountsWithPeriod(period *string)(*GetTeamsTeamCountsWithPeriodRequestBuilder) {
    return NewGetTeamsTeamCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsUserActivityCountsWithPeriod provides operations to call the getTeamsUserActivityCounts method.
// returns a *GetTeamsUserActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsUserActivityCountsWithPeriod(period *string)(*GetTeamsUserActivityCountsWithPeriodRequestBuilder) {
    return NewGetTeamsUserActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsUserActivityUserCountsWithPeriod provides operations to call the getTeamsUserActivityUserCounts method.
// returns a *GetTeamsUserActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserCountsWithPeriod(period *string)(*GetTeamsUserActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetTeamsUserActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetTeamsUserActivityUserDetailWithDate provides operations to call the getTeamsUserActivityUserDetail method.
// returns a *GetTeamsUserActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetTeamsUserActivityUserDetailWithDateRequestBuilder) {
    return NewGetTeamsUserActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetTeamsUserActivityUserDetailWithPeriod provides operations to call the getTeamsUserActivityUserDetail method.
// returns a *GetTeamsUserActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserDetailWithPeriod(period *string)(*GetTeamsUserActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetTeamsUserActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTime provides operations to call the getUserArchivedPrintJobs method.
// returns a *GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTimeRequestBuilder when successful
func (m *ReportsRequestBuilder) GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, userId *string)(*GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return NewGetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, endDateTime, startDateTime, userId)
}
// GetYammerActivityCountsWithPeriod provides operations to call the getYammerActivityCounts method.
// returns a *GetYammerActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerActivityCountsWithPeriod(period *string)(*GetYammerActivityCountsWithPeriodRequestBuilder) {
    return NewGetYammerActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerActivityUserCountsWithPeriod provides operations to call the getYammerActivityUserCounts method.
// returns a *GetYammerActivityUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerActivityUserCountsWithPeriod(period *string)(*GetYammerActivityUserCountsWithPeriodRequestBuilder) {
    return NewGetYammerActivityUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerActivityUserDetailWithDate provides operations to call the getYammerActivityUserDetail method.
// returns a *GetYammerActivityUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetYammerActivityUserDetailWithDateRequestBuilder) {
    return NewGetYammerActivityUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetYammerActivityUserDetailWithPeriod provides operations to call the getYammerActivityUserDetail method.
// returns a *GetYammerActivityUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerActivityUserDetailWithPeriod(period *string)(*GetYammerActivityUserDetailWithPeriodRequestBuilder) {
    return NewGetYammerActivityUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getYammerDeviceUsageDistributionUserCounts method.
// returns a *GetYammerDeviceUsageDistributionUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerDeviceUsageDistributionUserCountsWithPeriod(period *string)(*GetYammerDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return NewGetYammerDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerDeviceUsageUserCountsWithPeriod provides operations to call the getYammerDeviceUsageUserCounts method.
// returns a *GetYammerDeviceUsageUserCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserCountsWithPeriod(period *string)(*GetYammerDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return NewGetYammerDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerDeviceUsageUserDetailWithDate provides operations to call the getYammerDeviceUsageUserDetail method.
// returns a *GetYammerDeviceUsageUserDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetYammerDeviceUsageUserDetailWithDateRequestBuilder) {
    return NewGetYammerDeviceUsageUserDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetYammerDeviceUsageUserDetailWithPeriod provides operations to call the getYammerDeviceUsageUserDetail method.
// returns a *GetYammerDeviceUsageUserDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserDetailWithPeriod(period *string)(*GetYammerDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return NewGetYammerDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerGroupsActivityCountsWithPeriod provides operations to call the getYammerGroupsActivityCounts method.
// returns a *GetYammerGroupsActivityCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerGroupsActivityCountsWithPeriod(period *string)(*GetYammerGroupsActivityCountsWithPeriodRequestBuilder) {
    return NewGetYammerGroupsActivityCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerGroupsActivityDetailWithDate provides operations to call the getYammerGroupsActivityDetail method.
// returns a *GetYammerGroupsActivityDetailWithDateRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerGroupsActivityDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*GetYammerGroupsActivityDetailWithDateRequestBuilder) {
    return NewGetYammerGroupsActivityDetailWithDateRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, date)
}
// GetYammerGroupsActivityDetailWithPeriod provides operations to call the getYammerGroupsActivityDetail method.
// returns a *GetYammerGroupsActivityDetailWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerGroupsActivityDetailWithPeriod(period *string)(*GetYammerGroupsActivityDetailWithPeriodRequestBuilder) {
    return NewGetYammerGroupsActivityDetailWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// GetYammerGroupsActivityGroupCountsWithPeriod provides operations to call the getYammerGroupsActivityGroupCounts method.
// returns a *GetYammerGroupsActivityGroupCountsWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) GetYammerGroupsActivityGroupCountsWithPeriod(period *string)(*GetYammerGroupsActivityGroupCountsWithPeriodRequestBuilder) {
    return NewGetYammerGroupsActivityGroupCountsWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// ManagedDeviceEnrollmentFailureDetails provides operations to call the managedDeviceEnrollmentFailureDetails method.
// returns a *ManagedDeviceEnrollmentFailureDetailsRequestBuilder when successful
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentFailureDetails()(*ManagedDeviceEnrollmentFailureDetailsRequestBuilder) {
    return NewManagedDeviceEnrollmentFailureDetailsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipToken provides operations to call the managedDeviceEnrollmentFailureDetails method.
// returns a *ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipTokenRequestBuilder when successful
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipToken(filter *string, skip *int32, skipToken *string, top *int32)(*ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipTokenRequestBuilder) {
    return NewManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipTokenRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, filter, skip, skipToken, top)
}
// ManagedDeviceEnrollmentTopFailures provides operations to call the managedDeviceEnrollmentTopFailures method.
// returns a *ManagedDeviceEnrollmentTopFailuresRequestBuilder when successful
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentTopFailures()(*ManagedDeviceEnrollmentTopFailuresRequestBuilder) {
    return NewManagedDeviceEnrollmentTopFailuresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ManagedDeviceEnrollmentTopFailuresWithPeriod provides operations to call the managedDeviceEnrollmentTopFailures method.
// returns a *ManagedDeviceEnrollmentTopFailuresWithPeriodRequestBuilder when successful
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentTopFailuresWithPeriod(period *string)(*ManagedDeviceEnrollmentTopFailuresWithPeriodRequestBuilder) {
    return NewManagedDeviceEnrollmentTopFailuresWithPeriodRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, period)
}
// MonthlyPrintUsageByPrinter provides operations to manage the monthlyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
// returns a *MonthlyPrintUsageByPrinterRequestBuilder when successful
func (m *ReportsRequestBuilder) MonthlyPrintUsageByPrinter()(*MonthlyPrintUsageByPrinterRequestBuilder) {
    return NewMonthlyPrintUsageByPrinterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MonthlyPrintUsageByUser provides operations to manage the monthlyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
// returns a *MonthlyPrintUsageByUserRequestBuilder when successful
func (m *ReportsRequestBuilder) MonthlyPrintUsageByUser()(*MonthlyPrintUsageByUserRequestBuilder) {
    return NewMonthlyPrintUsageByUserRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Partners provides operations to manage the partners property of the microsoft.graph.reportRoot entity.
// returns a *PartnersRequestBuilder when successful
func (m *ReportsRequestBuilder) Partners()(*PartnersRequestBuilder) {
    return NewPartnersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update the properties of a reportRoot object.
// returns a ReportRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-deviceconfig-reportroot-update?view=graph-rest-1.0
func (m *ReportsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateReportRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable), nil
}
// Security provides operations to manage the security property of the microsoft.graph.reportRoot entity.
// returns a *SecurityRequestBuilder when successful
func (m *ReportsRequestBuilder) Security()(*SecurityRequestBuilder) {
    return NewSecurityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation read properties and relationships of the reportRoot object.
// returns a *RequestInformation when successful
func (m *ReportsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a reportRoot object.
// returns a *RequestInformation when successful
func (m *ReportsRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ReportsRequestBuilder when successful
func (m *ReportsRequestBuilder) WithUrl(rawUrl string)(*ReportsRequestBuilder) {
    return NewReportsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
