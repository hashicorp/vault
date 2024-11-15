package drives

import (
    "context"
    i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274 "strconv"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder provides operations to call the cell method.
type ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn1 provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnWithColumn1RequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ColumnWithColumn1(column1 *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnWithColumn1RequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnColumnWithColumn1RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column1)
}
// NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/cell(row={row},column={column})", pathParameters),
    }
    if column != nil {
        m.BaseRequestBuilder.PathParameters["column"] = i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274.FormatInt(int64(*column), 10)
    }
    if row != nil {
        m.BaseRequestBuilder.PathParameters["row"] = i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274.FormatInt(int64(*row), 10)
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderInternal(urlParams, requestAdapter, nil, nil)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get gets the range object containing the single cell based on row and column numbers. The cell can be outside the bounds of its parent range, so long as it's stays within the worksheet grid.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/worksheet-cell?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow1 provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowWithRow1RequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) RowWithRow1(row1 *int32)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowWithRow1RequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRowWithRow1RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row1)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation gets the range object containing the single cell based on row and column numbers. The cell can be outside the bounds of its parent range, so long as it's stays within the worksheet grid.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemCellWithRowWithColumnWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
