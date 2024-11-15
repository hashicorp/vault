package identitygovernance

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder provides operations to manage the stages property of the microsoft.graph.approval entity.
type AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetQueryParameters a collection of stages in the approval decision.
type AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetQueryParameters struct {
    // Include count of items
    Count *bool `uriparametername:"%24count"`
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Filter items by property values
    Filter *string `uriparametername:"%24filter"`
    // Order items by property values
    Orderby []string `uriparametername:"%24orderby"`
    // Search items by search phrases
    Search *string `uriparametername:"%24search"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
    // Skip the first n items
    Skip *int32 `uriparametername:"%24skip"`
    // Show only the first n items
    Top *int32 `uriparametername:"%24top"`
}
// AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetQueryParameters
}
// AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByApprovalStageId provides operations to manage the stages property of the microsoft.graph.approval entity.
// returns a *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesApprovalStageItemRequestBuilder when successful
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) ByApprovalStageId(approvalStageId string)(*AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesApprovalStageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if approvalStageId != "" {
        urlTplParams["approvalStage%2Did"] = approvalStageId
    }
    return NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesApprovalStageItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderInternal instantiates a new AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder and sets the default values.
func NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) {
    m := &AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/identityGovernance/appConsent/appConsentRequests/{appConsentRequest%2Did}/userConsentRequests/{userConsentRequest%2Did}/approval/stages{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder instantiates a new AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder and sets the default values.
func NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesCountRequestBuilder when successful
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) Count()(*AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesCountRequestBuilder) {
    return NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get a collection of stages in the approval decision.
// returns a ApprovalStageCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) Get(ctx context.Context, requestConfiguration *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApprovalStageCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageCollectionResponseable), nil
}
// Post create new navigation property to stages for identityGovernance
// returns a ApprovalStageable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageable, requestConfiguration *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateApprovalStageFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageable), nil
}
// ToGetRequestInformation a collection of stages in the approval decision.
// returns a *RequestInformation when successful
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to stages for identityGovernance
// returns a *RequestInformation when successful
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ApprovalStageable, requestConfiguration *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
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
// returns a *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder when successful
func (m *AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) WithUrl(rawUrl string)(*AppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder) {
    return NewAppConsentAppConsentRequestsItemUserConsentRequestsItemApprovalStagesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
