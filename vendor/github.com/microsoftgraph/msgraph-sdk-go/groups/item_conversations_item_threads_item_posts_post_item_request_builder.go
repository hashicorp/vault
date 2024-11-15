package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemConversationsItemThreadsItemPostsPostItemRequestBuilder provides operations to manage the posts property of the microsoft.graph.conversationThread entity.
type ItemConversationsItemThreadsItemPostsPostItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetQueryParameters get posts from groups
type ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetQueryParameters
}
// Attachments provides operations to manage the attachments property of the microsoft.graph.post entity.
// returns a *ItemConversationsItemThreadsItemPostsItemAttachmentsRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) Attachments()(*ItemConversationsItemThreadsItemPostsItemAttachmentsRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsItemAttachmentsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewItemConversationsItemThreadsItemPostsPostItemRequestBuilderInternal instantiates a new ItemConversationsItemThreadsItemPostsPostItemRequestBuilder and sets the default values.
func NewItemConversationsItemThreadsItemPostsPostItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) {
    m := &ItemConversationsItemThreadsItemPostsPostItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/groups/{group%2Did}/conversations/{conversation%2Did}/threads/{conversationThread%2Did}/posts/{post%2Did}{?%24expand,%24select}", pathParameters),
    }
    return m
}
// NewItemConversationsItemThreadsItemPostsPostItemRequestBuilder instantiates a new ItemConversationsItemThreadsItemPostsPostItemRequestBuilder and sets the default values.
func NewItemConversationsItemThreadsItemPostsPostItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemConversationsItemThreadsItemPostsPostItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.post entity.
// returns a *ItemConversationsItemThreadsItemPostsItemExtensionsRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) Extensions()(*ItemConversationsItemThreadsItemPostsItemExtensionsRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsItemExtensionsRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Forward provides operations to call the forward method.
// returns a *ItemConversationsItemThreadsItemPostsItemForwardRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) Forward()(*ItemConversationsItemThreadsItemPostsItemForwardRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsItemForwardRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Get get posts from groups
// returns a Postable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) Get(ctx context.Context, requestConfiguration *ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Postable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreatePostFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Postable), nil
}
// InReplyTo provides operations to manage the inReplyTo property of the microsoft.graph.post entity.
// returns a *ItemConversationsItemThreadsItemPostsItemInReplyToRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) InReplyTo()(*ItemConversationsItemThreadsItemPostsItemInReplyToRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsItemInReplyToRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// Reply provides operations to call the reply method.
// returns a *ItemConversationsItemThreadsItemPostsItemReplyRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) Reply()(*ItemConversationsItemThreadsItemPostsItemReplyRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsItemReplyRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// ToGetRequestInformation get posts from groups
// returns a *RequestInformation when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *ItemConversationsItemThreadsItemPostsPostItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder when successful
func (m *ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) WithUrl(rawUrl string)(*ItemConversationsItemThreadsItemPostsPostItemRequestBuilder) {
    return NewItemConversationsItemThreadsItemPostsPostItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
