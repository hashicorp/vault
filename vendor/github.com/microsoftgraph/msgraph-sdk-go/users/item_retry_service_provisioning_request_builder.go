package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemRetryServiceProvisioningRequestBuilder provides operations to call the retryServiceProvisioning method.
type ItemRetryServiceProvisioningRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemRetryServiceProvisioningRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemRetryServiceProvisioningRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemRetryServiceProvisioningRequestBuilderInternal instantiates a new ItemRetryServiceProvisioningRequestBuilder and sets the default values.
func NewItemRetryServiceProvisioningRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRetryServiceProvisioningRequestBuilder) {
    m := &ItemRetryServiceProvisioningRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/retryServiceProvisioning", pathParameters),
    }
    return m
}
// NewItemRetryServiceProvisioningRequestBuilder instantiates a new ItemRetryServiceProvisioningRequestBuilder and sets the default values.
func NewItemRetryServiceProvisioningRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemRetryServiceProvisioningRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemRetryServiceProvisioningRequestBuilderInternal(urlParams, requestAdapter)
}
// Post retry the user service provisioning.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/user-retryserviceprovisioning?view=graph-rest-1.0
func (m *ItemRetryServiceProvisioningRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemRetryServiceProvisioningRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation retry the user service provisioning.
// returns a *RequestInformation when successful
func (m *ItemRetryServiceProvisioningRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemRetryServiceProvisioningRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemRetryServiceProvisioningRequestBuilder when successful
func (m *ItemRetryServiceProvisioningRequestBuilder) WithUrl(rawUrl string)(*ItemRetryServiceProvisioningRequestBuilder) {
    return NewItemRetryServiceProvisioningRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
