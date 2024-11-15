package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewInstanceDecisionItemAzureRoleResource struct {
    AccessReviewInstanceDecisionItemResource
}
// NewAccessReviewInstanceDecisionItemAzureRoleResource instantiates a new AccessReviewInstanceDecisionItemAzureRoleResource and sets the default values.
func NewAccessReviewInstanceDecisionItemAzureRoleResource()(*AccessReviewInstanceDecisionItemAzureRoleResource) {
    m := &AccessReviewInstanceDecisionItemAzureRoleResource{
        AccessReviewInstanceDecisionItemResource: *NewAccessReviewInstanceDecisionItemResource(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewInstanceDecisionItemAzureRoleResource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewInstanceDecisionItemAzureRoleResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewInstanceDecisionItemAzureRoleResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInstanceDecisionItemAzureRoleResource(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewInstanceDecisionItemAzureRoleResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewInstanceDecisionItemResource.GetFieldDeserializers()
    res["scope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewInstanceDecisionItemResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScope(val.(AccessReviewInstanceDecisionItemResourceable))
        }
        return nil
    }
    return res
}
// GetScope gets the scope property value. Details of the scope this role is associated with.
// returns a AccessReviewInstanceDecisionItemResourceable when successful
func (m *AccessReviewInstanceDecisionItemAzureRoleResource) GetScope()(AccessReviewInstanceDecisionItemResourceable) {
    val, err := m.GetBackingStore().Get("scope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewInstanceDecisionItemResourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewInstanceDecisionItemAzureRoleResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewInstanceDecisionItemResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("scope", m.GetScope())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetScope sets the scope property value. Details of the scope this role is associated with.
func (m *AccessReviewInstanceDecisionItemAzureRoleResource) SetScope(value AccessReviewInstanceDecisionItemResourceable)() {
    err := m.GetBackingStore().Set("scope", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewInstanceDecisionItemAzureRoleResourceable interface {
    AccessReviewInstanceDecisionItemResourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetScope()(AccessReviewInstanceDecisionItemResourceable)
    SetScope(value AccessReviewInstanceDecisionItemResourceable)()
}
