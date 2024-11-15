package reports

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// GetOffice365ActivationsUserDetailRequestBuilder provides operations to call the getOffice365ActivationsUserDetail method.
type GetOffice365ActivationsUserDetailRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// GetOffice365ActivationsUserDetailRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type GetOffice365ActivationsUserDetailRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewGetOffice365ActivationsUserDetailRequestBuilderInternal instantiates a new GetOffice365ActivationsUserDetailRequestBuilder and sets the default values.
func NewGetOffice365ActivationsUserDetailRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*GetOffice365ActivationsUserDetailRequestBuilder) {
    m := &GetOffice365ActivationsUserDetailRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/reports/getOffice365ActivationsUserDetail()", pathParameters),
    }
    return m
}
// NewGetOffice365ActivationsUserDetailRequestBuilder instantiates a new GetOffice365ActivationsUserDetailRequestBuilder and sets the default values.
func NewGetOffice365ActivationsUserDetailRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*GetOffice365ActivationsUserDetailRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewGetOffice365ActivationsUserDetailRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get details about users who have activated Microsoft 365.
// returns a []byte when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/reportroot-getoffice365activationsuserdetail?view=graph-rest-1.0
func (m *GetOffice365ActivationsUserDetailRequestBuilder) Get(ctx context.Context, requestConfiguration *GetOffice365ActivationsUserDetailRequestBuilderGetRequestConfiguration)([]byte, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.SendPrimitive(ctx, requestInfo, "[]byte", errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.([]byte), nil
}
// ToGetRequestInformation get details about users who have activated Microsoft 365.
// returns a *RequestInformation when successful
func (m *GetOffice365ActivationsUserDetailRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *GetOffice365ActivationsUserDetailRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/octet-stream, application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *GetOffice365ActivationsUserDetailRequestBuilder when successful
func (m *GetOffice365ActivationsUserDetailRequestBuilder) WithUrl(rawUrl string)(*GetOffice365ActivationsUserDetailRequestBuilder) {
    return NewGetOffice365ActivationsUserDetailRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
