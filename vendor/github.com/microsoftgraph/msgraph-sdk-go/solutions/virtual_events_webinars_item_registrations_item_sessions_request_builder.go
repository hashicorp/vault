package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder provides operations to manage the sessions property of the microsoft.graph.virtualEventRegistration entity.
type VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetQueryParameters get a list of sessions summaries that a registrant registered for in a webinar. A session summary contains only the endDateTime, id, joinWebUrl, startDateTime, and subject of a virtual event session. The rest of session properties will be null. To get all the properties of a virtualEventSession, use the Get virtualEventSession method. 
type VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetQueryParameters struct {
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
// VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetQueryParameters
}
// ByVirtualEventSessionId provides operations to manage the sessions property of the microsoft.graph.virtualEventRegistration entity.
// returns a *VirtualEventsWebinarsItemRegistrationsItemSessionsVirtualEventSessionItemRequestBuilder when successful
func (m *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) ByVirtualEventSessionId(virtualEventSessionId string)(*VirtualEventsWebinarsItemRegistrationsItemSessionsVirtualEventSessionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if virtualEventSessionId != "" {
        urlTplParams["virtualEventSession%2Did"] = virtualEventSessionId
    }
    return NewVirtualEventsWebinarsItemRegistrationsItemSessionsVirtualEventSessionItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderInternal instantiates a new VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) {
    m := &VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/virtualEvents/webinars/{virtualEventWebinar%2Did}/registrations/{virtualEventRegistration%2Did}/sessions{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder instantiates a new VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *VirtualEventsWebinarsItemRegistrationsItemSessionsCountRequestBuilder when successful
func (m *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) Count()(*VirtualEventsWebinarsItemRegistrationsItemSessionsCountRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsItemSessionsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get a list of sessions summaries that a registrant registered for in a webinar. A session summary contains only the endDateTime, id, joinWebUrl, startDateTime, and subject of a virtual event session. The rest of session properties will be null. To get all the properties of a virtualEventSession, use the Get virtualEventSession method. 
// returns a VirtualEventSessionCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualeventregistration-list-sessions?view=graph-rest-1.0
func (m *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) Get(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventSessionCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateVirtualEventSessionCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.VirtualEventSessionCollectionResponseable), nil
}
// ToGetRequestInformation get a list of sessions summaries that a registrant registered for in a webinar. A session summary contains only the endDateTime, id, joinWebUrl, startDateTime, and subject of a virtual event session. The rest of session properties will be null. To get all the properties of a virtualEventSession, use the Get virtualEventSession method. 
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder when successful
func (m *VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) WithUrl(rawUrl string)(*VirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsItemSessionsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
