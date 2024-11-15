package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// SchoolsItemClassesEducationClassItemRequestBuilder builds and executes requests for operations under \education\schools\{educationSchool-id}\classes\{educationClass-id}
type SchoolsItemClassesEducationClassItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewSchoolsItemClassesEducationClassItemRequestBuilderInternal instantiates a new SchoolsItemClassesEducationClassItemRequestBuilder and sets the default values.
func NewSchoolsItemClassesEducationClassItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SchoolsItemClassesEducationClassItemRequestBuilder) {
    m := &SchoolsItemClassesEducationClassItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/schools/{educationSchool%2Did}/classes/{educationClass%2Did}", pathParameters),
    }
    return m
}
// NewSchoolsItemClassesEducationClassItemRequestBuilder instantiates a new SchoolsItemClassesEducationClassItemRequestBuilder and sets the default values.
func NewSchoolsItemClassesEducationClassItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*SchoolsItemClassesEducationClassItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewSchoolsItemClassesEducationClassItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *SchoolsItemClassesItemRefRequestBuilder when successful
func (m *SchoolsItemClassesEducationClassItemRequestBuilder) Ref()(*SchoolsItemClassesItemRefRequestBuilder) {
    return NewSchoolsItemClassesItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
