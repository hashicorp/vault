package solutions

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder provides operations to call the cancel method.
type VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderInternal instantiates a new VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) {
    m := &VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/solutions/virtualEvents/webinars/{virtualEventWebinar%2Did}/registrations(userId='{userId}')/cancel", pathParameters),
    }
    return m
}
// NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder instantiates a new VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder and sets the default values.
func NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderInternal(urlParams, requestAdapter)
}
// Post cancel a registrant's registration record for a webinar. 
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) Post(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation cancel a registrant's registration record for a webinar. 
// returns a *RequestInformation when successful
func (m *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder when successful
func (m *VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) WithUrl(rawUrl string)(*VirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder) {
    return NewVirtualEventsWebinarsItemRegistrationsWithUserIdCancelRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
