package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

// LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder provides operations to manage the dispositionReviewStages property of the microsoft.graph.security.retentionLabel entity.
type LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetQueryParameters when action at the end of retention is chosen as 'dispositionReview', dispositionReviewStages specifies a sequential set of stages with at least one reviewer in each stage.
type LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetQueryParameters
}
// LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderInternal instantiates a new LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder and sets the default values.
func NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) {
    m := &LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/labels/retentionLabels/{retentionLabel%2Did}/dispositionReviewStages/{dispositionReviewStage%2DstageNumber}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder instantiates a new LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder and sets the default values.
func NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete delete navigation property dispositionReviewStages for security
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get when action at the end of retention is chosen as 'dispositionReview', dispositionReviewStages specifies a sequential set of stages with at least one reviewer in each stage.
// returns a DispositionReviewStageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) Get(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateDispositionReviewStageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable), nil
}
// Patch update the navigation property dispositionReviewStages in security
// returns a DispositionReviewStageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) Patch(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderPatchRequestConfiguration)(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.CreateDispositionReviewStageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable), nil
}
// ToDeleteRequestInformation delete navigation property dispositionReviewStages for security
// returns a *RequestInformation when successful
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation when action at the end of retention is chosen as 'dispositionReview', dispositionReviewStages specifies a sequential set of stages with at least one reviewer in each stage.
// returns a *RequestInformation when successful
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property dispositionReviewStages in security
// returns a *RequestInformation when successful
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.DispositionReviewStageable, requestConfiguration *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder when successful
func (m *LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) WithUrl(rawUrl string)(*LabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder) {
    return NewLabelsRetentionLabelsItemDispositionReviewStagesDispositionReviewStageStageNumberItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
