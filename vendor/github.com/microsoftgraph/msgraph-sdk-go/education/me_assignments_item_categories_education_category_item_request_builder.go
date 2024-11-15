package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder builds and executes requests for operations under \education\me\assignments\{educationAssignment-id}\categories\{educationCategory-id}
type MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewMeAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal instantiates a new MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewMeAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    m := &MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/me/assignments/{educationAssignment%2Did}/categories/{educationCategory%2Did}", pathParameters),
    }
    return m
}
// NewMeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder instantiates a new MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder and sets the default values.
func NewMeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewMeAssignmentsItemCategoriesEducationCategoryItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *MeAssignmentsItemCategoriesItemRefRequestBuilder when successful
func (m *MeAssignmentsItemCategoriesEducationCategoryItemRequestBuilder) Ref()(*MeAssignmentsItemCategoriesItemRefRequestBuilder) {
    return NewMeAssignmentsItemCategoriesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
