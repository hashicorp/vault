package invitations

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// InvitedUserRequestBuilder provides operations to manage the invitedUser property of the microsoft.graph.invitation entity.
type InvitedUserRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// InvitedUserRequestBuilderGetQueryParameters the user created as part of the invitation creation. Read-only. The id property is required in the request body to reset a redemption status.
type InvitedUserRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// InvitedUserRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type InvitedUserRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *InvitedUserRequestBuilderGetQueryParameters
}
// NewInvitedUserRequestBuilderInternal instantiates a new InvitedUserRequestBuilder and sets the default values.
func NewInvitedUserRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InvitedUserRequestBuilder) {
    m := &InvitedUserRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/invitations/invitedUser{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewInvitedUserRequestBuilder instantiates a new InvitedUserRequestBuilder and sets the default values.
func NewInvitedUserRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*InvitedUserRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewInvitedUserRequestBuilderInternal(urlParams, requestAdapter)
}
// Get the user created as part of the invitation creation. Read-only. The id property is required in the request body to reset a redemption status.
// returns a Userable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *InvitedUserRequestBuilder) Get(ctx context.Context, requestConfiguration *InvitedUserRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
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
// returns a *InvitedUserMailboxSettingsRequestBuilder when successful
func (m *InvitedUserRequestBuilder) MailboxSettings()(*InvitedUserMailboxSettingsRequestBuilder) {
    return NewInvitedUserMailboxSettingsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ServiceProvisioningErrors the serviceProvisioningErrors property
// returns a *InvitedUserServiceProvisioningErrorsRequestBuilder when successful
func (m *InvitedUserRequestBuilder) ServiceProvisioningErrors()(*InvitedUserServiceProvisioningErrorsRequestBuilder) {
    return NewInvitedUserServiceProvisioningErrorsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation the user created as part of the invitation creation. Read-only. The id property is required in the request body to reset a redemption status.
// returns a *RequestInformation when successful
func (m *InvitedUserRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *InvitedUserRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *InvitedUserRequestBuilder when successful
func (m *InvitedUserRequestBuilder) WithUrl(rawUrl string)(*InvitedUserRequestBuilder) {
    return NewInvitedUserRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
