package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder provides operations to manage the subcategories property of the microsoft.graph.security.categoryTemplate entity.
type LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetQueryParameters read the properties and relationships of a subcategoryTemplate object.
type LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetQueryParameters
}
// LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderInternal instantiates a new LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder and sets the default values.
func NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) {
    m := &LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/labels/categories/{categoryTemplate%2Did}/subcategories/{subcategoryTemplate%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder instantiates a new LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder and sets the default values.
func NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property subcategories for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.ToDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// Get read the properties and relationships of a subcategoryTemplate object.
// returns a SubcategoryTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/security-subcategorytemplate-get?view=graph-rest-1.0
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) Get(ctx context.Context, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateSubcategoryTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable), nil
}
// Patch update the navigation property subcategories in security
// returns a SubcategoryTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateSubcategoryTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable), nil
}
// ToDeleteRequestInformation delete navigation property subcategories for security
// returns a *RequestInformation when successful
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a subcategoryTemplate object.
// returns a *RequestInformation when successful
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property subcategories in security
// returns a *RequestInformation when successful
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.SubcategoryTemplateable, requestConfiguration *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder when successful
func (m *LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) WithUrl(rawUrl string)(*LabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder) {
    return NewLabelsCategoriesItemSubcategoriesSubcategoryTemplateItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
