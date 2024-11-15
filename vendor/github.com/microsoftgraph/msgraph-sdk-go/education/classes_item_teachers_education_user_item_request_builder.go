package education

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// ClassesItemTeachersEducationUserItemRequestBuilder builds and executes requests for operations under \education\classes\{educationClass-id}\teachers\{educationUser-id}
type ClassesItemTeachersEducationUserItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// NewClassesItemTeachersEducationUserItemRequestBuilderInternal instantiates a new ClassesItemTeachersEducationUserItemRequestBuilder and sets the default values.
func NewClassesItemTeachersEducationUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemTeachersEducationUserItemRequestBuilder) {
    m := &ClassesItemTeachersEducationUserItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/education/classes/{educationClass%2Did}/teachers/{educationUser%2Did}", pathParameters),
    }
    return m
}
// NewClassesItemTeachersEducationUserItemRequestBuilder instantiates a new ClassesItemTeachersEducationUserItemRequestBuilder and sets the default values.
func NewClassesItemTeachersEducationUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ClassesItemTeachersEducationUserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewClassesItemTeachersEducationUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// Ref provides operations to manage the collection of educationRoot entities.
// returns a *ClassesItemTeachersItemRefRequestBuilder when successful
func (m *ClassesItemTeachersEducationUserItemRequestBuilder) Ref()(*ClassesItemTeachersItemRefRequestBuilder) {
    return NewClassesItemTeachersItemRefRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
