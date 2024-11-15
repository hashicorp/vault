package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder provides operations to call the range method.
type ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// BoundingRectWithAnotherRange provides operations to call the boundingRect method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressBoundingRectWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) BoundingRectWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressBoundingRectWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressBoundingRectWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// CellWithRowWithColumn provides operations to call the cell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressCellWithRowWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) CellWithRowWithColumn(column *int32, row *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressCellWithRowWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressCellWithRowWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column, row)
}
// Clear provides operations to call the clear method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressClearRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Clear()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressClearRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressClearRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfter provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ColumnsAfter()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsAfterWithCount provides operations to call the columnsAfter method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ColumnsAfterWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsAfterWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnsBefore provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ColumnsBefore()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ColumnsBeforeWithCount provides operations to call the columnsBefore method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ColumnsBeforeWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnsBeforeWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// ColumnWithColumn provides operations to call the column method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnWithColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ColumnWithColumn(column *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnWithColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressColumnWithColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, column)
}
// NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderInternal instantiates a new ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, address *string)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) {
    m := &ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/worksheets/{workbookWorksheet%2Did}/range(address='{address}')", pathParameters),
    }
    if address != nil {
        m.BaseRequestBuilder.PathParameters["address"] = *address
    }
    return m
}
// NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder instantiates a new ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// DeletePath provides operations to call the delete method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressDeleteRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) DeletePath()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressDeleteRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressDeleteRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireColumn provides operations to call the entireColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) EntireColumn()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// EntireRow provides operations to call the entireRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) EntireRow()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressEntireRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Format provides operations to manage the format property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressFormatRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Format()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressFormatRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressFormatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function range
// returns a WorkbookRangeable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookRangeable, error) {
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
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressInsertRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Insert()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressInsertRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressInsertRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IntersectionWithAnotherRange provides operations to call the intersection method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressIntersectionWithAnotherRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) IntersectionWithAnotherRange(anotherRange *string)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressIntersectionWithAnotherRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressIntersectionWithAnotherRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, anotherRange)
}
// LastCell provides operations to call the lastCell method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastCellRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) LastCell()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastCellRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastCellRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastColumn provides operations to call the lastColumn method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastColumnRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) LastColumn()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastColumnRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastColumnRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// LastRow provides operations to call the lastRow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) LastRow()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressLastRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Merge provides operations to call the merge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressMergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Merge()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressMergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressMergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// OffsetRangeWithRowOffsetWithColumnOffset provides operations to call the offsetRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) OffsetRangeWithRowOffsetWithColumnOffset(columnOffset *int32, rowOffset *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressOffsetRangeWithRowOffsetWithColumnOffsetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, columnOffset, rowOffset)
}
// ResizedRangeWithDeltaRowsWithDeltaColumns provides operations to call the resizedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ResizedRangeWithDeltaRowsWithDeltaColumns(deltaColumns *int32, deltaRows *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressResizedRangeWithDeltaRowsWithDeltaColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, deltaColumns, deltaRows)
}
// RowsAbove provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) RowsAbove()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsAboveWithCount provides operations to call the rowsAbove method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) RowsAboveWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsAboveWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowsBelow provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) RowsBelow()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RowsBelowWithCount provides operations to call the rowsBelow method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowWithCountRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) RowsBelowWithCount(count *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowWithCountRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowsBelowWithCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, count)
}
// RowWithRow provides operations to call the row method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowWithRowRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) RowWithRow(row *int32)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowWithRowRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRowWithRowRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, row)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Sort()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressSortRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function range
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// Unmerge provides operations to call the unmerge method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUnmergeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Unmerge()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUnmergeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressUnmergeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRange provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) UsedRange()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// UsedRangeWithValuesOnly provides operations to call the usedRange method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeWithValuesOnlyRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) UsedRangeWithValuesOnly(valuesOnly *bool)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeWithValuesOnlyRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressUsedRangeWithValuesOnlyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter, valuesOnly)
}
// VisibleView provides operations to call the visibleView method.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressVisibleViewRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) VisibleView()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressVisibleViewRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressVisibleViewRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookRange entity.
// returns a *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookWorksheetsItemRangeWithAddressRequestBuilder) Worksheet()(*ItemItemsItemWorkbookWorksheetsItemRangeWithAddressWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookWorksheetsItemRangeWithAddressWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
