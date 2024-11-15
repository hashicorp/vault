package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessPackageAssignmentRequestorSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessPackageAssignmentRequestorSettings instantiates a new AccessPackageAssignmentRequestorSettings and sets the default values.
func NewAccessPackageAssignmentRequestorSettings()(*AccessPackageAssignmentRequestorSettings) {
    m := &AccessPackageAssignmentRequestorSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessPackageAssignmentRequestorSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentRequestorSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentRequestorSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessPackageAssignmentRequestorSettings) GetAdditionalData()(map[string]any) {
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
// GetAllowCustomAssignmentSchedule gets the allowCustomAssignmentSchedule property value. False indicates that the requestor isn't permitted to include a schedule in their request.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetAllowCustomAssignmentSchedule()(*bool) {
    val, err := m.GetBackingStore().Get("allowCustomAssignmentSchedule")
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
func (m *AccessPackageAssignmentRequestorSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetEnableOnBehalfRequestorsToAddAccess gets the enableOnBehalfRequestorsToAddAccess property value. True allows on-behalf-of requestors to create a request to add access for another principal.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableOnBehalfRequestorsToAddAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableOnBehalfRequestorsToAddAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableOnBehalfRequestorsToRemoveAccess gets the enableOnBehalfRequestorsToRemoveAccess property value. True allows on-behalf-of requestors to create a request to remove access for another principal.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableOnBehalfRequestorsToRemoveAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableOnBehalfRequestorsToRemoveAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableOnBehalfRequestorsToUpdateAccess gets the enableOnBehalfRequestorsToUpdateAccess property value. True allows on-behalf-of requestors to create a request to update access for another principal.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableOnBehalfRequestorsToUpdateAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableOnBehalfRequestorsToUpdateAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableTargetsToSelfAddAccess gets the enableTargetsToSelfAddAccess property value. True allows requestors to create a request to add access for themselves.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableTargetsToSelfAddAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableTargetsToSelfAddAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableTargetsToSelfRemoveAccess gets the enableTargetsToSelfRemoveAccess property value. True allows requestors to create a request to remove their access.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableTargetsToSelfRemoveAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableTargetsToSelfRemoveAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnableTargetsToSelfUpdateAccess gets the enableTargetsToSelfUpdateAccess property value. True allows requestors to create a request to update their access.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestorSettings) GetEnableTargetsToSelfUpdateAccess()(*bool) {
    val, err := m.GetBackingStore().Get("enableTargetsToSelfUpdateAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAssignmentRequestorSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowCustomAssignmentSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCustomAssignmentSchedule(val)
        }
        return nil
    }
    res["enableOnBehalfRequestorsToAddAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableOnBehalfRequestorsToAddAccess(val)
        }
        return nil
    }
    res["enableOnBehalfRequestorsToRemoveAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableOnBehalfRequestorsToRemoveAccess(val)
        }
        return nil
    }
    res["enableOnBehalfRequestorsToUpdateAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableOnBehalfRequestorsToUpdateAccess(val)
        }
        return nil
    }
    res["enableTargetsToSelfAddAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableTargetsToSelfAddAccess(val)
        }
        return nil
    }
    res["enableTargetsToSelfRemoveAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableTargetsToSelfRemoveAccess(val)
        }
        return nil
    }
    res["enableTargetsToSelfUpdateAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnableTargetsToSelfUpdateAccess(val)
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
    res["onBehalfRequestors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubjectSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SubjectSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SubjectSetable)
                }
            }
            m.SetOnBehalfRequestors(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AccessPackageAssignmentRequestorSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOnBehalfRequestors gets the onBehalfRequestors property value. The principals who can request on-behalf-of others.
// returns a []SubjectSetable when successful
func (m *AccessPackageAssignmentRequestorSettings) GetOnBehalfRequestors()([]SubjectSetable) {
    val, err := m.GetBackingStore().Get("onBehalfRequestors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SubjectSetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentRequestorSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowCustomAssignmentSchedule", m.GetAllowCustomAssignmentSchedule())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableOnBehalfRequestorsToAddAccess", m.GetEnableOnBehalfRequestorsToAddAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableOnBehalfRequestorsToRemoveAccess", m.GetEnableOnBehalfRequestorsToRemoveAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableOnBehalfRequestorsToUpdateAccess", m.GetEnableOnBehalfRequestorsToUpdateAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableTargetsToSelfAddAccess", m.GetEnableTargetsToSelfAddAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableTargetsToSelfRemoveAccess", m.GetEnableTargetsToSelfRemoveAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enableTargetsToSelfUpdateAccess", m.GetEnableTargetsToSelfUpdateAccess())
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
    if m.GetOnBehalfRequestors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOnBehalfRequestors()))
        for i, v := range m.GetOnBehalfRequestors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("onBehalfRequestors", cast)
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
func (m *AccessPackageAssignmentRequestorSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCustomAssignmentSchedule sets the allowCustomAssignmentSchedule property value. False indicates that the requestor isn't permitted to include a schedule in their request.
func (m *AccessPackageAssignmentRequestorSettings) SetAllowCustomAssignmentSchedule(value *bool)() {
    err := m.GetBackingStore().Set("allowCustomAssignmentSchedule", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessPackageAssignmentRequestorSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetEnableOnBehalfRequestorsToAddAccess sets the enableOnBehalfRequestorsToAddAccess property value. True allows on-behalf-of requestors to create a request to add access for another principal.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableOnBehalfRequestorsToAddAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableOnBehalfRequestorsToAddAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableOnBehalfRequestorsToRemoveAccess sets the enableOnBehalfRequestorsToRemoveAccess property value. True allows on-behalf-of requestors to create a request to remove access for another principal.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableOnBehalfRequestorsToRemoveAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableOnBehalfRequestorsToRemoveAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableOnBehalfRequestorsToUpdateAccess sets the enableOnBehalfRequestorsToUpdateAccess property value. True allows on-behalf-of requestors to create a request to update access for another principal.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableOnBehalfRequestorsToUpdateAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableOnBehalfRequestorsToUpdateAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableTargetsToSelfAddAccess sets the enableTargetsToSelfAddAccess property value. True allows requestors to create a request to add access for themselves.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableTargetsToSelfAddAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableTargetsToSelfAddAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableTargetsToSelfRemoveAccess sets the enableTargetsToSelfRemoveAccess property value. True allows requestors to create a request to remove their access.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableTargetsToSelfRemoveAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableTargetsToSelfRemoveAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetEnableTargetsToSelfUpdateAccess sets the enableTargetsToSelfUpdateAccess property value. True allows requestors to create a request to update their access.
func (m *AccessPackageAssignmentRequestorSettings) SetEnableTargetsToSelfUpdateAccess(value *bool)() {
    err := m.GetBackingStore().Set("enableTargetsToSelfUpdateAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageAssignmentRequestorSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOnBehalfRequestors sets the onBehalfRequestors property value. The principals who can request on-behalf-of others.
func (m *AccessPackageAssignmentRequestorSettings) SetOnBehalfRequestors(value []SubjectSetable)() {
    err := m.GetBackingStore().Set("onBehalfRequestors", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentRequestorSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowCustomAssignmentSchedule()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetEnableOnBehalfRequestorsToAddAccess()(*bool)
    GetEnableOnBehalfRequestorsToRemoveAccess()(*bool)
    GetEnableOnBehalfRequestorsToUpdateAccess()(*bool)
    GetEnableTargetsToSelfAddAccess()(*bool)
    GetEnableTargetsToSelfRemoveAccess()(*bool)
    GetEnableTargetsToSelfUpdateAccess()(*bool)
    GetOdataType()(*string)
    GetOnBehalfRequestors()([]SubjectSetable)
    SetAllowCustomAssignmentSchedule(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetEnableOnBehalfRequestorsToAddAccess(value *bool)()
    SetEnableOnBehalfRequestorsToRemoveAccess(value *bool)()
    SetEnableOnBehalfRequestorsToUpdateAccess(value *bool)()
    SetEnableTargetsToSelfAddAccess(value *bool)()
    SetEnableTargetsToSelfRemoveAccess(value *bool)()
    SetEnableTargetsToSelfUpdateAccess(value *bool)()
    SetOdataType(value *string)()
    SetOnBehalfRequestors(value []SubjectSetable)()
}
