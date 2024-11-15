package employeeexperience

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// CommunitiesItemOwnersUserItemRequestBuilder provides operations to manage the owners property of the microsoft.graph.community entity.
type CommunitiesItemOwnersUserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// CommunitiesItemOwnersUserItemRequestBuilderGetQueryParameters the admins of the community. Limited to 100 users. If this property isn't specified when you create the community, the calling user is automatically assigned as the community owner.
type CommunitiesItemOwnersUserItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// CommunitiesItemOwnersUserItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type CommunitiesItemOwnersUserItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *CommunitiesItemOwnersUserItemRequestBuilderGetQueryParameters
}
// NewCommunitiesItemOwnersUserItemRequestBuilderInternal instantiates a new CommunitiesItemOwnersUserItemRequestBuilder and sets the default values.
func NewCommunitiesItemOwnersUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CommunitiesItemOwnersUserItemRequestBuilder) {
    m := &CommunitiesItemOwnersUserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/employeeExperience/communities/{community%2Did}/owners/{user%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewCommunitiesItemOwnersUserItemRequestBuilder instantiates a new CommunitiesItemOwnersUserItemRequestBuilder and sets the default values.
func NewCommunitiesItemOwnersUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*CommunitiesItemOwnersUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewCommunitiesItemOwnersUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the admins of the community. Limited to 100 users. If this property isn't specified when you create the community, the calling user is automatically assigned as the community owner.
// returns a Userable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *CommunitiesItemOwnersUserItemRequestBuilder) Get(ctx context.Context, requestConfiguration *CommunitiesItemOwnersUserItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable), nil
}
// MailboxSettings the mailboxSettings property
// returns a *CommunitiesItemOwnersItemMailboxSettingsRequestBuilder when successful
func (m *CommunitiesItemOwnersUserItemRequestBuilder) MailboxSettings()(*CommunitiesItemOwnersItemMailboxSettingsRequestBuilder) {
    return NewCommunitiesItemOwnersItemMailboxSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceProvisioningErrors the serviceProvisioningErrors property
// returns a *CommunitiesItemOwnersItemServiceProvisioningErrorsRequestBuilder when successful
func (m *CommunitiesItemOwnersUserItemRequestBuilder) ServiceProvisioningErrors()(*CommunitiesItemOwnersItemServiceProvisioningErrorsRequestBuilder) {
    return NewCommunitiesItemOwnersItemServiceProvisioningErrorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation the admins of the community. Limited to 100 users. If this property isn't specified when you create the community, the calling user is automatically assigned as the community owner.
// returns a *RequestInformation when successful
func (m *CommunitiesItemOwnersUserItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *CommunitiesItemOwnersUserItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *CommunitiesItemOwnersUserItemRequestBuilder when successful
func (m *CommunitiesItemOwnersUserItemRequestBuilder) WithUrl(rawUrl string)(*CommunitiesItemOwnersUserItemRequestBuilder) {
    return NewCommunitiesItemOwnersUserItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
