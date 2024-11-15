package drives

import (
    "context"
    i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274 "strconv"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder provides operations to call the itemAt method.
type ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ClearFilters provides operations to call the clearFilters method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexClearFiltersRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) ClearFilters()(*ItemItemsItemWorkbookTablesItemAtWithIndexClearFiltersRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexClearFiltersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Columns provides operations to manage the columns property of the microsoft.graph.workbookTable entity.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexColumnsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) Columns()(*ItemItemsItemWorkbookTablesItemAtWithIndexColumnsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexColumnsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, index *int32)(*ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) {
    m := &ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tables/itemAt(index={index})", pathParameters),
    }
    if index != nil {
        m.BaseRequestBuilder.PathParameters["index"] = i53ac87e8cb3cc9276228f74d38694a208cacb99bb8ceb705eeae99fb88d4d274.FormatInt(int64(*index), 10)
    }
    return m
}
// NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder instantiates a new ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// ConvertToRange provides operations to call the convertToRange method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexConvertToRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) ConvertToRange()(*ItemItemsItemWorkbookTablesItemAtWithIndexConvertToRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexConvertToRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DataBodyRange provides operations to call the dataBodyRange method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexDataBodyRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) DataBodyRange()(*ItemItemsItemWorkbookTablesItemAtWithIndexDataBodyRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexDataBodyRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get invoke function itemAt
// returns a WorkbookTableable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookTableable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookTableFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookTableable), nil
}
// HeaderRowRange provides operations to call the headerRowRange method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexHeaderRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) HeaderRowRange()(*ItemItemsItemWorkbookTablesItemAtWithIndexHeaderRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexHeaderRowRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RangeEscaped provides operations to call the range method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) RangeEscaped()(*ItemItemsItemWorkbookTablesItemAtWithIndexRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ReapplyFilters provides operations to call the reapplyFilters method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexReapplyFiltersRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) ReapplyFilters()(*ItemItemsItemWorkbookTablesItemAtWithIndexReapplyFiltersRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexReapplyFiltersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Rows provides operations to manage the rows property of the microsoft.graph.workbookTable entity.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexRowsRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) Rows()(*ItemItemsItemWorkbookTablesItemAtWithIndexRowsRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexRowsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Sort provides operations to manage the sort property of the microsoft.graph.workbookTable entity.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexSortRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) Sort()(*ItemItemsItemWorkbookTablesItemAtWithIndexSortRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexSortRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation invoke function itemAt
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// TotalRowRange provides operations to call the totalRowRange method.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexTotalRowRangeRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) TotalRowRange()(*ItemItemsItemWorkbookTablesItemAtWithIndexTotalRowRangeRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexTotalRowRangeRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
// Worksheet provides operations to manage the worksheet property of the microsoft.graph.workbookTable entity.
// returns a *ItemItemsItemWorkbookTablesItemAtWithIndexWorksheetRequestBuilder when successful
func (m *ItemItemsItemWorkbookTablesItemAtWithIndexRequestBuilder) Worksheet()(*ItemItemsItemWorkbookTablesItemAtWithIndexWorksheetRequestBuilder) {
    return NewItemItemsItemWorkbookTablesItemAtWithIndexWorksheetRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
