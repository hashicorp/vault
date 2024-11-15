package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder provides operations to call the headerRowRange method.
type ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/columns/{workbookTableColumn%2Did}/headerRowRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with the header row of the column.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/tablecolumn-headerrowrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with the header row of the column.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemHeaderRowRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
