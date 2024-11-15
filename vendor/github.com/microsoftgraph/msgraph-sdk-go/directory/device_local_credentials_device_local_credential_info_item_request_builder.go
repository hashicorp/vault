package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder provides operations to manage the deviceLocalCredentials property of the microsoft.graph.directory entity.
type DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetQueryParameters retrieve the properties of a deviceLocalCredentialInfo for a specified device object. 
type DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetQueryParameters
}
// DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderInternal instantiates a new DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder and sets the default values.
func NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) {
    m := &DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/deviceLocalCredentials/{deviceLocalCredentialInfo%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder instantiates a new DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder and sets the default values.
func NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property deviceLocalCredentials for directory
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
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
// Get retrieve the properties of a deviceLocalCredentialInfo for a specified device object. 
// returns a DeviceLocalCredentialInfoable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/devicelocalcredentialinfo-get?view=graph-rest-1.0
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) Get(ctx context.Context, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceLocalCredentialInfoFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable), nil
}
// Patch update the navigation property deviceLocalCredentials in directory
// returns a DeviceLocalCredentialInfoable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateDeviceLocalCredentialInfoFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable), nil
}
// ToDeleteRequestInformation delete navigation property deviceLocalCredentials for directory
// returns a *RequestInformation when successful
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the properties of a deviceLocalCredentialInfo for a specified device object. 
// returns a *RequestInformation when successful
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPatchRequestInformation update the navigation property deviceLocalCredentials in directory
// returns a *RequestInformation when successful
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.DeviceLocalCredentialInfoable, requestConfiguration *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder when successful
func (m *DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) WithUrl(rawUrl string)(*DeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder) {
    return NewDeviceLocalCredentialsDeviceLocalCredentialInfoItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
