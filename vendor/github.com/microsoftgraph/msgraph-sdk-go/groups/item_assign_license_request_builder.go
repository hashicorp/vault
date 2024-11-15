package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemAssignLicenseRequestBuilder provides operations to call the assignLicense method.
type ItemAssignLicenseRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemAssignLicenseRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemAssignLicenseRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemAssignLicenseRequestBuilderInternal instantiates a new ItemAssignLicenseRequestBuilder and sets the default values.
func NewItemAssignLicenseRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAssignLicenseRequestBuilder) {
    m := &ItemAssignLicenseRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/assignLicense", pathParameters),
    }
    return m
}
// NewItemAssignLicenseRequestBuilder instantiates a new ItemAssignLicenseRequestBuilder and sets the default values.
func NewItemAssignLicenseRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemAssignLicenseRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemAssignLicenseRequestBuilderInternal(urlParams, requestAdapter)
}
// Post add or remove licenses on the group. Licenses assigned to the group will be assigned to all users in the group. Group-based licensing is an alternative to direct user licensing. To learn more about group-based licensing, see What is group-based licensing in Microsoft Entra ID. To get the subscriptions available in the directory, perform a GET subscribedSkus request.
// returns a Groupable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/group-assignlicense?view=graph-rest-1.0
func (m *ItemAssignLicenseRequestBuilder) Post(ctx context.Context, body ItemAssignLicensePostRequestBodyable, requestConfiguration *ItemAssignLicenseRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateGroupFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Groupable), nil
}
// ToPostRequestInformation add or remove licenses on the group. Licenses assigned to the group will be assigned to all users in the group. Group-based licensing is an alternative to direct user licensing. To learn more about group-based licensing, see What is group-based licensing in Microsoft Entra ID. To get the subscriptions available in the directory, perform a GET subscribedSkus request.
// returns a *RequestInformation when successful
func (m *ItemAssignLicenseRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemAssignLicensePostRequestBodyable, requestConfiguration *ItemAssignLicenseRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemAssignLicenseRequestBuilder when successful
func (m *ItemAssignLicenseRequestBuilder) WithUrl(rawUrl string)(*ItemAssignLicenseRequestBuilder) {
    return NewItemAssignLicenseRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
