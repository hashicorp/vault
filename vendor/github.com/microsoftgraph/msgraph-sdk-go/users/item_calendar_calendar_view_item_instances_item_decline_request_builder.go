package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder provides operations to call the decline method.
type ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderInternal instantiates a new ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder and sets the default values.
func NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) {
    m := &ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/calendar/calendarView/{event%2Did}/instances/{event%2Did1}/decline", pathParameters),
    }
    return m
}
// NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder instantiates a new ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder and sets the default values.
func NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderInternal(urlParams, requestAdapter)
}
// Post decline invitation to the specified event in a user calendar. If the event allows proposals for new times, on declining the event, an invitee can choose to suggest an alternative time by including the proposedNewTime parameter. For more information on how to propose a time, and how to receive and accept a new time proposal, see Propose new meeting times.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/event-decline?view=graph-rest-1.0
func (m *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) Post(ctx context.Context, body ItemCalendarCalendarViewItemInstancesItemDeclinePostRequestBodyable, requestConfiguration *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
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
// ToPostRequestInformation decline invitation to the specified event in a user calendar. If the event allows proposals for new times, on declining the event, an invitee can choose to suggest an alternative time by including the proposedNewTime parameter. For more information on how to propose a time, and how to receive and accept a new time proposal, see Propose new meeting times.
// returns a *RequestInformation when successful
func (m *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemCalendarCalendarViewItemInstancesItemDeclinePostRequestBodyable, requestConfiguration *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder when successful
func (m *ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) WithUrl(rawUrl string)(*ItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder) {
    return NewItemCalendarCalendarViewItemInstancesItemDeclineRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
