package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder provides operations to call the sendVirtualAppointmentSms method.
type ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderInternal instantiates a new ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder and sets the default values.
func NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) {
    m := &ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/onlineMeetings/{onlineMeeting%2Did}/sendVirtualAppointmentSms", pathParameters),
    }
    return m
}
// NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder instantiates a new ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder and sets the default values.
func NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderInternal(urlParams, requestAdapter)
}
// Post send an SMS notification to external attendees when a Teams virtual appointment is confirmed, rescheduled, or canceled. This feature requires Teams premium. Attendees must have a valid United States phone number to receive these SMS notifications.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualappointment-sendvirtualappointmentsms?view=graph-rest-1.0
func (m *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) Post(ctx context.Context, body ItemOnlineMeetingsItemSendVirtualAppointmentSmsPostRequestBodyable, requestConfiguration *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation send an SMS notification to external attendees when a Teams virtual appointment is confirmed, rescheduled, or canceled. This feature requires Teams premium. Attendees must have a valid United States phone number to receive these SMS notifications.
// returns a *RequestInformation when successful
func (m *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemOnlineMeetingsItemSendVirtualAppointmentSmsPostRequestBodyable, requestConfiguration *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder when successful
func (m *ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) WithUrl(rawUrl string)(*ItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder) {
    return NewItemOnlineMeetingsItemSendVirtualAppointmentSmsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
