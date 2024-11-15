package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type CopyNotebookModel struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCopyNotebookModel instantiates a new CopyNotebookModel and sets the default values.
func NewCopyNotebookModel()(*CopyNotebookModel) {
    m := &CopyNotebookModel{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCopyNotebookModelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCopyNotebookModelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCopyNotebookModel(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CopyNotebookModel) GetAdditionalData()(map[string]any) {
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
func (m *CopyNotebookModel) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCreatedBy gets the createdBy property value. The createdBy property
// returns a *string when successful
func (m *CopyNotebookModel) GetCreatedBy()(*string) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedByIdentity gets the createdByIdentity property value. The createdByIdentity property
// returns a IdentitySetable when successful
func (m *CopyNotebookModel) GetCreatedByIdentity()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdByIdentity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedTime gets the createdTime property value. The createdTime property
// returns a *Time when successful
func (m *CopyNotebookModel) GetCreatedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CopyNotebookModel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val)
        }
        return nil
    }
    res["createdByIdentity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedByIdentity(val.(IdentitySetable))
        }
        return nil
    }
    res["createdTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedTime(val)
        }
        return nil
    }
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["isDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDefault(val)
        }
        return nil
    }
    res["isShared"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsShared(val)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val)
        }
        return nil
    }
    res["lastModifiedByIdentity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedByIdentity(val.(IdentitySetable))
        }
        return nil
    }
    res["lastModifiedTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedTime(val)
        }
        return nil
    }
    res["links"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateNotebookLinksFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLinks(val.(NotebookLinksable))
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
    res["sectionGroupsUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSectionGroupsUrl(val)
        }
        return nil
    }
    res["sectionsUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSectionsUrl(val)
        }
        return nil
    }
    res["self"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSelf(val)
        }
        return nil
    }
    res["userRole"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOnenoteUserRole)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserRole(val.(*OnenoteUserRole))
        }
        return nil
    }
    return res
}
// GetId gets the id property value. The id property
// returns a *string when successful
func (m *CopyNotebookModel) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsDefault gets the isDefault property value. The isDefault property
// returns a *bool when successful
func (m *CopyNotebookModel) GetIsDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsShared gets the isShared property value. The isShared property
// returns a *bool when successful
func (m *CopyNotebookModel) GetIsShared()(*bool) {
    val, err := m.GetBackingStore().Get("isShared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The lastModifiedBy property
// returns a *string when successful
func (m *CopyNotebookModel) GetLastModifiedBy()(*string) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedByIdentity gets the lastModifiedByIdentity property value. The lastModifiedByIdentity property
// returns a IdentitySetable when successful
func (m *CopyNotebookModel) GetLastModifiedByIdentity()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedByIdentity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedTime gets the lastModifiedTime property value. The lastModifiedTime property
// returns a *Time when successful
func (m *CopyNotebookModel) GetLastModifiedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLinks gets the links property value. The links property
// returns a NotebookLinksable when successful
func (m *CopyNotebookModel) GetLinks()(NotebookLinksable) {
    val, err := m.GetBackingStore().Get("links")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(NotebookLinksable)
    }
    return nil
}
// GetName gets the name property value. The name property
// returns a *string when successful
func (m *CopyNotebookModel) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *CopyNotebookModel) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSectionGroupsUrl gets the sectionGroupsUrl property value. The sectionGroupsUrl property
// returns a *string when successful
func (m *CopyNotebookModel) GetSectionGroupsUrl()(*string) {
    val, err := m.GetBackingStore().Get("sectionGroupsUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSectionsUrl gets the sectionsUrl property value. The sectionsUrl property
// returns a *string when successful
func (m *CopyNotebookModel) GetSectionsUrl()(*string) {
    val, err := m.GetBackingStore().Get("sectionsUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSelf gets the self property value. The self property
// returns a *string when successful
func (m *CopyNotebookModel) GetSelf()(*string) {
    val, err := m.GetBackingStore().Get("self")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUserRole gets the userRole property value. The userRole property
// returns a *OnenoteUserRole when successful
func (m *CopyNotebookModel) GetUserRole()(*OnenoteUserRole) {
    val, err := m.GetBackingStore().Get("userRole")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OnenoteUserRole)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CopyNotebookModel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("createdByIdentity", m.GetCreatedByIdentity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("createdTime", m.GetCreatedTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isDefault", m.GetIsDefault())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isShared", m.GetIsShared())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("lastModifiedByIdentity", m.GetLastModifiedByIdentity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("lastModifiedTime", m.GetLastModifiedTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("links", m.GetLinks())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
        err := writer.WriteStringValue("sectionGroupsUrl", m.GetSectionGroupsUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("sectionsUrl", m.GetSectionsUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("self", m.GetSelf())
        if err != nil {
            return err
        }
    }
    if m.GetUserRole() != nil {
        cast := (*m.GetUserRole()).String()
        err := writer.WriteStringValue("userRole", &cast)
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
func (m *CopyNotebookModel) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CopyNotebookModel) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCreatedBy sets the createdBy property value. The createdBy property
func (m *CopyNotebookModel) SetCreatedBy(value *string)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedByIdentity sets the createdByIdentity property value. The createdByIdentity property
func (m *CopyNotebookModel) SetCreatedByIdentity(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdByIdentity", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedTime sets the createdTime property value. The createdTime property
func (m *CopyNotebookModel) SetCreatedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdTime", value)
    if err != nil {
        panic(err)
    }
}
// SetId sets the id property value. The id property
func (m *CopyNotebookModel) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDefault sets the isDefault property value. The isDefault property
func (m *CopyNotebookModel) SetIsDefault(value *bool)() {
    err := m.GetBackingStore().Set("isDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetIsShared sets the isShared property value. The isShared property
func (m *CopyNotebookModel) SetIsShared(value *bool)() {
    err := m.GetBackingStore().Set("isShared", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The lastModifiedBy property
func (m *CopyNotebookModel) SetLastModifiedBy(value *string)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedByIdentity sets the lastModifiedByIdentity property value. The lastModifiedByIdentity property
func (m *CopyNotebookModel) SetLastModifiedByIdentity(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedByIdentity", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedTime sets the lastModifiedTime property value. The lastModifiedTime property
func (m *CopyNotebookModel) SetLastModifiedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLinks sets the links property value. The links property
func (m *CopyNotebookModel) SetLinks(value NotebookLinksable)() {
    err := m.GetBackingStore().Set("links", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name property
func (m *CopyNotebookModel) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *CopyNotebookModel) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionGroupsUrl sets the sectionGroupsUrl property value. The sectionGroupsUrl property
func (m *CopyNotebookModel) SetSectionGroupsUrl(value *string)() {
    err := m.GetBackingStore().Set("sectionGroupsUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetSectionsUrl sets the sectionsUrl property value. The sectionsUrl property
func (m *CopyNotebookModel) SetSectionsUrl(value *string)() {
    err := m.GetBackingStore().Set("sectionsUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetSelf sets the self property value. The self property
func (m *CopyNotebookModel) SetSelf(value *string)() {
    err := m.GetBackingStore().Set("self", value)
    if err != nil {
        panic(err)
    }
}
// SetUserRole sets the userRole property value. The userRole property
func (m *CopyNotebookModel) SetUserRole(value *OnenoteUserRole)() {
    err := m.GetBackingStore().Set("userRole", value)
    if err != nil {
        panic(err)
    }
}
type CopyNotebookModelable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCreatedBy()(*string)
    GetCreatedByIdentity()(IdentitySetable)
    GetCreatedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetId()(*string)
    GetIsDefault()(*bool)
    GetIsShared()(*bool)
    GetLastModifiedBy()(*string)
    GetLastModifiedByIdentity()(IdentitySetable)
    GetLastModifiedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLinks()(NotebookLinksable)
    GetName()(*string)
    GetOdataType()(*string)
    GetSectionGroupsUrl()(*string)
    GetSectionsUrl()(*string)
    GetSelf()(*string)
    GetUserRole()(*OnenoteUserRole)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCreatedBy(value *string)()
    SetCreatedByIdentity(value IdentitySetable)()
    SetCreatedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetId(value *string)()
    SetIsDefault(value *bool)()
    SetIsShared(value *bool)()
    SetLastModifiedBy(value *string)()
    SetLastModifiedByIdentity(value IdentitySetable)()
    SetLastModifiedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLinks(value NotebookLinksable)()
    SetName(value *string)()
    SetOdataType(value *string)()
    SetSectionGroupsUrl(value *string)()
    SetSectionsUrl(value *string)()
    SetSelf(value *string)()
    SetUserRole(value *OnenoteUserRole)()
}
