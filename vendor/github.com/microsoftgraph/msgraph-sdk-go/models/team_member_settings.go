package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TeamMemberSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTeamMemberSettings instantiates a new TeamMemberSettings and sets the default values.
func NewTeamMemberSettings()(*TeamMemberSettings) {
    m := &TeamMemberSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTeamMemberSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamMemberSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamMemberSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TeamMemberSettings) GetAdditionalData()(map[string]any) {
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
// GetAllowAddRemoveApps gets the allowAddRemoveApps property value. If set to true, members can add and remove apps.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowAddRemoveApps()(*bool) {
    val, err := m.GetBackingStore().Get("allowAddRemoveApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowCreatePrivateChannels gets the allowCreatePrivateChannels property value. If set to true, members can add and update private channels.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowCreatePrivateChannels()(*bool) {
    val, err := m.GetBackingStore().Get("allowCreatePrivateChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowCreateUpdateChannels gets the allowCreateUpdateChannels property value. If set to true, members can add and update channels.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowCreateUpdateChannels()(*bool) {
    val, err := m.GetBackingStore().Get("allowCreateUpdateChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowCreateUpdateRemoveConnectors gets the allowCreateUpdateRemoveConnectors property value. If set to true, members can add, update, and remove connectors.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowCreateUpdateRemoveConnectors()(*bool) {
    val, err := m.GetBackingStore().Get("allowCreateUpdateRemoveConnectors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowCreateUpdateRemoveTabs gets the allowCreateUpdateRemoveTabs property value. If set to true, members can add, update, and remove tabs.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowCreateUpdateRemoveTabs()(*bool) {
    val, err := m.GetBackingStore().Get("allowCreateUpdateRemoveTabs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowDeleteChannels gets the allowDeleteChannels property value. If set to true, members can delete channels.
// returns a *bool when successful
func (m *TeamMemberSettings) GetAllowDeleteChannels()(*bool) {
    val, err := m.GetBackingStore().Get("allowDeleteChannels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *TeamMemberSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamMemberSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowAddRemoveApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowAddRemoveApps(val)
        }
        return nil
    }
    res["allowCreatePrivateChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCreatePrivateChannels(val)
        }
        return nil
    }
    res["allowCreateUpdateChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCreateUpdateChannels(val)
        }
        return nil
    }
    res["allowCreateUpdateRemoveConnectors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCreateUpdateRemoveConnectors(val)
        }
        return nil
    }
    res["allowCreateUpdateRemoveTabs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCreateUpdateRemoveTabs(val)
        }
        return nil
    }
    res["allowDeleteChannels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowDeleteChannels(val)
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
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TeamMemberSettings) GetOdataType()(*string) {
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
func (m *TeamMemberSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowAddRemoveApps", m.GetAllowAddRemoveApps())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowCreatePrivateChannels", m.GetAllowCreatePrivateChannels())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowCreateUpdateChannels", m.GetAllowCreateUpdateChannels())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowCreateUpdateRemoveConnectors", m.GetAllowCreateUpdateRemoveConnectors())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowCreateUpdateRemoveTabs", m.GetAllowCreateUpdateRemoveTabs())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowDeleteChannels", m.GetAllowDeleteChannels())
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
func (m *TeamMemberSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowAddRemoveApps sets the allowAddRemoveApps property value. If set to true, members can add and remove apps.
func (m *TeamMemberSettings) SetAllowAddRemoveApps(value *bool)() {
    err := m.GetBackingStore().Set("allowAddRemoveApps", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCreatePrivateChannels sets the allowCreatePrivateChannels property value. If set to true, members can add and update private channels.
func (m *TeamMemberSettings) SetAllowCreatePrivateChannels(value *bool)() {
    err := m.GetBackingStore().Set("allowCreatePrivateChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCreateUpdateChannels sets the allowCreateUpdateChannels property value. If set to true, members can add and update channels.
func (m *TeamMemberSettings) SetAllowCreateUpdateChannels(value *bool)() {
    err := m.GetBackingStore().Set("allowCreateUpdateChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCreateUpdateRemoveConnectors sets the allowCreateUpdateRemoveConnectors property value. If set to true, members can add, update, and remove connectors.
func (m *TeamMemberSettings) SetAllowCreateUpdateRemoveConnectors(value *bool)() {
    err := m.GetBackingStore().Set("allowCreateUpdateRemoveConnectors", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCreateUpdateRemoveTabs sets the allowCreateUpdateRemoveTabs property value. If set to true, members can add, update, and remove tabs.
func (m *TeamMemberSettings) SetAllowCreateUpdateRemoveTabs(value *bool)() {
    err := m.GetBackingStore().Set("allowCreateUpdateRemoveTabs", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowDeleteChannels sets the allowDeleteChannels property value. If set to true, members can delete channels.
func (m *TeamMemberSettings) SetAllowDeleteChannels(value *bool)() {
    err := m.GetBackingStore().Set("allowDeleteChannels", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TeamMemberSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TeamMemberSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type TeamMemberSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowAddRemoveApps()(*bool)
    GetAllowCreatePrivateChannels()(*bool)
    GetAllowCreateUpdateChannels()(*bool)
    GetAllowCreateUpdateRemoveConnectors()(*bool)
    GetAllowCreateUpdateRemoveTabs()(*bool)
    GetAllowDeleteChannels()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    SetAllowAddRemoveApps(value *bool)()
    SetAllowCreatePrivateChannels(value *bool)()
    SetAllowCreateUpdateChannels(value *bool)()
    SetAllowCreateUpdateRemoveConnectors(value *bool)()
    SetAllowCreateUpdateRemoveTabs(value *bool)()
    SetAllowDeleteChannels(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
}
