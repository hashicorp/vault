package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemGroupLifecyclePoliciesItemAddGroupPostResponseable instead.
type ItemGroupLifecyclePoliciesItemAddGroupResponse struct {
    ItemGroupLifecyclePoliciesItemAddGroupPostResponse
}
// NewItemGroupLifecyclePoliciesItemAddGroupResponse instantiates a new ItemGroupLifecyclePoliciesItemAddGroupResponse and sets the default values.
func NewItemGroupLifecyclePoliciesItemAddGroupResponse()(*ItemGroupLifecyclePoliciesItemAddGroupResponse) {
    m := &ItemGroupLifecyclePoliciesItemAddGroupResponse{
        ItemGroupLifecyclePoliciesItemAddGroupPostResponse: *NewItemGroupLifecyclePoliciesItemAddGroupPostResponse(),
    }
    return m
}
// CreateItemGroupLifecyclePoliciesItemAddGroupResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemGroupLifecyclePoliciesItemAddGroupResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemGroupLifecyclePoliciesItemAddGroupResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemGroupLifecyclePoliciesItemAddGroupPostResponseable instead.
type ItemGroupLifecyclePoliciesItemAddGroupResponseable interface {
    ItemGroupLifecyclePoliciesItemAddGroupPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
