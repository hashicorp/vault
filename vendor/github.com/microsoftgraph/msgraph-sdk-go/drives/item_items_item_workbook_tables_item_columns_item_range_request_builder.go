package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/columns/{workbookTableColumn%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with the entire column.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/tablecolumn-range?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with the entire column.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemColumnsItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
