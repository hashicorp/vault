package communications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder provides operations to call the sendVirtualAppointmentReminderSms method.
type OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderInternal instantiates a new OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder and sets the default values.
func NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) {
    m := &OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/communications/onlineMeetings/{onlineMeeting%2Did}/sendVirtualAppointmentReminderSms", pathParameters),
    }
    return m
}
// NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder instantiates a new OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder and sets the default values.
func NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderInternal(urlParams, requestAdapter)
}
// Post send an SMS reminder to external attendees for a Teams virtual appointment. This feature requires Teams premium and attendees must have a valid United States phone number to receive SMS notifications.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/virtualappointment-sendvirtualappointmentremindersms?view=graph-rest-1.0
func (m *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) Post(ctx context.Context, body OnlineMeetingsItemSendVirtualAppointmentReminderSmsPostRequestBodyable, requestConfiguration *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation send an SMS reminder to external attendees for a Teams virtual appointment. This feature requires Teams premium and attendees must have a valid United States phone number to receive SMS notifications.
// returns a *RequestInformation when successful
func (m *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) ToPostRequestInformation(ctx context.Context, body OnlineMeetingsItemSendVirtualAppointmentReminderSmsPostRequestBodyable, requestConfiguration *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder when successful
func (m *OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) WithUrl(rawUrl string)(*OnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder) {
    return NewOnlineMeetingsItemSendVirtualAppointmentReminderSmsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
