package external

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i648e92ed22999203da3c8fad3bc63deefe974fd0d511e7f830d70ea0aff57ffc "github.com/microsoftgraph/msgraph-sdk-go/models/externalconnectors"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder provides operations to manage the performedBy property of the microsoft.graph.externalConnectors.externalActivity entity.
type ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetQueryParameters represents an identity used to identify who is responsible for the activity.
type ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetQueryParameters
}
// NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderInternal instantiates a new ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder and sets the default values.
func NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) {
    m := &ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/external/connections/{externalConnection%2Did}/items/{externalItem%2Did}/activities/{externalActivity%2Did}/performedBy{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder instantiates a new ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder and sets the default values.
func NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderInternal(urlParams, requestAdapter)
}
// Get represents an identity used to identify who is responsible for the activity.
// returns a Identityable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) Get(ctx context.Context, requestConfiguration *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetRequestConfiguration)(i648e92ed22999203da3c8fad3bc63deefe974fd0d511e7f830d70ea0aff57ffc.Identityable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, i648e92ed22999203da3c8fad3bc63deefe974fd0d511e7f830d70ea0aff57ffc.CreateIdentityFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(i648e92ed22999203da3c8fad3bc63deefe974fd0d511e7f830d70ea0aff57ffc.Identityable), nil
}
// ToGetRequestInformation represents an identity used to identify who is responsible for the activity.
// returns a *RequestInformation when successful
func (m *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder when successful
func (m *ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) WithUrl(rawUrl string)(*ConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder) {
    return NewConnectionsItemItemsItemActivitiesItemPerformedByRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
