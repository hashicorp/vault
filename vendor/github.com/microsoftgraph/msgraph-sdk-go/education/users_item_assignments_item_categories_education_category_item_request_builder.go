package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder builds and executes requests for operations under \education\users\{educationUser-id}\assignments\{educationAssignment-id}\categories\{educationCategory-id}
type UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewUsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal instantiates a new UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewUsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    m := &UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/users/{educationUser%2Did}/assignments/{educationAssignment%2Did}/categories/{educationCategory%2Did}", pathParameters),
    }
    return m
}
// NewUsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder instantiates a new UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewUsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *UsersItemAssignmentsItemCategoriesItemRefRequestBuilder when successful
func (m *UsersItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) Ref()(*UsersItemAssignmentsItemCategoriesItemRefRequestBuilder) {
    return NewUsersItemAssignmentsItemCategoriesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
