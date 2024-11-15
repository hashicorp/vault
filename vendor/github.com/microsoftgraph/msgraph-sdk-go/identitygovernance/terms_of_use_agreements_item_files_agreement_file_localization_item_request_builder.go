package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder provides operations to manage the files property of the microsoft.graph.agreement entity.
type TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetQueryParameters pDFs linked to this agreement. This property is in the process of being deprecated. Use the  file property instead. Supports $expand.
type TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetQueryParameters
}
// TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderInternal instantiates a new TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder and sets the default values.
func NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) {
    m := &TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/termsOfUse/agreements/{agreement%2Did}/files/{agreementFileLocalization%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder instantiates a new TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder and sets the default values.
func NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property files for identityGovernance
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get pDFs linked to this agreement. This property is in the process of being deprecated. Use the  file property instead. Supports $expand.
// returns a AgreementFileLocalizationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAgreementFileLocalizationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable), nil
}
// Patch update the navigation property files in identityGovernance
// returns a AgreementFileLocalizationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateAgreementFileLocalizationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable), nil
}
// ToDeleteRequestInformation delete navigation property files for identityGovernance
// returns a *RequestInformation when successful
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation pDFs linked to this agreement. This property is in the process of being deprecated. Use the  file property instead. Supports $expand.
// returns a *RequestInformation when successful
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property files in identityGovernance
// returns a *RequestInformation when successful
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.AgreementFileLocalizationable, requestConfiguration *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// Versions provides operations to manage the versions property of the microsoft.graph.agreementFileLocalization entity.
// returns a *TermsOfUseAgreementsItemFilesItemVersionsRequestBuilder when successful
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) Versions()(*TermsOfUseAgreementsItemFilesItemVersionsRequestBuilder) {
    return NewTermsOfUseAgreementsItemFilesItemVersionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder when successful
func (m *TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) WithUrl(rawUrl string)(*TermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder) {
    return NewTermsOfUseAgreementsItemFilesAgreementFileLocalizationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
