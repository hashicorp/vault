package drives

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder provides operations to call the tableRowOperationResult method.
type ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderInternal instantiates a new ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, key *string)(*ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) {
    m := &ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/drives/{drive%2Did}/items/{driveItem%2Did}/workbook/tableRowOperationResult(key='{key}')", pathParameters),
    }
    if key != nil {
        m.BaseRequestBuilder.PathParameters["key"] = *key
    }
    return m
}
// NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder instantiates a new ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder and sets the default values.
func NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// Get this function is the last in a series of steps to create a workbookTableRow resource asynchronously. A best practice to create multiple table rows is to batch them in one create tableRow operation and carry out the operation asynchronously. An asynchronous request to create table rows involves the following steps:1. Issue an async Create tableRow request and get the query URL returned in the Location response header.2. Use the query URL returned from step 1 to issue the Get workbookOperation request and get the operation ID for step 3.     Alternatively, for convenience, after you get a succeeded operationStatus result, you can get the query URL from the resourceLocation property of the workbookOperation returned in the response, and apply the query URL to step 3. 3. Use the query URL returned from step 2 as the GET request URL for this function tableRowOperationResult. A successful function call returns the new table rows in a workbookTableRow resource. This function does not do anything if called independently.
// returns a WorkbookTableRowable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/workbook-tablerowoperationresult?view=graph-rest-1.0
func (m *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookTableRowable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWorkbookTableRowFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WorkbookTableRowable), nil
}
// ToGetRequestInformation this function is the last in a series of steps to create a workbookTableRow resource asynchronously. A best practice to create multiple table rows is to batch them in one create tableRow operation and carry out the operation asynchronously. An asynchronous request to create table rows involves the following steps:1. Issue an async Create tableRow request and get the query URL returned in the Location response header.2. Use the query URL returned from step 1 to issue the Get workbookOperation request and get the operation ID for step 3.     Alternatively, for convenience, after you get a succeeded operationStatus result, you can get the query URL from the resourceLocation property of the workbookOperation returned in the response, and apply the query URL to step 3. 3. Use the query URL returned from step 2 as the GET request URL for this function tableRowOperationResult. A successful function call returns the new table rows in a workbookTableRow resource. This function does not do anything if called independently.
// returns a *RequestInformation when successful
func (m *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder when successful
func (m *ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) WithUrl(rawUrl string)(*ItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder) {
    return NewItemItemsItemWorkbookTableRowOperationResultWithKeyRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
