package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder builds and executes requests for operations under \education\classes\{educationClass-id}\assignments\{educationAssignment-id}\categories\{educationCategory-id}
type ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal instantiates a new ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    m := &ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/classes/{educationClass%2Did}/assignments/{educationAssignment%2Did}/categories/{educationCategory%2Did}", pathParameters),
    }
    return m
}
// NewClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder instantiates a new ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *ClassesItemAssignmentsItemCategoriesItemRefRequestBuilder when successful
func (m *ClassesItemAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) Ref()(*ClassesItemAssignmentsItemCategoriesItemRefRequestBuilder) {
    return NewClassesItemAssignmentsItemCategoriesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
