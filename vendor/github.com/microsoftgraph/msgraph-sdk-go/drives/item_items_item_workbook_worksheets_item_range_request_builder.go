package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function range
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// Insert provides operations to call the insert method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function range
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
