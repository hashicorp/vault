package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder provides operations to call the headerRowRange method.
type ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/tables/{workbookTable%2Did}/headerRowRange()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with header row of the table.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/table-headerrowrange?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with header row of the table.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemHeaderRowRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
