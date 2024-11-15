package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemActivity struct {
    Entity
}
// NewItemActivity instantiates a new ItemActivity and sets the default values.
func NewItemActivity()(*ItemActivity) {
    m := &ItemActivity{
        Entity: *NewEntity(),
    }
    return m
}
// CreateItemActivityFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemActivityFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemActivity(), nil
}
// GetAccess gets the access property value. An item was accessed.
// returns a AccessActionable when successful
func (m *ItemActivity) GetAccess()(AccessActionable) {
    val, err := m.GetBackingStore().Get("access")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessActionable)
    }
    return nil
}
// GetActivityDateTime gets the activityDateTime property value. Details about when the activity took place. Read-only.
// returns a *Time when successful
func (m *ItemActivity) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("activityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetActor gets the actor property value. Identity of who performed the action. Read-only.
// returns a IdentitySetable when successful
func (m *ItemActivity) GetActor()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("actor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetDriveItem gets the driveItem property value. Exposes the driveItem that was the target of this activity.
// returns a DriveItemable when successful
func (m *ItemActivity) GetDriveItem()(DriveItemable) {
    val, err := m.GetBackingStore().Get("driveItem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemActivity) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["access"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccess(val.(AccessActionable))
        }
        return nil
    }
    res["activityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityDateTime(val)
        }
        return nil
    }
    res["actor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActor(val.(IdentitySetable))
        }
        return nil
    }
    res["driveItem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriveItem(val.(DriveItemable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ItemActivity) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("access", m.GetAccess())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("actor", m.GetActor())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("driveItem", m.GetDriveItem())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccess sets the access property value. An item was accessed.
func (m *ItemActivity) SetAccess(value AccessActionable)() {
    err := m.GetBackingStore().Set("access", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityDateTime sets the activityDateTime property value. Details about when the activity took place. Read-only.
func (m *ItemActivity) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("activityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetActor sets the actor property value. Identity of who performed the action. Read-only.
func (m *ItemActivity) SetActor(value IdentitySetable)() {
    err := m.GetBackingStore().Set("actor", value)
    if err != nil {
        panic(err)
    }
}
// SetDriveItem sets the driveItem property value. Exposes the driveItem that was the target of this activity.
func (m *ItemActivity) SetDriveItem(value DriveItemable)() {
    err := m.GetBackingStore().Set("driveItem", value)
    if err != nil {
        panic(err)
    }
}
type ItemActivityable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccess()(AccessActionable)
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetActor()(IdentitySetable)
    GetDriveItem()(DriveItemable)
    SetAccess(value AccessActionable)()
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetActor(value IdentitySetable)()
    SetDriveItem(value DriveItemable)()
}
