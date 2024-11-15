package directory

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder provides operations to manage the customSecurityAttributeDefinitions property of the microsoft.graph.directory entity.
type CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetQueryParameters read the properties and relationships of a customSecurityAttributeDefinition object.
type CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetQueryParameters
}
// CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// AllowedValues provides operations to manage the allowedValues property of the microsoft.graph.customSecurityAttributeDefinition entity.
// returns a *CustomSecurityAttributeDefinitionsItemAllowedValuesRequestBuilder when successful
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) AllowedValues()(*CustomSecurityAttributeDefinitionsItemAllowedValuesRequestBuilder) {
    return NewCustomSecurityAttributeDefinitionsItemAllowedValuesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderInternal instantiates a new CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder and sets the default values.
func NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) {
    m := &CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/directory/customSecurityAttributeDefinitions/{customSecurityAttributeDefinition%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder instantiates a new CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder and sets the default values.
func NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property customSecurityAttributeDefinitions for directory
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get read the properties and relationships of a customSecurityAttributeDefinition object.
// returns a CustomSecurityAttributeDefinitionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/customsecurityattributedefinition-get?view=graph-rest-1.0
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) Get(ctx context.Context, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomSecurityAttributeDefinitionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable), nil
}
// Patch update the properties of a customSecurityAttributeDefinition object.
// returns a CustomSecurityAttributeDefinitionable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/customsecurityattributedefinition-update?view=graph-rest-1.0
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateCustomSecurityAttributeDefinitionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable), nil
}
// ToDeleteRequestInformation delete navigation property customSecurityAttributeDefinitions for directory
// returns a *RequestInformation when successful
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation read the properties and relationships of a customSecurityAttributeDefinition object.
// returns a *RequestInformation when successful
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the properties of a customSecurityAttributeDefinition object.
// returns a *RequestInformation when successful
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CustomSecurityAttributeDefinitionable, requestConfiguration *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder when successful
func (m *CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) WithUrl(rawUrl string)(*CustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder) {
    return NewCustomSecurityAttributeDefinitionsCustomSecurityAttributeDefinitionItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
