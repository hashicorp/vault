package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder provides operations to call the lastCell method.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/usedRange(valuesOnly={valuesOnly})/lastCell()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderInternal(urlParams, requestAdapter)
}
// Get invoke function lastCell
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookRangeFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable), nil
}
// ToGetRequestInformation invoke function lastCell
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
