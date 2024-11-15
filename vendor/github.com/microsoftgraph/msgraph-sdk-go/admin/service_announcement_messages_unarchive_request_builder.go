package admin

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ServiceAnnouncementMessagesUnarchiveRequestBuilder provides operations to call the unarchive method.
type ServiceAnnouncementMessagesUnarchiveRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ServiceAnnouncementMessagesUnarchiveRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ServiceAnnouncementMessagesUnarchiveRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewServiceAnnouncementMessagesUnarchiveRequestBuilderInternal instantiates a new ServiceAnnouncementMessagesUnarchiveRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesUnarchiveRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesUnarchiveRequestBuilder) {
    m := &ServiceAnnouncementMessagesUnarchiveRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/admin/serviceAnnouncement/messages/unarchive", pathParameters),
    }
    return m
}
// NewServiceAnnouncementMessagesUnarchiveRequestBuilder instantiates a new ServiceAnnouncementMessagesUnarchiveRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesUnarchiveRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesUnarchiveRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewServiceAnnouncementMessagesUnarchiveRequestBuilderInternal(urlParams, requestAdapter)
}
// Post unarchive a list of serviceUpdateMessages for the signed in user.
// Deprecated: This method is obsolete. Use PostAsUnarchivePostResponse instead.
// returns a ServiceAnnouncementMessagesUnarchiveResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-unarchive?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesUnarchiveRequestBuilder) Post(ctx context.Context, body ServiceAnnouncementMessagesUnarchivePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnarchiveRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesUnarchiveResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesUnarchiveResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesUnarchiveResponseable), nil
}
// PostAsUnarchivePostResponse unarchive a list of serviceUpdateMessages for the signed in user.
// returns a ServiceAnnouncementMessagesUnarchivePostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-unarchive?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesUnarchiveRequestBuilder) PostAsUnarchivePostResponse(ctx context.Context, body ServiceAnnouncementMessagesUnarchivePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnarchiveRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesUnarchivePostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesUnarchivePostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesUnarchivePostResponseable), nil
}
// ToPostRequestInformation unarchive a list of serviceUpdateMessages for the signed in user.
// returns a *RequestInformation when successful
func (m *ServiceAnnouncementMessagesUnarchiveRequestBuilder) ToPostRequestInformation(ctx context.Context, body ServiceAnnouncementMessagesUnarchivePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnarchiveRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ServiceAnnouncementMessagesUnarchiveRequestBuilder when successful
func (m *ServiceAnnouncementMessagesUnarchiveRequestBuilder) WithUrl(rawUrl string)(*ServiceAnnouncementMessagesUnarchiveRequestBuilder) {
    return NewServiceAnnouncementMessagesUnarchiveRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
