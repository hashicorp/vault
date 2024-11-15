package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder provides operations to call the deleteUserFromSharedAppleDevice method.
type ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderInternal instantiates a new ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder and sets the default values.
func NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) {
    m := &ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/managedDevices/{managedDevice%2Did}/deleteUserFromSharedAppleDevice", pathParameters),
    }
    return m
}
// NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder instantiates a new ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder and sets the default values.
func NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderInternal(urlParams, requestAdapter)
}
// Post delete user from shared Apple device
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-1.0
func (m *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) Post(ctx context.Context, body ItemManagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBodyable, requestConfiguration *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation delete user from shared Apple device
// returns a *RequestInformation when successful
func (m *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemManagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBodyable, requestConfiguration *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder when successful
func (m *ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) WithUrl(rawUrl string)(*ItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder) {
    return NewItemManagedDevicesItemDeleteUserFromSharedAppleDeviceRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
