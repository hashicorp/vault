package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder provides operations to count the resources in the collection.
type VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetQueryParameters get the number of the resource
type VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetQueryParameters struct {
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
}
// VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetQueryParameters
}
// NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderInternal instantiates a new VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder and sets the default values.
func NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) {
    m := &VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/virtualEvents/townhalls/{virtualEventTownhall%2Did}/sessions/{virtualEventSession%2Did}/attendanceReports/{meetingAttendanceReport%2Did}/attendanceRecords/$count{?%24filter,%24search}", pathParameters),
    }
    return m
}
// NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder instantiates a new VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder and sets the default values.
func NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get the number of the resource
// returns a *int32 when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) Get(ctx context.Context, requestConfiguration *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetRequestConfiguration)(*int32, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.SendPrimitive(ctx, requestInfo, "int32", errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(*int32), nil
}
// ToGetRequestInformation get the number of the resource
// returns a *RequestInformation when successful
func (m *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "text/plain;q=0.9")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder when successful
func (m *VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) WithUrl(rawUrl string)(*VirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder) {
    return NewVirtualEventsTownhallsItemSessionsItemAttendanceReportsItemAttendanceRecordsCountRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
