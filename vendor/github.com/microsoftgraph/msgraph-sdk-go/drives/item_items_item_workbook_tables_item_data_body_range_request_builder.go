package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder provides operations to call the dataBodyRange method.
type ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemDataBodyRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/dataBodyRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemDataBodyRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemDataBodyRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemDataBodyRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemDataBodyRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with the data body of the table.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/table-databodyrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemDataBodyRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemDataBodyRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemDataBodyRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemDataBodyRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemDataBodyRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemDataBodyRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with the data body of the table.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemDataBodyRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemDataBodyRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemDataBodyRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemDataBodyRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemDataBodyRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemDataBodyRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
