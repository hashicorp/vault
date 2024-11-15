package storage

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder provides operations to manage the services property of the microsoft.graph.unifiedStorageQuota entity.
type SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetQueryParameters get services from storage
type SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetQueryParameters
}
// SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderInternal instantiates a new SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder and sets the default values.
func NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) {
    m := &SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/storage/settings/quota/services/{serviceStorageQuotaBreakdown%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder instantiates a new SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder and sets the default values.
func NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property services for storage
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get get services from storage
// returns a ServiceStorageQuotaBreakdownable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) Get(ctx context.Context, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateServiceStorageQuotaBreakdownFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable), nil
}
// Patch update the navigation property services in storage
// returns a ServiceStorageQuotaBreakdownable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateServiceStorageQuotaBreakdownFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable), nil
}
// ToDeleteRequestInformation delete navigation property services for storage
// returns a *RequestInformation when successful
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation get services from storage
// returns a *RequestInformation when successful
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property services in storage
// returns a *RequestInformation when successful
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ServiceStorageQuotaBreakdownable, requestConfiguration *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder when successful
func (m *SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) WithUrl(rawUrl string)(*SettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder) {
    return NewSettingsQuotaServicesServiceStorageQuotaBreakdownItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
