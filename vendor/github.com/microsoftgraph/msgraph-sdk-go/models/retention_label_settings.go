package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type RetentionLabelSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRetentionLabelSettings instantiates a new RetentionLabelSettings and sets the default values.
func NewRetentionLabelSettings()(*RetentionLabelSettings) {
    m := &RetentionLabelSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRetentionLabelSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRetentionLabelSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRetentionLabelSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *RetentionLabelSettings) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *RetentionLabelSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RetentionLabelSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["isContentUpdateAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsContentUpdateAllowed(val)
        }
        return nil
    }
    res["isDeleteAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDeleteAllowed(val)
        }
        return nil
    }
    res["isLabelUpdateAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLabelUpdateAllowed(val)
        }
        return nil
    }
    res["isMetadataUpdateAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMetadataUpdateAllowed(val)
        }
        return nil
    }
    res["isRecordLocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRecordLocked(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    return res
}
// GetIsContentUpdateAllowed gets the isContentUpdateAllowed property value. Specifies whether updates to document content are allowed. Read-only.
// returns a *bool when successful
func (m *RetentionLabelSettings) GetIsContentUpdateAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("isContentUpdateAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDeleteAllowed gets the isDeleteAllowed property value. Specifies whether the document deletion is allowed. Read-only.
// returns a *bool when successful
func (m *RetentionLabelSettings) GetIsDeleteAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("isDeleteAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsLabelUpdateAllowed gets the isLabelUpdateAllowed property value. Specifies whether you're allowed to change the retention label on the document. Read-only.
// returns a *bool when successful
func (m *RetentionLabelSettings) GetIsLabelUpdateAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("isLabelUpdateAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMetadataUpdateAllowed gets the isMetadataUpdateAllowed property value. Specifies whether updates to the item metadata (for example, the Title field) are blocked. Read-only.
// returns a *bool when successful
func (m *RetentionLabelSettings) GetIsMetadataUpdateAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("isMetadataUpdateAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRecordLocked gets the isRecordLocked property value. Specifies whether the item is locked. Read-write.
// returns a *bool when successful
func (m *RetentionLabelSettings) GetIsRecordLocked()(*bool) {
    val, err := m.GetBackingStore().Get("isRecordLocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *RetentionLabelSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RetentionLabelSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("isContentUpdateAllowed", m.GetIsContentUpdateAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isDeleteAllowed", m.GetIsDeleteAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isLabelUpdateAllowed", m.GetIsLabelUpdateAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isMetadataUpdateAllowed", m.GetIsMetadataUpdateAllowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isRecordLocked", m.GetIsRecordLocked())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *RetentionLabelSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *RetentionLabelSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsContentUpdateAllowed sets the isContentUpdateAllowed property value. Specifies whether updates to document content are allowed. Read-only.
func (m *RetentionLabelSettings) SetIsContentUpdateAllowed(value *bool)() {
    err := m.GetBackingStore().Set("isContentUpdateAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDeleteAllowed sets the isDeleteAllowed property value. Specifies whether the document deletion is allowed. Read-only.
func (m *RetentionLabelSettings) SetIsDeleteAllowed(value *bool)() {
    err := m.GetBackingStore().Set("isDeleteAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetIsLabelUpdateAllowed sets the isLabelUpdateAllowed property value. Specifies whether you're allowed to change the retention label on the document. Read-only.
func (m *RetentionLabelSettings) SetIsLabelUpdateAllowed(value *bool)() {
    err := m.GetBackingStore().Set("isLabelUpdateAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMetadataUpdateAllowed sets the isMetadataUpdateAllowed property value. Specifies whether updates to the item metadata (for example, the Title field) are blocked. Read-only.
func (m *RetentionLabelSettings) SetIsMetadataUpdateAllowed(value *bool)() {
    err := m.GetBackingStore().Set("isMetadataUpdateAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRecordLocked sets the isRecordLocked property value. Specifies whether the item is locked. Read-write.
func (m *RetentionLabelSettings) SetIsRecordLocked(value *bool)() {
    err := m.GetBackingStore().Set("isRecordLocked", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RetentionLabelSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type RetentionLabelSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsContentUpdateAllowed()(*bool)
    GetIsDeleteAllowed()(*bool)
    GetIsLabelUpdateAllowed()(*bool)
    GetIsMetadataUpdateAllowed()(*bool)
    GetIsRecordLocked()(*bool)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsContentUpdateAllowed(value *bool)()
    SetIsDeleteAllowed(value *bool)()
    SetIsLabelUpdateAllowed(value *bool)()
    SetIsMetadataUpdateAllowed(value *bool)()
    SetIsRecordLocked(value *bool)()
    SetOdataType(value *string)()
}
