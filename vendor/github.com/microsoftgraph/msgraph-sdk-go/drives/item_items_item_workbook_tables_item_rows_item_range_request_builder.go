package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemRowsItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookTablesItemRowsItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/{workbookTable%2Did}/rows/{workbookTableRow%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookTablesItemRowsItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookTablesItemRowsItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookTablesItemRowsItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookTablesItemRowsItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function range
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookTablesItemRowsItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookTablesItemRowsItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookTablesItemRowsItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookTablesItemRowsItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookTablesItemRowsItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemRowsItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function range
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookTablesItemRowsItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookTablesItemRowsItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookTablesItemRowsItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemRowsItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemRowsItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemRowsItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
