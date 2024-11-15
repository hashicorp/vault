package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder provides operations to manage the localizedNotificationMessages property of the microsoft.graph.notificationMessageTemplate entity.
type NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetQueryParameters read properties and relationships of the localizedNotificationMessage object.
type NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetQueryParameters
}
// NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderInternal instantiates a new NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder and sets the default values.
func NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) {
    m := &NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/notificationMessageTemplates/{notificationMessageTemplate%2Did}/localizedNotificationMessages/{localizedNotificationMessage%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder instantiates a new NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder and sets the default values.
func NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete deletes a localizedNotificationMessage.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-notification-localizednotificationmessage-delete?view=graph-rest-1.0
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read properties and relationships of the localizedNotificationMessage object.
// returns a LocalizedNotificationMessageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-notification-localizednotificationmessage-get?view=graph-rest-1.0
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) Get(ctx context.Context, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateLocalizedNotificationMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable), nil
}
// Patch update the properties of a localizedNotificationMessage object.
// returns a LocalizedNotificationMessageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-notification-localizednotificationmessage-update?view=graph-rest-1.0
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateLocalizedNotificationMessageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable), nil
}
// ToDeleteRequestInformation deletes a localizedNotificationMessage.
// returns a *RequestInformation when successful
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read properties and relationships of the localizedNotificationMessage object.
// returns a *RequestInformation when successful
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a localizedNotificationMessage object.
// returns a *RequestInformation when successful
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.LocalizedNotificationMessageable, requestConfiguration *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder when successful
func (m *NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) WithUrl(rawUrl string)(*NotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder) {
    return NewNotificationMessageTemplatesItemLocalizedNotificationMessagesLocalizedNotificationMessageItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
