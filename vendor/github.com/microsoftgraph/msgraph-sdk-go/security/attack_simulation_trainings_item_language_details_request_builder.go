package security

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// AttackSimulationTrainingsItemLanguageDetailsRequestBuilder provides operations to manage the languageDetails property of the microsoft.graph.training entity.
type AttackSimulationTrainingsItemLanguageDetailsRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetQueryParameters language specific details on a training.
type AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetQueryParameters struct {
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
// AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetQueryParameters
}
// AttackSimulationTrainingsItemLanguageDetailsRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type AttackSimulationTrainingsItemLanguageDetailsRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByTrainingLanguageDetailId provides operations to manage the languageDetails property of the microsoft.graph.training entity.
// returns a *AttackSimulationTrainingsItemLanguageDetailsTrainingLanguageDetailItemRequestBuilder when successful
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) ByTrainingLanguageDetailId(trainingLanguageDetailId string)(*AttackSimulationTrainingsItemLanguageDetailsTrainingLanguageDetailItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if trainingLanguageDetailId != "" {
        urlTplParams["trainingLanguageDetail%2Did"] = trainingLanguageDetailId
    }
    return NewAttackSimulationTrainingsItemLanguageDetailsTrainingLanguageDetailItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilderInternal instantiates a new AttackSimulationTrainingsItemLanguageDetailsRequestBuilder and sets the default values.
func NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) {
    m := &AttackSimulationTrainingsItemLanguageDetailsRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/security/attackSimulation/trainings/{training%2Did}/languageDetails{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilder instantiates a new AttackSimulationTrainingsItemLanguageDetailsRequestBuilder and sets the default values.
func NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *AttackSimulationTrainingsItemLanguageDetailsCountRequestBuilder when successful
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) Count()(*AttackSimulationTrainingsItemLanguageDetailsCountRequestBuilder) {
    return NewAttackSimulationTrainingsItemLanguageDetailsCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get language specific details on a training.
// returns a TrainingLanguageDetailCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) Get(ctx context.Context, requestConfiguration *AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTrainingLanguageDetailCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailCollectionResponseable), nil
}
// Post create new navigation property to languageDetails for security
// returns a TrainingLanguageDetailable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailable, requestConfiguration *AttackSimulationTrainingsItemLanguageDetailsRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateTrainingLanguageDetailFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailable), nil
}
// ToGetRequestInformation language specific details on a training.
// returns a *RequestInformation when successful
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *AttackSimulationTrainingsItemLanguageDetailsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to languageDetails for security
// returns a *RequestInformation when successful
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.TrainingLanguageDetailable, requestConfiguration *AttackSimulationTrainingsItemLanguageDetailsRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder when successful
func (m *AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) WithUrl(rawUrl string)(*AttackSimulationTrainingsItemLanguageDetailsRequestBuilder) {
    return NewAttackSimulationTrainingsItemLanguageDetailsRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
