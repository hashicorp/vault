package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/names/{workbookNamedItem%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get returns the range object that is associated with the name. Throws an exception if the named item's type is not a range.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/nameditem-range?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation returns the range object that is associated with the name. Throws an exception if the named item's type is not a range.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemNamesItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemNamesItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemNamesItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
