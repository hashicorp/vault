package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource struct {
    AccessReviewInstanceDecisionItemResource
}
// NewAccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource instantiates a new AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource and sets the default values.
func NewAccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource()(*AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) {
    m := &AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource{
        AccessReviewInstanceDecisionItemResource: *NewAccessReviewInstanceDecisionItemResource(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource(), nil
}
// GetAccessPackageDisplayName gets the accessPackageDisplayName property value. Display name of the access package to which access has been granted.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) GetAccessPackageDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("accessPackageDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAccessPackageId gets the accessPackageId property value. Identifier of the access package to which access has been granted.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) GetAccessPackageId()(*string) {
    val, err := m.GetBackingStore().Get("accessPackageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewInstanceDecisionItemResource.GetFieldDeserializers()
    res["accessPackageDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessPackageDisplayName(val)
        }
        return nil
    }
    res["accessPackageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessPackageId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewInstanceDecisionItemResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("accessPackageDisplayName", m.GetAccessPackageDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("accessPackageId", m.GetAccessPackageId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackageDisplayName sets the accessPackageDisplayName property value. Display name of the access package to which access has been granted.
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) SetAccessPackageDisplayName(value *string)() {
    err := m.GetBackingStore().Set("accessPackageDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAccessPackageId sets the accessPackageId property value. Identifier of the access package to which access has been granted.
func (m *AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResource) SetAccessPackageId(value *string)() {
    err := m.GetBackingStore().Set("accessPackageId", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewInstanceDecisionItemAccessPackageAssignmentPolicyResourceable interface {
    AccessReviewInstanceDecisionItemResourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackageDisplayName()(*string)
    GetAccessPackageId()(*string)
    SetAccessPackageDisplayName(value *string)()
    SetAccessPackageId(value *string)()
}
