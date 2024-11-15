package admin

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// PeopleRequestBuilder provides operations to manage the people property of the microsoft.graph.admin entity.
type PeopleRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// PeopleRequestBuilderGetQueryParameters retrieve the properties and relationships of a peopleAdminSettings object.
type PeopleRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PeopleRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PeopleRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PeopleRequestBuilderGetQueryParameters
}
// NewPeopleRequestBuilderInternal instantiates a new PeopleRequestBuilder and sets the default values.
func NewPeopleRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PeopleRequestBuilder) {
    m := &PeopleRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/admin/people{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewPeopleRequestBuilder instantiates a new PeopleRequestBuilder and sets the default values.
func NewPeopleRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PeopleRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPeopleRequestBuilderInternal(urlParams, requestAdapter)
}
// Get retrieve the properties and relationships of a peopleAdminSettings object.
// returns a PeopleAdminSettingsable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/peopleadminsettings-get?view=graph-rest-1.0
func (m *PeopleRequestBuilder) Get(ctx context.Context, requestConfiguration *PeopleRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PeopleAdminSettingsable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePeopleAdminSettingsFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PeopleAdminSettingsable), nil
}
// ItemInsights provides operations to manage the itemInsights property of the microsoft.graph.peopleAdminSettings entity.
// returns a *PeopleItemInsightsRequestBuilder when successful
func (m *PeopleRequestBuilder) ItemInsights()(*PeopleItemInsightsRequestBuilder) {
    return NewPeopleItemInsightsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ProfileCardProperties provides operations to manage the profileCardProperties property of the microsoft.graph.peopleAdminSettings entity.
// returns a *PeopleProfileCardPropertiesRequestBuilder when successful
func (m *PeopleRequestBuilder) ProfileCardProperties()(*PeopleProfileCardPropertiesRequestBuilder) {
    return NewPeopleProfileCardPropertiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Pronouns provides operations to manage the pronouns property of the microsoft.graph.peopleAdminSettings entity.
// returns a *PeoplePronounsRequestBuilder when successful
func (m *PeopleRequestBuilder) Pronouns()(*PeoplePronounsRequestBuilder) {
    return NewPeoplePronounsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation retrieve the properties and relationships of a peopleAdminSettings object.
// returns a *RequestInformation when successful
func (m *PeopleRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *PeopleRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *PeopleRequestBuilder when successful
func (m *PeopleRequestBuilder) WithUrl(rawUrl string)(*PeopleRequestBuilder) {
    return NewPeopleRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
