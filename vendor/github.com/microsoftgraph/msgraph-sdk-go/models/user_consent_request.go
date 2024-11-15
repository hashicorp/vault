package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserConsentRequest struct {
    Request
}
// NewUserConsentRequest instantiates a new UserConsentRequest and sets the default values.
func NewUserConsentRequest()(*UserConsentRequest) {
    m := &UserConsentRequest{
        Request: *NewRequest(),
    }
    return m
}
// CreateUserConsentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserConsentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserConsentRequest(), nil
}
// GetApproval gets the approval property value. Approval decisions associated with a request.
// returns a Approvalable when successful
func (m *UserConsentRequest) GetApproval()(Approvalable) {
    val, err := m.GetBackingStore().Get("approval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Approvalable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserConsentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Request.GetFieldDeserializers()
    res["approval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateApprovalFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApproval(val.(Approvalable))
        }
        return nil
    }
    res["reason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReason(val)
        }
        return nil
    }
    return res
}
// GetReason gets the reason property value. The user's justification for requiring access to the app. Supports $filter (eq only) and $orderby.
// returns a *string when successful
func (m *UserConsentRequest) GetReason()(*string) {
    val, err := m.GetBackingStore().Get("reason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserConsentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Request.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("approval", m.GetApproval())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("reason", m.GetReason())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApproval sets the approval property value. Approval decisions associated with a request.
func (m *UserConsentRequest) SetApproval(value Approvalable)() {
    err := m.GetBackingStore().Set("approval", value)
    if err != nil {
        panic(err)
    }
}
// SetReason sets the reason property value. The user's justification for requiring access to the app. Supports $filter (eq only) and $orderby.
func (m *UserConsentRequest) SetReason(value *string)() {
    err := m.GetBackingStore().Set("reason", value)
    if err != nil {
        panic(err)
    }
}
type UserConsentRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Requestable
    GetApproval()(Approvalable)
    GetReason()(*string)
    SetApproval(value Approvalable)()
    SetReason(value *string)()
}
