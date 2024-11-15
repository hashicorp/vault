package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder provides operations to call the rowsAbove method.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/usedRange(valuesOnly={valuesOnly})/rowsAbove()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderInternal(urlParams, requestAdapter)
}
// Get invoke function rowsAbove
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// ToGetRequestInformation invoke function rowsAbove
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
