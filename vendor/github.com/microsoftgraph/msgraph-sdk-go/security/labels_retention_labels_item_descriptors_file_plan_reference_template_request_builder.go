package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder provides operations to manage the filePlanReferenceTemplate property of the microsoft.graph.security.filePlanDescriptor entity.
type LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetQueryParameters specifies a unique alpha-numeric identifier for an organization’s retention schedule.
type LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetQueryParameters
}
// NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderInternal instantiates a new LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder and sets the default values.
func NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) {
    m := &LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/labels/retentionLabels/{retentionLabel%2Did}/descriptors/filePlanReferenceTemplate{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder instantiates a new LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder and sets the default values.
func NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderInternal(urlParams, requestAdapter)
}
// Get specifies a unique alpha-numeric identifier for an organization’s retention schedule.
// returns a FilePlanReferenceTemplateable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) Get(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.FilePlanReferenceTemplateable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateFilePlanReferenceTemplateFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.FilePlanReferenceTemplateable), nil
}
// ToGetRequestInformation specifies a unique alpha-numeric identifier for an organization’s retention schedule.
// returns a *RequestInformation when successful
func (m *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder when successful
func (m *LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) WithUrl(rawUrl string)(*LabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder) {
    return NewLabelsRetentionLabelsItemDescriptorsFilePlanReferenceTemplateRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
