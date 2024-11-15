package organization

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSetMobileDeviceManagementAuthorityRequestBuilder provides operations to call the setMobileDeviceManagementAuthority method.
type ItemSetMobileDeviceManagementAuthorityRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSetMobileDeviceManagementAuthorityRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSetMobileDeviceManagementAuthorityRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSetMobileDeviceManagementAuthorityRequestBuilderInternal instantiates a new ItemSetMobileDeviceManagementAuthorityRequestBuilder and sets the default values.
func NewItemSetMobileDeviceManagementAuthorityRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSetMobileDeviceManagementAuthorityRequestBuilder) {
    m := &ItemSetMobileDeviceManagementAuthorityRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/organization/{organization%2Did}/setMobileDeviceManagementAuthority", pathParameters),
    }
    return m
}
// NewItemSetMobileDeviceManagementAuthorityRequestBuilder instantiates a new ItemSetMobileDeviceManagementAuthorityRequestBuilder and sets the default values.
func NewItemSetMobileDeviceManagementAuthorityRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSetMobileDeviceManagementAuthorityRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSetMobileDeviceManagementAuthorityRequestBuilderInternal(urlParams, requestAdapter)
}
// Post set mobile device management authority
// Deprecated: This method is obsolete. Use PostAsSetMobileDeviceManagementAuthorityPostResponse instead.
// returns a ItemSetMobileDeviceManagementAuthorityResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-onboarding-organization-setmobiledevicemanagementauthority?view=graph-rest-1.0
func (m *ItemSetMobileDeviceManagementAuthorityRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemSetMobileDeviceManagementAuthorityRequestBuilderPostRequestConfiguration)(ItemSetMobileDeviceManagementAuthorityResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSetMobileDeviceManagementAuthorityResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSetMobileDeviceManagementAuthorityResponseable), nil
}
// PostAsSetMobileDeviceManagementAuthorityPostResponse set mobile device management authority
// returns a ItemSetMobileDeviceManagementAuthorityPostResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-onboarding-organization-setmobiledevicemanagementauthority?view=graph-rest-1.0
func (m *ItemSetMobileDeviceManagementAuthorityRequestBuilder) PostAsSetMobileDeviceManagementAuthorityPostResponse(ctx context.Context, requestConfiguration *ItemSetMobileDeviceManagementAuthorityRequestBuilderPostRequestConfiguration)(ItemSetMobileDeviceManagementAuthorityPostResponseable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, CreateItemSetMobileDeviceManagementAuthorityPostResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ItemSetMobileDeviceManagementAuthorityPostResponseable), nil
}
// ToPostRequestInformation set mobile device management authority
// returns a *RequestInformation when successful
func (m *ItemSetMobileDeviceManagementAuthorityRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemSetMobileDeviceManagementAuthorityRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemSetMobileDeviceManagementAuthorityRequestBuilder when successful
func (m *ItemSetMobileDeviceManagementAuthorityRequestBuilder) WithUrl(rawUrl string)(*ItemSetMobileDeviceManagementAuthorityRequestBuilder) {
    return NewItemSetMobileDeviceManagementAuthorityRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
