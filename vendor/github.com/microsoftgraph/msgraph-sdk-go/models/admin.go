package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Admin struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAdmin instantiates a new Admin and sets the default values.
func NewAdmin()(*Admin) {
    m := &Admin{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAdminFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAdminFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdmin(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Admin) GetAdditionalData()(map[string]any) {
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
func (m *Admin) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEdge gets the edge property value. A container for Microsoft Edge resources. Read-only.
// returns a Edgeable when successful
func (m *Admin) GetEdge()(Edgeable) {
    val, err := m.GetBackingStore().Get("edge")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Edgeable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Admin) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["edge"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdgeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdge(val.(Edgeable))
        }
        return nil
    }
    res["microsoft365Apps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAdminMicrosoft365AppsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoft365Apps(val.(AdminMicrosoft365Appsable))
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
    res["people"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePeopleAdminSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPeople(val.(PeopleAdminSettingsable))
        }
        return nil
    }
    res["reportSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAdminReportSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReportSettings(val.(AdminReportSettingsable))
        }
        return nil
    }
    res["serviceAnnouncement"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateServiceAnnouncementFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceAnnouncement(val.(ServiceAnnouncementable))
        }
        return nil
    }
    res["sharepoint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharepointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharepoint(val.(Sharepointable))
        }
        return nil
    }
    return res
}
// GetMicrosoft365Apps gets the microsoft365Apps property value. A container for the Microsoft 365 apps admin functionality.
// returns a AdminMicrosoft365Appsable when successful
func (m *Admin) GetMicrosoft365Apps()(AdminMicrosoft365Appsable) {
    val, err := m.GetBackingStore().Get("microsoft365Apps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AdminMicrosoft365Appsable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Admin) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPeople gets the people property value. Represents a setting to control people-related admin settings in the tenant.
// returns a PeopleAdminSettingsable when successful
func (m *Admin) GetPeople()(PeopleAdminSettingsable) {
    val, err := m.GetBackingStore().Get("people")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PeopleAdminSettingsable)
    }
    return nil
}
// GetReportSettings gets the reportSettings property value. A container for administrative resources to manage reports.
// returns a AdminReportSettingsable when successful
func (m *Admin) GetReportSettings()(AdminReportSettingsable) {
    val, err := m.GetBackingStore().Get("reportSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AdminReportSettingsable)
    }
    return nil
}
// GetServiceAnnouncement gets the serviceAnnouncement property value. A container for service communications resources. Read-only.
// returns a ServiceAnnouncementable when successful
func (m *Admin) GetServiceAnnouncement()(ServiceAnnouncementable) {
    val, err := m.GetBackingStore().Get("serviceAnnouncement")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ServiceAnnouncementable)
    }
    return nil
}
// GetSharepoint gets the sharepoint property value. The sharepoint property
// returns a Sharepointable when successful
func (m *Admin) GetSharepoint()(Sharepointable) {
    val, err := m.GetBackingStore().Get("sharepoint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Sharepointable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Admin) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("edge", m.GetEdge())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("microsoft365Apps", m.GetMicrosoft365Apps())
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
        err := writer.WriteObjectValue("people", m.GetPeople())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("reportSettings", m.GetReportSettings())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("serviceAnnouncement", m.GetServiceAnnouncement())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("sharepoint", m.GetSharepoint())
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
func (m *Admin) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Admin) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEdge sets the edge property value. A container for Microsoft Edge resources. Read-only.
func (m *Admin) SetEdge(value Edgeable)() {
    err := m.GetBackingStore().Set("edge", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoft365Apps sets the microsoft365Apps property value. A container for the Microsoft 365 apps admin functionality.
func (m *Admin) SetMicrosoft365Apps(value AdminMicrosoft365Appsable)() {
    err := m.GetBackingStore().Set("microsoft365Apps", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Admin) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPeople sets the people property value. Represents a setting to control people-related admin settings in the tenant.
func (m *Admin) SetPeople(value PeopleAdminSettingsable)() {
    err := m.GetBackingStore().Set("people", value)
    if err != nil {
        panic(err)
    }
}
// SetReportSettings sets the reportSettings property value. A container for administrative resources to manage reports.
func (m *Admin) SetReportSettings(value AdminReportSettingsable)() {
    err := m.GetBackingStore().Set("reportSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceAnnouncement sets the serviceAnnouncement property value. A container for service communications resources. Read-only.
func (m *Admin) SetServiceAnnouncement(value ServiceAnnouncementable)() {
    err := m.GetBackingStore().Set("serviceAnnouncement", value)
    if err != nil {
        panic(err)
    }
}
// SetSharepoint sets the sharepoint property value. The sharepoint property
func (m *Admin) SetSharepoint(value Sharepointable)() {
    err := m.GetBackingStore().Set("sharepoint", value)
    if err != nil {
        panic(err)
    }
}
type Adminable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEdge()(Edgeable)
    GetMicrosoft365Apps()(AdminMicrosoft365Appsable)
    GetOdataType()(*string)
    GetPeople()(PeopleAdminSettingsable)
    GetReportSettings()(AdminReportSettingsable)
    GetServiceAnnouncement()(ServiceAnnouncementable)
    GetSharepoint()(Sharepointable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEdge(value Edgeable)()
    SetMicrosoft365Apps(value AdminMicrosoft365Appsable)()
    SetOdataType(value *string)()
    SetPeople(value PeopleAdminSettingsable)()
    SetReportSettings(value AdminReportSettingsable)()
    SetServiceAnnouncement(value ServiceAnnouncementable)()
    SetSharepoint(value Sharepointable)()
}
