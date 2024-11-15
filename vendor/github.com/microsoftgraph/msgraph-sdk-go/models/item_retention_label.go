package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemRetentionLabel struct {
    Entity
}
// NewItemRetentionLabel instantiates a new ItemRetentionLabel and sets the default values.
func NewItemRetentionLabel()(*ItemRetentionLabel) {
    m := &ItemRetentionLabel{
        Entity: *NewEntity(),
    }
    return m
}
// CreateItemRetentionLabelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemRetentionLabelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemRetentionLabel(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemRetentionLabel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["isLabelAppliedExplicitly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLabelAppliedExplicitly(val)
        }
        return nil
    }
    res["labelAppliedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabelAppliedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["labelAppliedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabelAppliedDateTime(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["retentionSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRetentionLabelSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRetentionSettings(val.(RetentionLabelSettingsable))
        }
        return nil
    }
    return res
}
// GetIsLabelAppliedExplicitly gets the isLabelAppliedExplicitly property value. Specifies whether the label is applied explicitly on the item. True indicates that the label is applied explicitly; otherwise, the label is inherited from its parent. Read-only.
// returns a *bool when successful
func (m *ItemRetentionLabel) GetIsLabelAppliedExplicitly()(*bool) {
    val, err := m.GetBackingStore().Get("isLabelAppliedExplicitly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLabelAppliedBy gets the labelAppliedBy property value. Identity of the user who applied the label. Read-only.
// returns a IdentitySetable when successful
func (m *ItemRetentionLabel) GetLabelAppliedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("labelAppliedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLabelAppliedDateTime gets the labelAppliedDateTime property value. The date and time when the label was applied on the item. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *ItemRetentionLabel) GetLabelAppliedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("labelAppliedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetName gets the name property value. The retention label on the document. Read-write.
// returns a *string when successful
func (m *ItemRetentionLabel) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRetentionSettings gets the retentionSettings property value. The retention settings enforced on the item. Read-write.
// returns a RetentionLabelSettingsable when successful
func (m *ItemRetentionLabel) GetRetentionSettings()(RetentionLabelSettingsable) {
    val, err := m.GetBackingStore().Get("retentionSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RetentionLabelSettingsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemRetentionLabel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isLabelAppliedExplicitly", m.GetIsLabelAppliedExplicitly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("labelAppliedBy", m.GetLabelAppliedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("labelAppliedDateTime", m.GetLabelAppliedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("retentionSettings", m.GetRetentionSettings())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsLabelAppliedExplicitly sets the isLabelAppliedExplicitly property value. Specifies whether the label is applied explicitly on the item. True indicates that the label is applied explicitly; otherwise, the label is inherited from its parent. Read-only.
func (m *ItemRetentionLabel) SetIsLabelAppliedExplicitly(value *bool)() {
    err := m.GetBackingStore().Set("isLabelAppliedExplicitly", value)
    if err != nil {
        panic(err)
    }
}
// SetLabelAppliedBy sets the labelAppliedBy property value. Identity of the user who applied the label. Read-only.
func (m *ItemRetentionLabel) SetLabelAppliedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("labelAppliedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLabelAppliedDateTime sets the labelAppliedDateTime property value. The date and time when the label was applied on the item. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *ItemRetentionLabel) SetLabelAppliedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("labelAppliedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The retention label on the document. Read-write.
func (m *ItemRetentionLabel) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetRetentionSettings sets the retentionSettings property value. The retention settings enforced on the item. Read-write.
func (m *ItemRetentionLabel) SetRetentionSettings(value RetentionLabelSettingsable)() {
    err := m.GetBackingStore().Set("retentionSettings", value)
    if err != nil {
        panic(err)
    }
}
type ItemRetentionLabelable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsLabelAppliedExplicitly()(*bool)
    GetLabelAppliedBy()(IdentitySetable)
    GetLabelAppliedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetName()(*string)
    GetRetentionSettings()(RetentionLabelSettingsable)
    SetIsLabelAppliedExplicitly(value *bool)()
    SetLabelAppliedBy(value IdentitySetable)()
    SetLabelAppliedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetName(value *string)()
    SetRetentionSettings(value RetentionLabelSettingsable)()
}
