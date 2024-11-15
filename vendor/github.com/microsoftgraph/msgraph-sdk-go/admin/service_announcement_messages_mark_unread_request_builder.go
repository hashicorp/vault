package admin

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ServiceAnnouncementMessagesMarkUnreadRequestBuilder provides operations to call the markUnread method.
type ServiceAnnouncementMessagesMarkUnreadRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ServiceAnnouncementMessagesMarkUnreadRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ServiceAnnouncementMessagesMarkUnreadRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewServiceAnnouncementMessagesMarkUnreadRequestBuilderInternal instantiates a new ServiceAnnouncementMessagesMarkUnreadRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesMarkUnreadRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesMarkUnreadRequestBuilder) {
    m := &ServiceAnnouncementMessagesMarkUnreadRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/admin/serviceAnnouncement/messages/markUnread", pathParameters),
    }
    return m
}
// NewServiceAnnouncementMessagesMarkUnreadRequestBuilder instantiates a new ServiceAnnouncementMessagesMarkUnreadRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesMarkUnreadRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesMarkUnreadRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewServiceAnnouncementMessagesMarkUnreadRequestBuilderInternal(urlParams, requestAdapter)
}
// Post mark a list of serviceUpdateMessages as unread for the signed in user.
// Deprecated: This method is obsolete. Use PostAsMarkUnreadPostResponse instead.
// returns a ServiceAnnouncementMessagesMarkUnreadResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-markunread?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesMarkUnreadRequestBuilder) Post(ctx context.Context, body ServiceAnnouncementMessagesMarkUnreadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkUnreadRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesMarkUnreadResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesMarkUnreadResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesMarkUnreadResponseable), nil
}
// PostAsMarkUnreadPostResponse mark a list of serviceUpdateMessages as unread for the signed in user.
// returns a ServiceAnnouncementMessagesMarkUnreadPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-markunread?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesMarkUnreadRequestBuilder) PostAsMarkUnreadPostResponse(ctx context.Context, body ServiceAnnouncementMessagesMarkUnreadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkUnreadRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesMarkUnreadPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesMarkUnreadPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesMarkUnreadPostResponseable), nil
}
// ToPostRequestInformation mark a list of serviceUpdateMessages as unread for the signed in user.
// returns a *RequestInformation when successful
func (m *ServiceAnnouncementMessagesMarkUnreadRequestBuilder) ToPostRequestInformation(ctx context.Context, body ServiceAnnouncementMessagesMarkUnreadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkUnreadRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ServiceAnnouncementMessagesMarkUnreadRequestBuilder when successful
func (m *ServiceAnnouncementMessagesMarkUnreadRequestBuilder) WithUrl(rawUrl string)(*ServiceAnnouncementMessagesMarkUnreadRequestBuilder) {
    return NewServiceAnnouncementMessagesMarkUnreadRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
