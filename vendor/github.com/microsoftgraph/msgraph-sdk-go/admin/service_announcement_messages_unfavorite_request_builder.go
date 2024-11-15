package admin

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ServiceAnnouncementMessagesUnfavoriteRequestBuilder provides operations to call the unfavorite method.
type ServiceAnnouncementMessagesUnfavoriteRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ServiceAnnouncementMessagesUnfavoriteRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ServiceAnnouncementMessagesUnfavoriteRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewServiceAnnouncementMessagesUnfavoriteRequestBuilderInternal instantiates a new ServiceAnnouncementMessagesUnfavoriteRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesUnfavoriteRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesUnfavoriteRequestBuilder) {
    m := &ServiceAnnouncementMessagesUnfavoriteRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/admin/serviceAnnouncement/messages/unfavorite", pathParameters),
    }
    return m
}
// NewServiceAnnouncementMessagesUnfavoriteRequestBuilder instantiates a new ServiceAnnouncementMessagesUnfavoriteRequestBuilder and sets the default values.
func NewServiceAnnouncementMessagesUnfavoriteRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ServiceAnnouncementMessagesUnfavoriteRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewServiceAnnouncementMessagesUnfavoriteRequestBuilderInternal(urlParams, requestAdapter)
}
// Post remove the favorite status of serviceUpdateMessages for the signed in user.
// Deprecated: This method is obsolete. Use PostAsUnfavoritePostResponse instead.
// returns a ServiceAnnouncementMessagesUnfavoriteResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-unfavorite?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesUnfavoriteRequestBuilder) Post(ctx context.Context, body ServiceAnnouncementMessagesUnfavoritePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnfavoriteRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesUnfavoriteResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesUnfavoriteResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesUnfavoriteResponseable), nil
}
// PostAsUnfavoritePostResponse remove the favorite status of serviceUpdateMessages for the signed in user.
// returns a ServiceAnnouncementMessagesUnfavoritePostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/serviceupdatemessage-unfavorite?view=graph-rest-1.0
func (m *ServiceAnnouncementMessagesUnfavoriteRequestBuilder) PostAsUnfavoritePostResponse(ctx context.Context, body ServiceAnnouncementMessagesUnfavoritePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnfavoriteRequestBuilderPostRequestConfiguration)(ServiceAnnouncementMessagesUnfavoritePostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateServiceAnnouncementMessagesUnfavoritePostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ServiceAnnouncementMessagesUnfavoritePostResponseable), nil
}
// ToPostRequestInformation remove the favorite status of serviceUpdateMessages for the signed in user.
// returns a *RequestInformation when successful
func (m *ServiceAnnouncementMessagesUnfavoriteRequestBuilder) ToPostRequestInformation(ctx context.Context, body ServiceAnnouncementMessagesUnfavoritePostRequestBodyable, requestConfiguration *ServiceAnnouncementMessagesUnfavoriteRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ServiceAnnouncementMessagesUnfavoriteRequestBuilder when successful
func (m *ServiceAnnouncementMessagesUnfavoriteRequestBuilder) WithUrl(rawUrl string)(*ServiceAnnouncementMessagesUnfavoriteRequestBuilder) {
    return NewServiceAnnouncementMessagesUnfavoriteRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
