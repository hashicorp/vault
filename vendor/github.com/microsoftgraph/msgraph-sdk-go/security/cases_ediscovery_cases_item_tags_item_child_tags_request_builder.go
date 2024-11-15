package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder provides operations to manage the childTags property of the microsoft.graph.security.ediscoveryReviewTag entity.
type CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetQueryParameters returns the tags that are a child of a tag.
type CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetQueryParameters struct {
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
// CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetQueryParameters
}
// ByEdiscoveryReviewTagId1 provides operations to manage the childTags property of the microsoft.graph.security.ediscoveryReviewTag entity.
// returns a *CasesEdiscoveryCasesItemTagsItemChildTagsEdiscoveryReviewTagItemRequestBuilder when successful
func (m *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) ByEdiscoveryReviewTagId1(ediscoveryReviewTagId1 string)(*CasesEdiscoveryCasesItemTagsItemChildTagsEdiscoveryReviewTagItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if ediscoveryReviewTagId1 != "" {
        urlTplParams["ediscoveryReviewTag%2Did1"] = ediscoveryReviewTagId1
    }
    return NewCasesEdiscoveryCasesItemTagsItemChildTagsEdiscoveryReviewTagItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderInternal instantiates a new CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) {
    m := &CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/cases/ediscoveryCases/{ediscoveryCase%2Did}/tags/{ediscoveryReviewTag%2Did}/childTags{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder instantiates a new CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder and sets the default values.
func NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *CasesEdiscoveryCasesItemTagsItemChildTagsCountRequestBuilder when successful
func (m *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) Count()(*CasesEdiscoveryCasesItemTagsItemChildTagsCountRequestBuilder) {
    return NewCasesEdiscoveryCasesItemTagsItemChildTagsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get returns the tags that are a child of a tag.
// returns a EdiscoveryReviewTagCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) Get(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryReviewTagCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateEdiscoveryReviewTagCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.EdiscoveryReviewTagCollectionResponseable), nil
}
// ToGetRequestInformation returns the tags that are a child of a tag.
// returns a *RequestInformation when successful
func (m *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder when successful
func (m *CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) WithUrl(rawUrl string)(*CasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder) {
    return NewCasesEdiscoveryCasesItemTagsItemChildTagsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
