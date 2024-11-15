package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder provides operations to call the totalRowRange method.
type ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemTotalRowRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/totalRowRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemTotalRowRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemTotalRowRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemTotalRowRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemTotalRowRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with totals row of the table.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/table-totalrowrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemTotalRowRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemTotalRowRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemTotalRowRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemTotalRowRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemTotalRowRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemTotalRowRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with totals row of the table.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemTotalRowRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemTotalRowRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemTotalRowRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemTotalRowRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemTotalRowRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
