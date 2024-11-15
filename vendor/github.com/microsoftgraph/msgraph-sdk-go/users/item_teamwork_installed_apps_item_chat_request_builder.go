package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemTeamworkInstalledAppsItemChatRequestBuilder provides operations to manage the chat property of the microsoft.graph.userScopeTeamsAppInstallation entity.
type ItemTeamworkInstalledAppsItemChatRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemTeamworkInstalledAppsItemChatRequestBuilderGetQueryParameters retrieve the chat of the specified user and Teams app.
type ItemTeamworkInstalledAppsItemChatRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemTeamworkInstalledAppsItemChatRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamworkInstalledAppsItemChatRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemTeamworkInstalledAppsItemChatRequestBuilderGetQueryParameters
}
// NewItemTeamworkInstalledAppsItemChatRequestBuilderInternal instantiates a new ItemTeamworkInstalledAppsItemChatRequestBuilder and sets the default values.
func NewItemTeamworkInstalledAppsItemChatRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamworkInstalledAppsItemChatRequestBuilder) {
    m := &ItemTeamworkInstalledAppsItemChatRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/teamwork/installedApps/{userScopeTeamsAppInstallation%2Did}/chat{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemTeamworkInstalledAppsItemChatRequestBuilder instantiates a new ItemTeamworkInstalledAppsItemChatRequestBuilder and sets the default values.
func NewItemTeamworkInstalledAppsItemChatRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamworkInstalledAppsItemChatRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTeamworkInstalledAppsItemChatRequestBuilderInternal(urlParams, requestAdapter)
}
// Get retrieve the chat of the specified user and Teams app.
// returns a Chatable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/userscopeteamsappinstallation-get-chat?view=graph-rest-1.0
func (m *ItemTeamworkInstalledAppsItemChatRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsItemChatRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateChatFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Chatable), nil
}
// ToGetRequestInformation retrieve the chat of the specified user and Teams app.
// returns a *RequestInformation when successful
func (m *ItemTeamworkInstalledAppsItemChatRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsItemChatRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemTeamworkInstalledAppsItemChatRequestBuilder when successful
func (m *ItemTeamworkInstalledAppsItemChatRequestBuilder) WithUrl(rawUrl string)(*ItemTeamworkInstalledAppsItemChatRequestBuilder) {
    return NewItemTeamworkInstalledAppsItemChatRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
