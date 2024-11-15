package serviceprincipals

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder provides operations to manage the oauth2PermissionGrants property of the microsoft.graph.servicePrincipal entity.
type ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetQueryParameters delegated permission grants authorizing this service principal to access an API on behalf of a signed-in user. Read-only. Nullable.
type ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetQueryParameters
}
// NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderInternal instantiates a new ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder and sets the default values.
func NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) {
    m := &ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/servicePrincipals/{servicePrincipal%2Did}/oauth2PermissionGrants/{oAuth2PermissionGrant%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder instantiates a new ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder and sets the default values.
func NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get delegated permission grants authorizing this service principal to access an API on behalf of a signed-in user. Read-only. Nullable.
// returns a OAuth2PermissionGrantable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OAuth2PermissionGrantable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOAuth2PermissionGrantFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OAuth2PermissionGrantable), nil
}
// ToGetRequestInformation delegated permission grants authorizing this service principal to access an API on behalf of a signed-in user. Read-only. Nullable.
// returns a *RequestInformation when successful
func (m *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder when successful
func (m *ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) WithUrl(rawUrl string)(*ItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder) {
    return NewItemOauth2PermissionGrantsOAuth2PermissionGrantItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
