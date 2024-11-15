package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/tables/{workbookTable%2Did}/columns/{workbookTableColumn%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object associated with the entire column.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/tablecolumn-range?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object associated with the entire column.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemTablesItemColumnsItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
