package communications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder provides operations to manage the attendanceRecords property of the microsoft.graph.meetingAttendanceReport entity.
type OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetQueryParameters list of attendance records of an attendance report. Read-only.
type OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetQueryParameters
}
// OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderInternal instantiates a new OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder and sets the default values.
func NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) {
    m := &OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/communications/onlineMeetings/{onlineMeeting%2Did}/attendanceReports/{meetingAttendanceReport%2Did}/attendanceRecords/{attendanceRecord%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder instantiates a new OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder and sets the default values.
func NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property attendanceRecords for communications
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get list of attendance records of an attendance report. Read-only.
// returns a AttendanceRecordable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) Get(ctx context.Context, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAttendanceRecordFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable), nil
}
// Patch update the navigation property attendanceRecords in communications
// returns a AttendanceRecordable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAttendanceRecordFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable), nil
}
// ToDeleteRequestInformation delete navigation property attendanceRecords for communications
// returns a *RequestInformation when successful
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation list of attendance records of an attendance report. Read-only.
// returns a *RequestInformation when successful
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property attendanceRecords in communications
// returns a *RequestInformation when successful
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AttendanceRecordable, requestConfiguration *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder when successful
func (m *OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) WithUrl(rawUrl string)(*OnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder) {
    return NewOnlineMeetingsItemAttendanceReportsItemAttendanceRecordsAttendanceRecordItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
