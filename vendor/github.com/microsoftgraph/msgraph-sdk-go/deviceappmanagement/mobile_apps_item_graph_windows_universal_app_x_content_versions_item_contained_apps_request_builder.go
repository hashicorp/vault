package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder provides operations to manage the containedApps property of the microsoft.graph.mobileAppContent entity.
type MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetQueryParameters the collection of contained apps in a MobileLobApp acting as a package.
type MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetQueryParameters
}
// MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByMobileContainedAppId provides operations to manage the containedApps property of the microsoft.graph.mobileAppContent entity.
// returns a *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsMobileContainedAppItemRequestBuilder when successful
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) ByMobileContainedAppId(mobileContainedAppId string)(*MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsMobileContainedAppItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if mobileContainedAppId != "" {
        urlTplParams["mobileContainedApp%2Did"] = mobileContainedAppId
    }
    return NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsMobileContainedAppItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderInternal instantiates a new MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) {
    m := &MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.windowsUniversalAppX/contentVersions/{mobileAppContent%2Did}/containedApps{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder instantiates a new MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsCountRequestBuilder when successful
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) Count()(*MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsCountRequestBuilder) {
    return NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the collection of contained apps in a MobileLobApp acting as a package.
// returns a MobileContainedAppCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) Get(ctx context.Context, requestConfiguration *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMobileContainedAppCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppCollectionResponseable), nil
}
// Post create new navigation property to containedApps for deviceAppManagement
// returns a MobileContainedAppable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppable, requestConfiguration *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMobileContainedAppFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppable), nil
}
// ToGetRequestInformation the collection of contained apps in a MobileLobApp acting as a package.
// returns a *RequestInformation when successful
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to containedApps for deviceAppManagement
// returns a *RequestInformation when successful
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileContainedAppable, requestConfiguration *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder when successful
func (m *MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder) {
    return NewMobileAppsItemGraphWindowsUniversalAppXContentVersionsItemContainedAppsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
