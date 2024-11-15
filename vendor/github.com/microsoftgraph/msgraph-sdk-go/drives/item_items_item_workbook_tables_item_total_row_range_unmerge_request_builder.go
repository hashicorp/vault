package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder provides operations to call the unmerge method.
type ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/totalRowRange()/unmerge", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderInternal(urlParams, requestAdapter)
}
// Post invoke action unmerge
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderPostRequestConfiguration)(error) {
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
// ToPostRequestInformation invoke action unmerge
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) ToPostRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
