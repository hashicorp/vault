package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder provides operations to manage the deploymentSummary property of the microsoft.graph.defaultManagedAppProtection entity.
type DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetQueryParameters navigation property to deployment summary of the configuration.
type DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetQueryParameters
}
// DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderInternal instantiates a new DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder and sets the default values.
func NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) {
    m := &DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/defaultManagedAppProtections/{defaultManagedAppProtection%2Did}/deploymentSummary{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder instantiates a new DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder and sets the default values.
func NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property deploymentSummary for deviceAppManagement
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) Delete(ctx context.Context, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get navigation property to deployment summary of the configuration.
// returns a ManagedAppPolicyDeploymentSummaryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) Get(ctx context.Context, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable), nil
}
// Patch update the navigation property deploymentSummary in deviceAppManagement
// returns a ManagedAppPolicyDeploymentSummaryable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateManagedAppPolicyDeploymentSummaryFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable), nil
}
// ToDeleteRequestInformation delete navigation property deploymentSummary for deviceAppManagement
// returns a *RequestInformation when successful
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation navigation property to deployment summary of the configuration.
// returns a *RequestInformation when successful
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property deploymentSummary in deviceAppManagement
// returns a *RequestInformation when successful
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ManagedAppPolicyDeploymentSummaryable, requestConfiguration *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder when successful
func (m *DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) WithUrl(rawUrl string)(*DefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder) {
    return NewDefaultManagedAppProtectionsItemDeploymentSummaryRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
