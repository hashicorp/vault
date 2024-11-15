package communications

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CallsLogTeleconferenceDeviceQualityRequestBuilder provides operations to call the logTeleconferenceDeviceQuality method.
type CallsLogTeleconferenceDeviceQualityRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CallsLogTeleconferenceDeviceQualityRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CallsLogTeleconferenceDeviceQualityRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewCallsLogTeleconferenceDeviceQualityRequestBuilderInternal instantiates a new CallsLogTeleconferenceDeviceQualityRequestBuilder and sets the default values.
func NewCallsLogTeleconferenceDeviceQualityRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CallsLogTeleconferenceDeviceQualityRequestBuilder) {
    m := &CallsLogTeleconferenceDeviceQualityRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/communications/calls/logTeleconferenceDeviceQuality", pathParameters),
    }
    return m
}
// NewCallsLogTeleconferenceDeviceQualityRequestBuilder instantiates a new CallsLogTeleconferenceDeviceQualityRequestBuilder and sets the default values.
func NewCallsLogTeleconferenceDeviceQualityRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CallsLogTeleconferenceDeviceQualityRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCallsLogTeleconferenceDeviceQualityRequestBuilderInternal(urlParams, requestAdapter)
}
// Post log video teleconferencing device quality data. The Cloud Video Interop (CVI) bot represents video teleconferencing (VTC) devices and acts as a back-to-back agent for a VTC device in a conference call. Because a CVI bot is in the middle of the VTC and Microsoft Teams infrastructure as a VTC proxy, it has two media legs. One media leg is between the CVI bot and Teams infrastructure, such as Teams conference server or a Teams client. The other media leg is between the CVI bot and the VTC device.  The third-party partners own the VTC media leg and the Teams infrastructure cannot access the quality data of the third-party call leg.  This method is only for the CVI partners to provide their media quality data.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/call-logteleconferencedevicequality?view=graph-rest-1.0
func (m *CallsLogTeleconferenceDeviceQualityRequestBuilder) Post(ctx context.Context, body CallsLogTeleconferenceDeviceQualityPostRequestBodyable, requestConfiguration *CallsLogTeleconferenceDeviceQualityRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation log video teleconferencing device quality data. The Cloud Video Interop (CVI) bot represents video teleconferencing (VTC) devices and acts as a back-to-back agent for a VTC device in a conference call. Because a CVI bot is in the middle of the VTC and Microsoft Teams infrastructure as a VTC proxy, it has two media legs. One media leg is between the CVI bot and Teams infrastructure, such as Teams conference server or a Teams client. The other media leg is between the CVI bot and the VTC device.  The third-party partners own the VTC media leg and the Teams infrastructure cannot access the quality data of the third-party call leg.  This method is only for the CVI partners to provide their media quality data.
// returns a *RequestInformation when successful
func (m *CallsLogTeleconferenceDeviceQualityRequestBuilder) ToPostRequestInformation(ctx context.Context, body CallsLogTeleconferenceDeviceQualityPostRequestBodyable, requestConfiguration *CallsLogTeleconferenceDeviceQualityRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CallsLogTeleconferenceDeviceQualityRequestBuilder when successful
func (m *CallsLogTeleconferenceDeviceQualityRequestBuilder) WithUrl(rawUrl string)(*CallsLogTeleconferenceDeviceQualityRequestBuilder) {
    return NewCallsLogTeleconferenceDeviceQualityRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
