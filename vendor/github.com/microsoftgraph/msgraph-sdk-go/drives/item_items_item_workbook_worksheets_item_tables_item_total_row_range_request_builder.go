package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder provides operations to call the totalRowRange method.
type ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/tables/{workbookTable%2Did}/totalRowRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with totals row of the table.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/table-totalrowrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with totals row of the table.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemTotalRowRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
