package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EntitlementManagementSettings struct {
    Entity
}
// NewEntitlementManagementSettings instantiates a new EntitlementManagementSettings and sets the default values.
func NewEntitlementManagementSettings()(*EntitlementManagementSettings) {
    m := &EntitlementManagementSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEntitlementManagementSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEntitlementManagementSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEntitlementManagementSettings(), nil
}
// GetDurationUntilExternalUserDeletedAfterBlocked gets the durationUntilExternalUserDeletedAfterBlocked property value. If externalUserLifecycleAction is blockSignInAndDelete, the duration, typically many days, after an external user is blocked from sign in before their account is deleted.
// returns a *ISODuration when successful
func (m *EntitlementManagementSettings) GetDurationUntilExternalUserDeletedAfterBlocked()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("durationUntilExternalUserDeletedAfterBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetExternalUserLifecycleAction gets the externalUserLifecycleAction property value. Automatic action that the service should take when an external user's last access package assignment is removed. The possible values are: none, blockSignIn, blockSignInAndDelete, unknownFutureValue.
// returns a *AccessPackageExternalUserLifecycleAction when successful
func (m *EntitlementManagementSettings) GetExternalUserLifecycleAction()(*AccessPackageExternalUserLifecycleAction) {
    val, err := m.GetBackingStore().Get("externalUserLifecycleAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageExternalUserLifecycleAction)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EntitlementManagementSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["durationUntilExternalUserDeletedAfterBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationUntilExternalUserDeletedAfterBlocked(val)
        }
        return nil
    }
    res["externalUserLifecycleAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessPackageExternalUserLifecycleAction)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalUserLifecycleAction(val.(*AccessPackageExternalUserLifecycleAction))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *EntitlementManagementSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteISODurationValue("durationUntilExternalUserDeletedAfterBlocked", m.GetDurationUntilExternalUserDeletedAfterBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetExternalUserLifecycleAction() != nil {
        cast := (*m.GetExternalUserLifecycleAction()).String()
        err = writer.WriteStringValue("externalUserLifecycleAction", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDurationUntilExternalUserDeletedAfterBlocked sets the durationUntilExternalUserDeletedAfterBlocked property value. If externalUserLifecycleAction is blockSignInAndDelete, the duration, typically many days, after an external user is blocked from sign in before their account is deleted.
func (m *EntitlementManagementSettings) SetDurationUntilExternalUserDeletedAfterBlocked(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("durationUntilExternalUserDeletedAfterBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalUserLifecycleAction sets the externalUserLifecycleAction property value. Automatic action that the service should take when an external user's last access package assignment is removed. The possible values are: none, blockSignIn, blockSignInAndDelete, unknownFutureValue.
func (m *EntitlementManagementSettings) SetExternalUserLifecycleAction(value *AccessPackageExternalUserLifecycleAction)() {
    err := m.GetBackingStore().Set("externalUserLifecycleAction", value)
    if err != nil {
        panic(err)
    }
}
type EntitlementManagementSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDurationUntilExternalUserDeletedAfterBlocked()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetExternalUserLifecycleAction()(*AccessPackageExternalUserLifecycleAction)
    SetDurationUntilExternalUserDeletedAfterBlocked(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetExternalUserLifecycleAction(value *AccessPackageExternalUserLifecycleAction)()
}
