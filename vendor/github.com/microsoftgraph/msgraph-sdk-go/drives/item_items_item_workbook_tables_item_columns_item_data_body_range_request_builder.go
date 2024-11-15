package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder provides operations to call the dataBodyRange method.
type ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/columns/{workbookTableColumn%2Did}/dataBodyRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with the data body of the column.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/tablecolumn-databodyrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with the data body of the column.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemColumnsItemDataBodyRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
