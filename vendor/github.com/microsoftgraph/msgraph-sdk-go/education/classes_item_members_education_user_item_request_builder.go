package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ClassesItemMembersEducationUserItemRequestBuilder builds and executes requests for operations under \education\classes\{educationClass-id}\members\{educationUser-id}
type ClassesItemMembersEducationUserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewClassesItemMembersEducationUserItemRequestBuilderInternal instantiates a new ClassesItemMembersEducationUserItemRequestBuilder and sets the default values.
func NewClassesItemMembersEducationUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemMembersEducationUserItemRequestBuilder) {
    m := &ClassesItemMembersEducationUserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/classes/{educationClass%2Did}/members/{educationUser%2Did}", pathParameters),
    }
    return m
}
// NewClassesItemMembersEducationUserItemRequestBuilder instantiates a new ClassesItemMembersEducationUserItemRequestBuilder and sets the default values.
func NewClassesItemMembersEducationUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemMembersEducationUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClassesItemMembersEducationUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *ClassesItemMembersItemRefRequestBuilder when successful
func (m *ClassesItemMembersEducationUserItemRequestBuilder) Ref()(*ClassesItemMembersItemRefRequestBuilder) {
    return NewClassesItemMembersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
