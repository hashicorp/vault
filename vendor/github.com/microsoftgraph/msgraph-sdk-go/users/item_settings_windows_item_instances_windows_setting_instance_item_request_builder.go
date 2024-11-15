package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder provides operations to manage the instances property of the microsoft.graph.windowsSetting entity.
type ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetQueryParameters a collection of setting values for a given windowsSetting.
type ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetQueryParameters
}
// ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderInternal instantiates a new ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder and sets the default values.
func NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) {
    m := &ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/settings/windows/{windowsSetting%2Did}/instances/{windowsSettingInstance%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder instantiates a new ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder and sets the default values.
func NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property instances for users
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get a collection of setting values for a given windowsSetting.
// returns a WindowsSettingInstanceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsSettingInstanceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable), nil
}
// Patch update the navigation property instances in users
// returns a WindowsSettingInstanceable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsSettingInstanceFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable), nil
}
// ToDeleteRequestInformation delete navigation property instances for users
// returns a *RequestInformation when successful
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation a collection of setting values for a given windowsSetting.
// returns a *RequestInformation when successful
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property instances in users
// returns a *RequestInformation when successful
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsSettingInstanceable, requestConfiguration *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder when successful
func (m *ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) WithUrl(rawUrl string)(*ItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder) {
    return NewItemSettingsWindowsItemInstancesWindowsSettingInstanceItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
