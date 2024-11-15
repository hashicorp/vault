package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// SecurityRequestBuilder provides operations to manage the security singleton.
type SecurityRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// SecurityRequestBuilderGetQueryParameters get security
type SecurityRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// SecurityRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *SecurityRequestBuilderGetQueryParameters
}
// SecurityRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type SecurityRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Alerts provides operations to manage the alerts property of the microsoft.graph.security entity.
// returns a *AlertsRequestBuilder when successful
func (m *SecurityRequestBuilder) Alerts()(*AlertsRequestBuilder) {
    return NewAlertsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Alerts_v2 provides operations to manage the alerts_v2 property of the microsoft.graph.security entity.
// returns a *Alerts_v2RequestBuilder when successful
func (m *SecurityRequestBuilder) Alerts_v2()(*Alerts_v2RequestBuilder) {
    return NewAlerts_v2RequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// AttackSimulation provides operations to manage the attackSimulation property of the microsoft.graph.security entity.
// returns a *AttackSimulationRequestBuilder when successful
func (m *SecurityRequestBuilder) AttackSimulation()(*AttackSimulationRequestBuilder) {
    return NewAttackSimulationRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Cases provides operations to manage the cases property of the microsoft.graph.security entity.
// returns a *CasesRequestBuilder when successful
func (m *SecurityRequestBuilder) Cases()(*CasesRequestBuilder) {
    return NewCasesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewSecurityRequestBuilderInternal instantiates a new SecurityRequestBuilder and sets the default values.
func NewSecurityRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityRequestBuilder) {
    m := &SecurityRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewSecurityRequestBuilder instantiates a new SecurityRequestBuilder and sets the default values.
func NewSecurityRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SecurityRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSecurityRequestBuilderInternal(urlParams, requestAdapter)
}
// Get get security
// returns a Securityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SecurityRequestBuilder) Get(ctx context.Context, requestConfiguration *SecurityRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecurityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable), nil
}
// Identities provides operations to manage the identities property of the microsoft.graph.security entity.
// returns a *IdentitiesRequestBuilder when successful
func (m *SecurityRequestBuilder) Identities()(*IdentitiesRequestBuilder) {
    return NewIdentitiesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Incidents provides operations to manage the incidents property of the microsoft.graph.security entity.
// returns a *IncidentsRequestBuilder when successful
func (m *SecurityRequestBuilder) Incidents()(*IncidentsRequestBuilder) {
    return NewIncidentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Labels provides operations to manage the labels property of the microsoft.graph.security entity.
// returns a *LabelsRequestBuilder when successful
func (m *SecurityRequestBuilder) Labels()(*LabelsRequestBuilder) {
    return NewLabelsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// MicrosoftGraphSecurityRunHuntingQuery provides operations to call the runHuntingQuery method.
// returns a *MicrosoftGraphSecurityRunHuntingQueryRequestBuilder when successful
func (m *SecurityRequestBuilder) MicrosoftGraphSecurityRunHuntingQuery()(*MicrosoftGraphSecurityRunHuntingQueryRequestBuilder) {
    return NewMicrosoftGraphSecurityRunHuntingQueryRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Patch update security
// returns a Securityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *SecurityRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, requestConfiguration *SecurityRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateSecurityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable), nil
}
// SecureScoreControlProfiles provides operations to manage the secureScoreControlProfiles property of the microsoft.graph.security entity.
// returns a *SecureScoreControlProfilesRequestBuilder when successful
func (m *SecurityRequestBuilder) SecureScoreControlProfiles()(*SecureScoreControlProfilesRequestBuilder) {
    return NewSecureScoreControlProfilesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SecureScores provides operations to manage the secureScores property of the microsoft.graph.security entity.
// returns a *SecureScoresRequestBuilder when successful
func (m *SecurityRequestBuilder) SecureScores()(*SecureScoresRequestBuilder) {
    return NewSecureScoresRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// SubjectRightsRequests provides operations to manage the subjectRightsRequests property of the microsoft.graph.security entity.
// returns a *SubjectRightsRequestsRequestBuilder when successful
func (m *SecurityRequestBuilder) SubjectRightsRequests()(*SubjectRightsRequestsRequestBuilder) {
    return NewSubjectRightsRequestsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ThreatIntelligence provides operations to manage the threatIntelligence property of the microsoft.graph.security entity.
// returns a *ThreatIntelligenceRequestBuilder when successful
func (m *SecurityRequestBuilder) ThreatIntelligence()(*ThreatIntelligenceRequestBuilder) {
    return NewThreatIntelligenceRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation get security
// returns a *RequestInformation when successful
func (m *SecurityRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *SecurityRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update security
// returns a *RequestInformation when successful
func (m *SecurityRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Securityable, requestConfiguration *SecurityRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Triggers provides operations to manage the triggers property of the microsoft.graph.security entity.
// returns a *TriggersRequestBuilder when successful
func (m *SecurityRequestBuilder) Triggers()(*TriggersRequestBuilder) {
    return NewTriggersRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TriggerTypes provides operations to manage the triggerTypes property of the microsoft.graph.security entity.
// returns a *TriggerTypesRequestBuilder when successful
func (m *SecurityRequestBuilder) TriggerTypes()(*TriggerTypesRequestBuilder) {
    return NewTriggerTypesRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *SecurityRequestBuilder when successful
func (m *SecurityRequestBuilder) WithUrl(rawUrl string)(*SecurityRequestBuilder) {
    return NewSecurityRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
