package sites

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder provides operations to call the copyToSectionGroup method.
type ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderInternal instantiates a new ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder and sets the default values.
func NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) {
    m := &ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/sites/{site%2Did}/onenote/notebooks/{notebook%2Did}/sectionGroups/{sectionGroup%2Did}/sections/{onenoteSection%2Did}/copyToSectionGroup", pathParameters),
    }
    return m
}
// NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder instantiates a new ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder and sets the default values.
func NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderInternal(urlParams, requestAdapter)
}
// Post for Copy operations, you follow an asynchronous calling pattern:  First call the Copy action, and then poll the operation endpoint for the result.
// returns a OnenoteOperationable when successful
// returns a ODataError error when the service returns a 4XX or 5XX status code
// [Find more info here]
// 
// [Find more info here]: https://learn.microsoft.com/graph/api/section-copytosectiongroup?view=graph-rest-1.0
func (m *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) Post(ctx context.Context, body ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupPostRequestBodyable, requestConfiguration *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderPostRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnenoteOperationable, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "XXX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateOnenoteOperationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.OnenoteOperationable), nil
}
// ToPostRequestInformation for Copy operations, you follow an asynchronous calling pattern:  First call the Copy action, and then poll the operation endpoint for the result.
// returns a *RequestInformation when successful
func (m *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) ToPostRequestInformation(ctx context.Context, body ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupPostRequestBodyable, requestConfiguration *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// returns a *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder when successful
func (m *ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) WithUrl(rawUrl string)(*ItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder) {
    return NewItemOnenoteNotebooksItemSectionGroupsItemSectionsItemCopyToSectionGroupRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
