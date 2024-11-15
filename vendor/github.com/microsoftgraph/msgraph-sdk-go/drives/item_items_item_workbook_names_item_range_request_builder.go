package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookNamesItemRangeRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookNamesItemRangeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookNamesItemRangeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookNamesItemRangeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookNamesItemRangeBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookNamesItemRangeBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookNamesItemRangeCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookNamesItemRangeCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookNamesItemRangeClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Clear()(*ItemItemsItemWorkbookNamesItemRangeClearRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookNamesItemRangeColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookNamesItemRangeColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookNamesItemRangeColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookNamesItemRangeColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookNamesItemRangeColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookNamesItemRangeColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookNamesItemRangeColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookNamesItemRangeColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookNamesItemRangeColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookNamesItemRangeColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookNamesItemRangeRequestBuilderInternal instantiates a new ItemItemsItemWorkbookNamesItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookNamesItemRangeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookNamesItemRangeRequestBuilder) {
    m := &ItemItemsItemWorkbookNamesItemRangeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/names/{workbookNamedItem%2Did}/range()", pathParameters),
    }
    return m
}
// NewItemItemsItemWorkbookNamesItemRangeRequestBuilder instantiates a new ItemItemsItemWorkbookNamesItemRangeRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookNamesItemRangeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookNamesItemRangeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookNamesItemRangeRequestBuilderInternal(urlParams, requestAdapter)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookNamesItemRangeDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) DeletePath()(*ItemItemsItemWorkbookNamesItemRangeDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookNamesItemRangeEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookNamesItemRangeEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookNamesItemRangeEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) EntireRow()(*ItemItemsItemWorkbookNamesItemRangeEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookNamesItemRangeFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Format()(*ItemItemsItemWorkbookNamesItemRangeFormatRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get returns the range object that is associated with the name. Throws an exception if the named item's type is not a range.
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/nameditem-range?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookNamesItemRangeRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookNamesItemRangeInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Insert()(*ItemItemsItemWorkbookNamesItemRangeInsertRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookNamesItemRangeIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookNamesItemRangeIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookNamesItemRangeLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) LastCell()(*ItemItemsItemWorkbookNamesItemRangeLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookNamesItemRangeLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) LastColumn()(*ItemItemsItemWorkbookNamesItemRangeLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookNamesItemRangeLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) LastRow()(*ItemItemsItemWorkbookNamesItemRangeLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookNamesItemRangeMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Merge()(*ItemItemsItemWorkbookNamesItemRangeMergeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookNamesItemRangeRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookNamesItemRangeRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookNamesItemRangeRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookNamesItemRangeRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookNamesItemRangeRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookNamesItemRangeRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookNamesItemRangeRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookNamesItemRangeRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookNamesItemRangeRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookNamesItemRangeRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookNamesItemRangeSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Sort()(*ItemItemsItemWorkbookNamesItemRangeSortRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation returns the range object that is associated with the name. Throws an exception if the named item's type is not a range.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookNamesItemRangeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookNamesItemRangeUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Unmerge()(*ItemItemsItemWorkbookNamesItemRangeUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookNamesItemRangeUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) UsedRange()(*ItemItemsItemWorkbookNamesItemRangeUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookNamesItemRangeUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookNamesItemRangeUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookNamesItemRangeVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) VisibleView()(*ItemItemsItemWorkbookNamesItemRangeVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookNamesItemRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookNamesItemRangeRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookNamesItemRangeWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookNamesItemRangeRequestBuilder) Worksheet()(*ItemItemsItemWorkbookNamesItemRangeWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookNamesItemRangeWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
