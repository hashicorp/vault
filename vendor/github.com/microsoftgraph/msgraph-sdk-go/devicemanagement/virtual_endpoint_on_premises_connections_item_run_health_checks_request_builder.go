package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder provides operations to call the runHealthChecks method.
type VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderInternal instantiates a new VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder and sets the default values.
func NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) {
    m := &VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/virtualEndpoint/onPremisesConnections/{cloudPcOnPremisesConnection%2Did}/runHealthChecks", pathParameters),
    }
    return m
}
// NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder instantiates a new VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder and sets the default values.
func NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderInternal(urlParams, requestAdapter)
}
// Post run health checks on the cloudPcOnPremisesConnection object. It triggers a new health check for the cloudPcOnPremisesConnection (../resources/cloudpconpremisesconnection.md) object and changes the healthCheckStatus and [healthCheckStatusDetail properties when check finished.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/cloudpconpremisesconnection-runhealthcheck?view=graph-rest-1.0
func (m *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) Post(ctx context.Context, requestConfiguration *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation run health checks on the cloudPcOnPremisesConnection object. It triggers a new health check for the cloudPcOnPremisesConnection (../resources/cloudpconpremisesconnection.md) object and changes the healthCheckStatus and [healthCheckStatusDetail properties when check finished.
// returns a *RequestInformation when successful
func (m *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder when successful
func (m *VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) WithUrl(rawUrl string)(*VirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder) {
    return NewVirtualEndpointOnPremisesConnectionsItemRunHealthChecksRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
