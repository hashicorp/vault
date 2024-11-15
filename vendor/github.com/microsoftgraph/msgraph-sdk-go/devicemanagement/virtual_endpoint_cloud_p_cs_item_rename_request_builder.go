package devicemanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// VirtualEndpointCloudPCsItemRenameRequestBuilder provides operations to call the rename method.
type VirtualEndpointCloudPCsItemRenameRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// VirtualEndpointCloudPCsItemRenameRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type VirtualEndpointCloudPCsItemRenameRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewVirtualEndpointCloudPCsItemRenameRequestBuilderInternal instantiates a new VirtualEndpointCloudPCsItemRenameRequestBuilder and sets the default values.
func NewVirtualEndpointCloudPCsItemRenameRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointCloudPCsItemRenameRequestBuilder) {
    m := &VirtualEndpointCloudPCsItemRenameRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceManagement/virtualEndpoint/cloudPCs/{cloudPC%2Did}/rename", pathParameters),
    }
    return m
}
// NewVirtualEndpointCloudPCsItemRenameRequestBuilder instantiates a new VirtualEndpointCloudPCsItemRenameRequestBuilder and sets the default values.
func NewVirtualEndpointCloudPCsItemRenameRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*VirtualEndpointCloudPCsItemRenameRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewVirtualEndpointCloudPCsItemRenameRequestBuilderInternal(urlParams, requestAdapter)
}
// Post rename a specific cloudPC object. Use this API to update the displayName of a Cloud PC entity.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/cloudpc-rename?view=graph-rest-1.0
func (m *VirtualEndpointCloudPCsItemRenameRequestBuilder) Post(ctx context.Context, body VirtualEndpointCloudPCsItemRenamePostRequestBodyable, requestConfiguration *VirtualEndpointCloudPCsItemRenameRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation rename a specific cloudPC object. Use this API to update the displayName of a Cloud PC entity.
// returns a *RequestInformation when successful
func (m *VirtualEndpointCloudPCsItemRenameRequestBuilder) ToPostRequestInformation(ctx context.Context, body VirtualEndpointCloudPCsItemRenamePostRequestBodyable, requestConfiguration *VirtualEndpointCloudPCsItemRenameRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *VirtualEndpointCloudPCsItemRenameRequestBuilder when successful
func (m *VirtualEndpointCloudPCsItemRenameRequestBuilder) WithUrl(rawUrl string)(*VirtualEndpointCloudPCsItemRenameRequestBuilder) {
    return NewVirtualEndpointCloudPCsItemRenameRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
