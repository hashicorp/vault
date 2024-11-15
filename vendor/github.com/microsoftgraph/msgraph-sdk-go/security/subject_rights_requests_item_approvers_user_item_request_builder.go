package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SubjectRightsRequestsItemApproversUserItemRequestBuilder provides operations to manage the approvers property of the microsoft.graph.subjectRightsRequest entity.
type SubjectRightsRequestsItemApproversUserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SubjectRightsRequestsItemApproversUserItemRequestBuilderGetQueryParameters collection of users who can approve the request. Currently only supported for requests of type delete.
type SubjectRightsRequestsItemApproversUserItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SubjectRightsRequestsItemApproversUserItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SubjectRightsRequestsItemApproversUserItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SubjectRightsRequestsItemApproversUserItemRequestBuilderGetQueryParameters
}
// NewSubjectRightsRequestsItemApproversUserItemRequestBuilderInternal instantiates a new SubjectRightsRequestsItemApproversUserItemRequestBuilder and sets the default values.
func NewSubjectRightsRequestsItemApproversUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SubjectRightsRequestsItemApproversUserItemRequestBuilder) {
    m := &SubjectRightsRequestsItemApproversUserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/subjectRightsRequests/{subjectRightsRequest%2Did}/approvers/{user%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSubjectRightsRequestsItemApproversUserItemRequestBuilder instantiates a new SubjectRightsRequestsItemApproversUserItemRequestBuilder and sets the default values.
func NewSubjectRightsRequestsItemApproversUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SubjectRightsRequestsItemApproversUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSubjectRightsRequestsItemApproversUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Get collection of users who can approve the request. Currently only supported for requests of type delete.
// returns a Userable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SubjectRightsRequestsItemApproversUserItemRequestBuilder) Get(ctx context.Context, requestConfiguration *SubjectRightsRequestsItemApproversUserItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
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
// returns a *SubjectRightsRequestsItemApproversItemMailboxSettingsRequestBuilder when successful
func (m *SubjectRightsRequestsItemApproversUserItemRequestBuilder) MailboxSettings()(*SubjectRightsRequestsItemApproversItemMailboxSettingsRequestBuilder) {
    return NewSubjectRightsRequestsItemApproversItemMailboxSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceProvisioningErrors the serviceProvisioningErrors property
// returns a *SubjectRightsRequestsItemApproversItemServiceProvisioningErrorsRequestBuilder when successful
func (m *SubjectRightsRequestsItemApproversUserItemRequestBuilder) ServiceProvisioningErrors()(*SubjectRightsRequestsItemApproversItemServiceProvisioningErrorsRequestBuilder) {
    return NewSubjectRightsRequestsItemApproversItemServiceProvisioningErrorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation collection of users who can approve the request. Currently only supported for requests of type delete.
// returns a *RequestInformation when successful
func (m *SubjectRightsRequestsItemApproversUserItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SubjectRightsRequestsItemApproversUserItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *SubjectRightsRequestsItemApproversUserItemRequestBuilder when successful
func (m *SubjectRightsRequestsItemApproversUserItemRequestBuilder) WithUrl(rawUrl string)(*SubjectRightsRequestsItemApproversUserItemRequestBuilder) {
    return NewSubjectRightsRequestsItemApproversUserItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
