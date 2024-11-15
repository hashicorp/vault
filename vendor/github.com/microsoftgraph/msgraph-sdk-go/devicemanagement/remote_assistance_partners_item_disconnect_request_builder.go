package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// RemoteAssistancePartnersItemDisconnectRequestBuilder provides operations to call the disconnect method.
type RemoteAssistancePartnersItemDisconnectRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// RemoteAssistancePartnersItemDisconnectRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type RemoteAssistancePartnersItemDisconnectRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewRemoteAssistancePartnersItemDisconnectRequestBuilderInternal instantiates a new RemoteAssistancePartnersItemDisconnectRequestBuilder and sets the default values.
func NewRemoteAssistancePartnersItemDisconnectRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RemoteAssistancePartnersItemDisconnectRequestBuilder) {
    m := &RemoteAssistancePartnersItemDisconnectRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/remoteAssistancePartners/{remoteAssistancePartner%2Did}/disconnect", pathParameters),
    }
    return m
}
// NewRemoteAssistancePartnersItemDisconnectRequestBuilder instantiates a new RemoteAssistancePartnersItemDisconnectRequestBuilder and sets the default values.
func NewRemoteAssistancePartnersItemDisconnectRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*RemoteAssistancePartnersItemDisconnectRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewRemoteAssistancePartnersItemDisconnectRequestBuilderInternal(urlParams, requestAdapter)
}
// Post a request to remove the active TeamViewer connector
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-remoteassistance-remoteassistancepartner-disconnect?view=graph-rest-1.0
func (m *RemoteAssistancePartnersItemDisconnectRequestBuilder) Post(ctx context.Context, requestConfiguration *RemoteAssistancePartnersItemDisconnectRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation a request to remove the active TeamViewer connector
// returns a *RequestInformation when successful
func (m *RemoteAssistancePartnersItemDisconnectRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *RemoteAssistancePartnersItemDisconnectRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *RemoteAssistancePartnersItemDisconnectRequestBuilder when successful
func (m *RemoteAssistancePartnersItemDisconnectRequestBuilder) WithUrl(rawUrl string)(*RemoteAssistancePartnersItemDisconnectRequestBuilder) {
    return NewRemoteAssistancePartnersItemDisconnectRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
