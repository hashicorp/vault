package drives

import (
    "context"
    i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274 "strconv"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder provides operations to call the usedRange method.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/usedRange(valuesOnly={valuesOnly})", pathParameters),
    }
    if valuesOnly != nil {
        m.BaseRequestBuilder.PathParameters["valuesOnly"] = i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274.FormatBool(*valuesOnly)
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function usedRange
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlySortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlySortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlySortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function usedRange
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemUsedRangeWithValuesOnlyWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
