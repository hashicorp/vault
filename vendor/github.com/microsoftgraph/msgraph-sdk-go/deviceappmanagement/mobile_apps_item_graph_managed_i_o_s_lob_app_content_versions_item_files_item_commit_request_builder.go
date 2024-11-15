package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder provides operations to call the commit method.
type MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderInternal instantiates a new MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder and sets the default values.
func NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) {
    m := &MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.managedIOSLobApp/contentVersions/{mobileAppContent%2Did}/files/{mobileAppContentFile%2Did}/commit", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder instantiates a new MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder and sets the default values.
func NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderInternal(urlParams, requestAdapter)
}
// Post commits a file of a given app.
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) Post(ctx context.Context, body MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitPostRequestBodyable, requestConfiguration *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration)(error) {
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
func (m *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) ToPostRequestInformation(ctx context.Context, body MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitPostRequestBodyable, requestConfiguration *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder when successful
func (m *MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder) {
    return NewMobileAppsItemGraphManagedIOSLobAppContentVersionsItemFilesItemCommitRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
