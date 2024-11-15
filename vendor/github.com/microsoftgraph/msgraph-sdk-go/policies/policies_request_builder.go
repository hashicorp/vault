package policies

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// PoliciesRequestBuilder provides operations to manage the policyRoot singleton.
type PoliciesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// PoliciesRequestBuilderGetQueryParameters get policies
type PoliciesRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// PoliciesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PoliciesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *PoliciesRequestBuilderGetQueryParameters
}
// PoliciesRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type PoliciesRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ActivityBasedTimeoutPolicies provides operations to manage the activityBasedTimeoutPolicies property of the microsoft.graph.policyRoot entity.
// returns a *ActivityBasedTimeoutPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) ActivityBasedTimeoutPolicies()(*ActivityBasedTimeoutPoliciesRequestBuilder) {
    return NewActivityBasedTimeoutPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AdminConsentRequestPolicy provides operations to manage the adminConsentRequestPolicy property of the microsoft.graph.policyRoot entity.
// returns a *AdminConsentRequestPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) AdminConsentRequestPolicy()(*AdminConsentRequestPolicyRequestBuilder) {
    return NewAdminConsentRequestPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AppManagementPolicies provides operations to manage the appManagementPolicies property of the microsoft.graph.policyRoot entity.
// returns a *AppManagementPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) AppManagementPolicies()(*AppManagementPoliciesRequestBuilder) {
    return NewAppManagementPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationFlowsPolicy provides operations to manage the authenticationFlowsPolicy property of the microsoft.graph.policyRoot entity.
// returns a *AuthenticationFlowsPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) AuthenticationFlowsPolicy()(*AuthenticationFlowsPolicyRequestBuilder) {
    return NewAuthenticationFlowsPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationMethodsPolicy provides operations to manage the authenticationMethodsPolicy property of the microsoft.graph.policyRoot entity.
// returns a *AuthenticationMethodsPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) AuthenticationMethodsPolicy()(*AuthenticationMethodsPolicyRequestBuilder) {
    return NewAuthenticationMethodsPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthenticationStrengthPolicies provides operations to manage the authenticationStrengthPolicies property of the microsoft.graph.policyRoot entity.
// returns a *AuthenticationStrengthPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) AuthenticationStrengthPolicies()(*AuthenticationStrengthPoliciesRequestBuilder) {
    return NewAuthenticationStrengthPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AuthorizationPolicy provides operations to manage the authorizationPolicy property of the microsoft.graph.policyRoot entity.
// returns a *AuthorizationPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) AuthorizationPolicy()(*AuthorizationPolicyRequestBuilder) {
    return NewAuthorizationPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ClaimsMappingPolicies provides operations to manage the claimsMappingPolicies property of the microsoft.graph.policyRoot entity.
// returns a *ClaimsMappingPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) ClaimsMappingPolicies()(*ClaimsMappingPoliciesRequestBuilder) {
    return NewClaimsMappingPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ConditionalAccessPolicies provides operations to manage the conditionalAccessPolicies property of the microsoft.graph.policyRoot entity.
// returns a *ConditionalAccessPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) ConditionalAccessPolicies()(*ConditionalAccessPoliciesRequestBuilder) {
    return NewConditionalAccessPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewPoliciesRequestBuilderInternal instantiates a new PoliciesRequestBuilder and sets the default values.
func NewPoliciesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PoliciesRequestBuilder) {
    m := &PoliciesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/policies{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewPoliciesRequestBuilder instantiates a new PoliciesRequestBuilder and sets the default values.
func NewPoliciesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PoliciesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPoliciesRequestBuilderInternal(urlParams, requestAdapter)
}
// CrossTenantAccessPolicy provides operations to manage the crossTenantAccessPolicy property of the microsoft.graph.policyRoot entity.
// returns a *CrossTenantAccessPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) CrossTenantAccessPolicy()(*CrossTenantAccessPolicyRequestBuilder) {
    return NewCrossTenantAccessPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DefaultAppManagementPolicy provides operations to manage the defaultAppManagementPolicy property of the microsoft.graph.policyRoot entity.
// returns a *DefaultAppManagementPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) DefaultAppManagementPolicy()(*DefaultAppManagementPolicyRequestBuilder) {
    return NewDefaultAppManagementPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// DeviceRegistrationPolicy provides operations to manage the deviceRegistrationPolicy property of the microsoft.graph.policyRoot entity.
// returns a *DeviceRegistrationPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) DeviceRegistrationPolicy()(*DeviceRegistrationPolicyRequestBuilder) {
    return NewDeviceRegistrationPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// FeatureRolloutPolicies provides operations to manage the featureRolloutPolicies property of the microsoft.graph.policyRoot entity.
// returns a *FeatureRolloutPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) FeatureRolloutPolicies()(*FeatureRolloutPoliciesRequestBuilder) {
    return NewFeatureRolloutPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get policies
// returns a PolicyRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *PoliciesRequestBuilder) Get(ctx context.Context, requestConfiguration *PoliciesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePolicyRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable), nil
}
// HomeRealmDiscoveryPolicies provides operations to manage the homeRealmDiscoveryPolicies property of the microsoft.graph.policyRoot entity.
// returns a *HomeRealmDiscoveryPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) HomeRealmDiscoveryPolicies()(*HomeRealmDiscoveryPoliciesRequestBuilder) {
    return NewHomeRealmDiscoveryPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// IdentitySecurityDefaultsEnforcementPolicy provides operations to manage the identitySecurityDefaultsEnforcementPolicy property of the microsoft.graph.policyRoot entity.
// returns a *IdentitySecurityDefaultsEnforcementPolicyRequestBuilder when successful
func (m *PoliciesRequestBuilder) IdentitySecurityDefaultsEnforcementPolicy()(*IdentitySecurityDefaultsEnforcementPolicyRequestBuilder) {
    return NewIdentitySecurityDefaultsEnforcementPolicyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update policies
// returns a PolicyRootable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *PoliciesRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, requestConfiguration *PoliciesRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePolicyRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable), nil
}
// PermissionGrantPolicies provides operations to manage the permissionGrantPolicies property of the microsoft.graph.policyRoot entity.
// returns a *PermissionGrantPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) PermissionGrantPolicies()(*PermissionGrantPoliciesRequestBuilder) {
    return NewPermissionGrantPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoleManagementPolicies provides operations to manage the roleManagementPolicies property of the microsoft.graph.policyRoot entity.
// returns a *RoleManagementPoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) RoleManagementPolicies()(*RoleManagementPoliciesRequestBuilder) {
    return NewRoleManagementPoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// RoleManagementPolicyAssignments provides operations to manage the roleManagementPolicyAssignments property of the microsoft.graph.policyRoot entity.
// returns a *RoleManagementPolicyAssignmentsRequestBuilder when successful
func (m *PoliciesRequestBuilder) RoleManagementPolicyAssignments()(*RoleManagementPolicyAssignmentsRequestBuilder) {
    return NewRoleManagementPolicyAssignmentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation get policies
// returns a *RequestInformation when successful
func (m *PoliciesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *PoliciesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// TokenIssuancePolicies provides operations to manage the tokenIssuancePolicies property of the microsoft.graph.policyRoot entity.
// returns a *TokenIssuancePoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) TokenIssuancePolicies()(*TokenIssuancePoliciesRequestBuilder) {
    return NewTokenIssuancePoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TokenLifetimePolicies provides operations to manage the tokenLifetimePolicies property of the microsoft.graph.policyRoot entity.
// returns a *TokenLifetimePoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) TokenLifetimePolicies()(*TokenLifetimePoliciesRequestBuilder) {
    return NewTokenLifetimePoliciesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToPatchRequestInformation update policies
// returns a *RequestInformation when successful
func (m *PoliciesRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.PolicyRootable, requestConfiguration *PoliciesRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *PoliciesRequestBuilder when successful
func (m *PoliciesRequestBuilder) WithUrl(rawUrl string)(*PoliciesRequestBuilder) {
    return NewPoliciesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
