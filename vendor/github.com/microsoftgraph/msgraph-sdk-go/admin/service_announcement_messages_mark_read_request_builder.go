package admin

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ServiceAnnouncementMessagesMarkReadRequestBuilder provides operations to call the markRead method.
type ServiceAnnouncementMessagesMarkReadRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ServiceAnnouncementMessagesMarkReadRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ServiceAnnouncementMessagesMarkReadRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewServiceAnnouncementMessagesMarkReadRequestBuilderInternal instantiates a new ServiceAnnouncementMessagesMarkReadRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesMarkReadRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesMarkReadRequestBuilder) {
    m := &ServiceAnnouncementMessagesMarkReadRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/admin/serviceAnnouncement/messages/markRead", pathParameters),
    }
    return m
}
// NewServiceAnnouncementMessagesMarkReadRequestBuilder instantiates a new ServiceAnnouncementMessagesMarkReadRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesMarkReadRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesMarkReadRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewServiceAnnouncementMessagesMarkReadRequestBuilderInternal(urlParams, requestAdapter)
}
// Post mark a list of serviceUpdateMessages as read for the signed in user.
// Deprecated: This method is obsolete. Use PostAsMarkReadPostResponse instead.
// returns a ServiceAnnouncementMessagesMarkReadResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-markread?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesMarkReadRequestBuilder) Post(ctx context.Context, body ServiceAnnouncementMessagesMarkReadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkReadRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesMarkReadResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesMarkReadResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesMarkReadResponseable), nil
}
// PostAsMarkReadPostResponse mark a list of serviceUpdateMessages as read for the signed in user.
// returns a ServiceAnnouncementMessagesMarkReadPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-markread?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesMarkReadRequestBuilder) PostAsMarkReadPostResponse(ctx context.Context, body ServiceAnnouncementMessagesMarkReadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkReadRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesMarkReadPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesMarkReadPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesMarkReadPostResponseable), nil
}
// ToPostRequestInformation mark a list of serviceUpdateMessages as read for the signed in user.
// returns a *RequestInformation when successful
func (m *ServiceAnnouncementMessagesMarkReadRequestBuilder) ToPostRequestInformation(ctx context.Context, body ServiceAnnouncementMessagesMarkReadPostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesMarkReadRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ServiceAnnouncementMessagesMarkReadRequestBuilder when successful
func (m *ServiceAnnouncementMessagesMarkReadRequestBuilder) WithUrl(rawUrl string)(*ServiceAnnouncementMessagesMarkReadRequestBuilder) {
    return NewServiceAnnouncementMessagesMarkReadRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
