package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder provides operations to call the commit method.
type MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderInternal instantiates a new MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) {
    m := &MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.windowsMobileMSI/contentVersions/{mobileAppContent%2Did}/files/{mobileAppContentFile%2Did}/commit", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder instantiates a new MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderInternal(urlParams, requestAdapter)
}
// Post commits a file of a given app.
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) Post(ctx context.Context, body MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitPostRequestBodyable, requestConfiguration *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation commits a file of a given app.
// returns a *RequestInformation when successful
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) ToPostRequestInformation(ctx context.Context, body MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitPostRequestBodyable, requestConfiguration *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder when successful
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder) {
    return NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemCommitRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
