package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder provides operations to call the headerRowRange method.
type ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/headerRowRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with header row of the table.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/table-headerrowrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with header row of the table.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemHeaderRowRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemHeaderRowRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemHeaderRowRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemHeaderRowRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
