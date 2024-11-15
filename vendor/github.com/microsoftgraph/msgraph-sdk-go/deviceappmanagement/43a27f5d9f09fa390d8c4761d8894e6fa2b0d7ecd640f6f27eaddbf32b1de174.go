package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder provides operations to call the renewUpload method.
type MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderInternal instantiates a new MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder and sets the default values.
func NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    m := &MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.managedAndroidLobApp/contentVersions/{mobileAppContent%2Did}/files/{mobileAppContentFile%2Did}/renewUpload", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder instantiates a new MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder and sets the default values.
func NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderInternal(urlParams, requestAdapter)
}
// Post renews the SAS URI for an application file upload.
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) Post(ctx context.Context, requestConfiguration *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration)(error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, requestConfiguration);
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
// ToPostRequestInformation renews the SAS URI for an application file upload.
// returns a *RequestInformation when successful
func (m *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder when successful
func (m *MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    return NewMobileAppsItemGraphManagedAndroidLobAppContentVersionsItemFilesItemRenewUploadRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
