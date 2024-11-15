package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type GroupLifecyclePolicy struct {
    Entity
}
// NewGroupLifecyclePolicy instantiates a new GroupLifecyclePolicy and sets the default values.
func NewGroupLifecyclePolicy()(*GroupLifecyclePolicy) {
    m := &GroupLifecyclePolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateGroupLifecyclePolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateGroupLifecyclePolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewGroupLifecyclePolicy(), nil
}
// GetAlternateNotificationEmails gets the alternateNotificationEmails property value. List of email address to send notifications for groups without owners. Multiple email address can be defined by separating email address with a semicolon.
// returns a *string when successful
func (m *GroupLifecyclePolicy) GetAlternateNotificationEmails()(*string) {
    val, err := m.GetBackingStore().Get("alternateNotificationEmails")
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
func (m *GroupLifecyclePolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["alternateNotificationEmails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlternateNotificationEmails(val)
        }
        return nil
    }
    res["groupLifetimeInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupLifetimeInDays(val)
        }
        return nil
    }
    res["managedGroupTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedGroupTypes(val)
        }
        return nil
    }
    return res
}
// GetGroupLifetimeInDays gets the groupLifetimeInDays property value. Number of days before a group expires and needs to be renewed. Once renewed, the group expiration is extended by the number of days defined.
// returns a *int32 when successful
func (m *GroupLifecyclePolicy) GetGroupLifetimeInDays()(*int32) {
    val, err := m.GetBackingStore().Get("groupLifetimeInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetManagedGroupTypes gets the managedGroupTypes property value. The group type for which the expiration policy applies. Possible values are All, Selected or None.
// returns a *string when successful
func (m *GroupLifecyclePolicy) GetManagedGroupTypes()(*string) {
    val, err := m.GetBackingStore().Get("managedGroupTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *GroupLifecyclePolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("alternateNotificationEmails", m.GetAlternateNotificationEmails())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("groupLifetimeInDays", m.GetGroupLifetimeInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managedGroupTypes", m.GetManagedGroupTypes())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAlternateNotificationEmails sets the alternateNotificationEmails property value. List of email address to send notifications for groups without owners. Multiple email address can be defined by separating email address with a semicolon.
func (m *GroupLifecyclePolicy) SetAlternateNotificationEmails(value *string)() {
    err := m.GetBackingStore().Set("alternateNotificationEmails", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupLifetimeInDays sets the groupLifetimeInDays property value. Number of days before a group expires and needs to be renewed. Once renewed, the group expiration is extended by the number of days defined.
func (m *GroupLifecyclePolicy) SetGroupLifetimeInDays(value *int32)() {
    err := m.GetBackingStore().Set("groupLifetimeInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedGroupTypes sets the managedGroupTypes property value. The group type for which the expiration policy applies. Possible values are All, Selected or None.
func (m *GroupLifecyclePolicy) SetManagedGroupTypes(value *string)() {
    err := m.GetBackingStore().Set("managedGroupTypes", value)
    if err != nil {
        panic(err)
    }
}
type GroupLifecyclePolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlternateNotificationEmails()(*string)
    GetGroupLifetimeInDays()(*int32)
    GetManagedGroupTypes()(*string)
    SetAlternateNotificationEmails(value *string)()
    SetGroupLifetimeInDays(value *int32)()
    SetManagedGroupTypes(value *string)()
}
