package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder provides operations to manage the secureScoreControlProfiles property of the microsoft.graph.security entity.
type SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetQueryParameters retrieve the properties and relationships of an securescorecontrolprofile object.
type SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetQueryParameters
}
// SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderInternal instantiates a new SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder and sets the default values.
func NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) {
    m := &SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/secureScoreControlProfiles/{secureScoreControlProfile%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder instantiates a new SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder and sets the default values.
func NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property secureScoreControlProfiles for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the properties and relationships of an securescorecontrolprofile object.
// returns a SecureScoreControlProfileable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securescorecontrolprofile-get?view=graph-rest-1.0
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) Get(ctx context.Context, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecureScoreControlProfileFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable), nil
}
// Patch update an editable secureScoreControlProfile object within any integrated solution to change various properties, such as assignedTo or tenantNote.
// returns a SecureScoreControlProfileable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/securescorecontrolprofile-update?view=graph-rest-1.0
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecureScoreControlProfileFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable), nil
}
// ToDeleteRequestInformation delete navigation property secureScoreControlProfiles for security
// returns a *RequestInformation when successful
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the properties and relationships of an securescorecontrolprofile object.
// returns a *RequestInformation when successful
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update an editable secureScoreControlProfile object within any integrated solution to change various properties, such as assignedTo or tenantNote.
// returns a *RequestInformation when successful
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.SecureScoreControlProfileable, requestConfiguration *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder when successful
func (m *SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) WithUrl(rawUrl string)(*SecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder) {
    return NewSecureScoreControlProfilesSecureScoreControlProfileItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
