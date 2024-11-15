package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder provides operations to call the protect method.
type ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/protection/protect", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderInternal(urlParams, requestAdapter)
}
// Post protect a worksheet. It throws if the worksheet has been protected.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/worksheetprotection-protect?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) Post(ctx context.Context, body ItemItemsItemWorkbookWorksheetsItemProtectionProtectPostRequestBodyable, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation protect a worksheet. It throws if the worksheet has been protected.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemItemsItemWorkbookWorksheetsItemProtectionProtectPostRequestBodyable, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemProtectionProtectRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
