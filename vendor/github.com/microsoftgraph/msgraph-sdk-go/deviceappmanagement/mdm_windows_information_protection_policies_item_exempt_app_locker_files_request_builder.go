package deviceappmanagement

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder provides operations to manage the exemptAppLockerFiles property of the microsoft.graph.windowsInformationProtection entity.
type MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetQueryParameters another way to input exempt apps through xml files
type MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetQueryParameters struct {
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
// MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetQueryParameters
}
// MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// ByWindowsInformationProtectionAppLockerFileId provides operations to manage the exemptAppLockerFiles property of the microsoft.graph.windowsInformationProtection entity.
// returns a *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesWindowsInformationProtectionAppLockerFileItemRequestBuilder when successful
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) ByWindowsInformationProtectionAppLockerFileId(windowsInformationProtectionAppLockerFileId string)(*MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesWindowsInformationProtectionAppLockerFileItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if windowsInformationProtectionAppLockerFileId != "" {
        urlTplParams["windowsInformationProtectionAppLockerFile%2Did"] = windowsInformationProtectionAppLockerFileId
    }
    return NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesWindowsInformationProtectionAppLockerFileItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderInternal instantiates a new MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder and sets the default values.
func NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) {
    m := &MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/deviceAppManagement/mdmWindowsInformationProtectionPolicies/{mdmWindowsInformationProtectionPolicy%2Did}/exemptAppLockerFiles{?%24count,%24expand,%24filter,%24orderby,%24search,%24select,%24skip,%24top}", pathParameters),
    }
    return m
}
// NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder instantiates a new MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder and sets the default values.
func NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesCountRequestBuilder when successful
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) Count()(*MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesCountRequestBuilder) {
    return NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get another way to input exempt apps through xml files
// returns a WindowsInformationProtectionAppLockerFileCollectionResponseable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) Get(ctx context.Context, requestConfiguration *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileCollectionResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsInformationProtectionAppLockerFileCollectionResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileCollectionResponseable), nil
}
// Post create new navigation property to exemptAppLockerFiles for deviceAppManagement
// returns a WindowsInformationProtectionAppLockerFileable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) Post(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileable, requestConfiguration *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateWindowsInformationProtectionAppLockerFileFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileable), nil
}
// ToGetRequestInformation another way to input exempt apps through xml files
// returns a *RequestInformation when successful
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// ToPostRequestInformation create new navigation property to exemptAppLockerFiles for deviceAppManagement
// returns a *RequestInformation when successful
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) ToPostRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.WindowsInformationProtectionAppLockerFileable, requestConfiguration *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder when successful
func (m *MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) WithUrl(rawUrl string)(*MdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder) {
    return NewMdmWindowsInformationProtectionPoliciesItemExemptAppLockerFilesRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
