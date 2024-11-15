package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder provides operations to manage the categories property of the microsoft.graph.mobileApp entity.
type MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetQueryParameters the list of categories for this app.
type MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetQueryParameters
}
// ByMobileAppCategoryId provides operations to manage the categories property of the microsoft.graph.mobileApp entity.
// returns a *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesMobileAppCategoryItemRequestBuilder when successful
func (m *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) ByMobileAppCategoryId(mobileAppCategoryId string)(*MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesMobileAppCategoryItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if mobileAppCategoryId != "" {
        urlTplParams["mobileAppCategory%2Did"] = mobileAppCategoryId
    }
    return NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesMobileAppCategoryItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderInternal instantiates a new MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder and sets the default values.
func NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) {
    m := &MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mobileApps/{mobileApp%2Did}/graph.microsoftStoreForBusinessApp/categories{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder instantiates a new MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder and sets the default values.
func NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesCountRequestBuilder when successful
func (m *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) Count()(*MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesCountRequestBuilder) {
    return NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get the list of categories for this app.
// returns a MobileAppCategoryCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) Get(ctx context.Context, requestConfiguration *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileAppCategoryCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateMobileAppCategoryCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.MobileAppCategoryCollectionResponseable), nil
}
// ToGetRequestInformation the list of categories for this app.
// returns a *RequestInformation when successful
func (m *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder when successful
func (m *MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) WithUrl(rawUrl string)(*MobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder) {
    return NewMobileAppsItemGraphMicrosoftStoreForBusinessAppCategoriesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
