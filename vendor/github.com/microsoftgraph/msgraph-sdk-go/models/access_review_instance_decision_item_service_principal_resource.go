package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewInstanceDecisionItemServicePrincipalResource struct {
    AccessReviewInstanceDecisionItemResource
}
// NewAccessReviewInstanceDecisionItemServicePrincipalResource instantiates a new AccessReviewInstanceDecisionItemServicePrincipalResource and sets the default values.
func NewAccessReviewInstanceDecisionItemServicePrincipalResource()(*AccessReviewInstanceDecisionItemServicePrincipalResource) {
    m := &AccessReviewInstanceDecisionItemServicePrincipalResource{
        AccessReviewInstanceDecisionItemResource: *NewAccessReviewInstanceDecisionItemResource(),
    }
    odataTypeValue := "#microsoft.graph.accessReviewInstanceDecisionItemServicePrincipalResource"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAccessReviewInstanceDecisionItemServicePrincipalResourceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewInstanceDecisionItemServicePrincipalResourceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInstanceDecisionItemServicePrincipalResource(), nil
}
// GetAppId gets the appId property value. The globally unique identifier of the application to which access has been granted.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) GetAppId()(*string) {
    val, err := m.GetBackingStore().Get("appId")
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
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AccessReviewInstanceDecisionItemResource.GetFieldDeserializers()
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AccessReviewInstanceDecisionItemResource.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appId", m.GetAppId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppId sets the appId property value. The globally unique identifier of the application to which access has been granted.
func (m *AccessReviewInstanceDecisionItemServicePrincipalResource) SetAppId(value *string)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewInstanceDecisionItemServicePrincipalResourceable interface {
    AccessReviewInstanceDecisionItemResourceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppId()(*string)
    SetAppId(value *string)()
}
