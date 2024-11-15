package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder provides operations to call the merge method.
type ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/rows/{workbookTableRow%2Did}/range()/merge", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderInternal(urlParams, requestAdapter)
}
// Post invoke action merge
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) Post(ctx context.Context, body ItemItemsItemWorkbookTablesItemRowsItemRangeMergePostRequestBodyable, requestConfiguration *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation invoke action merge
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemItemsItemWorkbookTablesItemRowsItemRangeMergePostRequestBodyable, requestConfiguration *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
