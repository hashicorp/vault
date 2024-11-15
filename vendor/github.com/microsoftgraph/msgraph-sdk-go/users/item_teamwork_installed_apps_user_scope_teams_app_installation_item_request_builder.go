package users

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder provides operations to manage the installedApps property of the microsoft.graph.userTeamwork entity.
type ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetQueryParameters retrieve the app installed in the personal scope of the specified user.
type ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetQueryParameters
}
// ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Chat provides operations to manage the chat property of the microsoft.graph.userScopeTeamsAppInstallation entity.
// returns a *ItemTeamworkInstalledAppsItemChatRequestBuilder when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) Chat()(*ItemTeamworkInstalledAppsItemChatRequestBuilder) {
    return NewItemTeamworkInstalledAppsItemChatRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderInternal instantiates a new ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder and sets the default values.
func NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) {
    m := &ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/users/{user%2Did}/teamwork/installedApps/{userScopeTeamsAppInstallation%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder instantiates a new ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder and sets the default values.
func NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Delete uninstall an app from the personal scope of the specified user.
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/userteamwork-delete-installedapps?view=graph-rest-1.0
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderDeleteRequestConfiguration)(error) {
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
// Get retrieve the app installed in the personal scope of the specified user.
// returns a UserScopeTeamsAppInstallationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/userteamwork-get-installedapps?view=graph-rest-1.0
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserScopeTeamsAppInstallationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable), nil
}
// Patch update the navigation property installedApps in users
// returns a UserScopeTeamsAppInstallationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable, error) {
    requestInfo, err := m.ToPatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserScopeTeamsAppInstallationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable), nil
}
// TeamsApp provides operations to manage the teamsApp property of the microsoft.graph.teamsAppInstallation entity.
// returns a *ItemTeamworkInstalledAppsItemTeamsAppRequestBuilder when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) TeamsApp()(*ItemTeamworkInstalledAppsItemTeamsAppRequestBuilder) {
    return NewItemTeamworkInstalledAppsItemTeamsAppRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// TeamsAppDefinition provides operations to manage the teamsAppDefinition property of the microsoft.graph.teamsAppInstallation entity.
// returns a *ItemTeamworkInstalledAppsItemTeamsAppDefinitionRequestBuilder when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) TeamsAppDefinition()(*ItemTeamworkInstalledAppsItemTeamsAppDefinitionRequestBuilder) {
    return NewItemTeamworkInstalledAppsItemTeamsAppDefinitionRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToDeleteRequestInformation uninstall an app from the personal scope of the specified user.
// returns a *RequestInformation when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) ToDeleteRequestInformation(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToGetRequestInformation retrieve the app installed in the personal scope of the specified user.
// returns a *RequestInformation when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPatchRequestInformation update the navigation property installedApps in users
// returns a *RequestInformation when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) ToPatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.UserScopeTeamsAppInstallationable, requestConfiguration *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder when successful
func (m *ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) WithUrl(rawUrl string)(*ItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder) {
    return NewItemTeamworkInstalledAppsUserScopeTeamsAppInstallationItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
