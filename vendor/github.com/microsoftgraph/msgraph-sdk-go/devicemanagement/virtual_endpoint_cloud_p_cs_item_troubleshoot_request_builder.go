package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEndpointCloudPCsItemTroubleshootRequestBuilder provides operations to call the troubleshoot method.
type VirtualEndpointCloudPCsItemTroubleshootRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEndpointCloudPCsItemTroubleshootRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEndpointCloudPCsItemTroubleshootRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilderInternal instantiates a new VirtualEndpointCloudPCsItemTroubleshootRequestBuilder and sets the default values.
func NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) {
    m := &VirtualEndpointCloudPCsItemTroubleshootRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/virtualEndpoint/cloudPCs/{cloudPC%2Did}/troubleshoot", pathParameters),
    }
    return m
}
// NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilder instantiates a new VirtualEndpointCloudPCsItemTroubleshootRequestBuilder and sets the default values.
func NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilderInternal(urlParams, requestAdapter)
}
// Post troubleshoot a specific cloudPC object. Use this API to check the health status of the Cloud PC and the session host.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/cloudpc-troubleshoot?view=graph-rest-1.0
func (m *VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) Post(ctx context.Context, requestConfiguration *VirtualEndpointCloudPCsItemTroubleshootRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation troubleshoot a specific cloudPC object. Use this API to check the health status of the Cloud PC and the session host.
// returns a *RequestInformation when successful
func (m *VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *VirtualEndpointCloudPCsItemTroubleshootRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *VirtualEndpointCloudPCsItemTroubleshootRequestBuilder when successful
func (m *VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) WithUrl(rawUrl string)(*VirtualEndpointCloudPCsItemTroubleshootRequestBuilder) {
    return NewVirtualEndpointCloudPCsItemTroubleshootRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
