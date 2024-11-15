package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// SchoolsItemUsersEducationUserItemRequestBuilder builds and executes requests for operations under \education\schools\{educationSchool-id}\users\{educationUser-id}
type SchoolsItemUsersEducationUserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewSchoolsItemUsersEducationUserItemRequestBuilderInternal instantiates a new SchoolsItemUsersEducationUserItemRequestBuilder and sets the default values.
func NewSchoolsItemUsersEducationUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SchoolsItemUsersEducationUserItemRequestBuilder) {
    m := &SchoolsItemUsersEducationUserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/schools/{educationSchool%2Did}/users/{educationUser%2Did}", pathParameters),
    }
    return m
}
// NewSchoolsItemUsersEducationUserItemRequestBuilder instantiates a new SchoolsItemUsersEducationUserItemRequestBuilder and sets the default values.
func NewSchoolsItemUsersEducationUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SchoolsItemUsersEducationUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSchoolsItemUsersEducationUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *SchoolsItemUsersItemRefRequestBuilder when successful
func (m *SchoolsItemUsersEducationUserItemRequestBuilder) Ref()(*SchoolsItemUsersItemRefRequestBuilder) {
    return NewSchoolsItemUsersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
