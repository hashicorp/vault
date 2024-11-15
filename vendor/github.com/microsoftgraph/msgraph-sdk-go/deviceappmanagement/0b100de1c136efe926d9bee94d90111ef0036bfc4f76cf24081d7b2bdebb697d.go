package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder provides operations to call the renewUpload method.
type MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderInternal instantiates a new MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    m := &MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.windowsMobileMSI/contentVersions/{mobileAppContent%2Did}/files/{mobileAppContentFile%2Did}/renewUpload", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder instantiates a new MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder and sets the default values.
func NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderInternal(urlParams, requestAdapter)
}
// Post renews the SAS URI for an application file upload.
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) Post(ctx context.Context, requestConfiguration *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration)(error) {
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
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder when successful
func (m *MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder) {
    return NewMobileAppsItemGraphWindowsMobileMSIContentVersionsItemFilesItemRenewUploadRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
