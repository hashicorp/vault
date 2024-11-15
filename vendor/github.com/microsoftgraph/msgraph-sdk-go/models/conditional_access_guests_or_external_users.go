package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessGuestsOrExternalUsers struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessGuestsOrExternalUsers instantiates a new ConditionalAccessGuestsOrExternalUsers and sets the default values.
func NewConditionalAccessGuestsOrExternalUsers()(*ConditionalAccessGuestsOrExternalUsers) {
    m := &ConditionalAccessGuestsOrExternalUsers{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessGuestsOrExternalUsersFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessGuestsOrExternalUsersFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessGuestsOrExternalUsers(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessGuestsOrExternalUsers) GetAdditionalData()(map[string]any) {
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
func (m *ConditionalAccessGuestsOrExternalUsers) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetExternalTenants gets the externalTenants property value. The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property guestOrExternalUserTypes isn't null or an empty String.
// returns a ConditionalAccessExternalTenantsable when successful
func (m *ConditionalAccessGuestsOrExternalUsers) GetExternalTenants()(ConditionalAccessExternalTenantsable) {
    val, err := m.GetBackingStore().Get("externalTenants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessExternalTenantsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessGuestsOrExternalUsers) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["externalTenants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessExternalTenantsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalTenants(val.(ConditionalAccessExternalTenantsable))
        }
        return nil
    }
    res["guestOrExternalUserTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseConditionalAccessGuestOrExternalUserTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGuestOrExternalUserTypes(val.(*ConditionalAccessGuestOrExternalUserTypes))
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
// GetGuestOrExternalUserTypes gets the guestOrExternalUserTypes property value. The guestOrExternalUserTypes property
// returns a *ConditionalAccessGuestOrExternalUserTypes when successful
func (m *ConditionalAccessGuestsOrExternalUsers) GetGuestOrExternalUserTypes()(*ConditionalAccessGuestOrExternalUserTypes) {
    val, err := m.GetBackingStore().Get("guestOrExternalUserTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ConditionalAccessGuestOrExternalUserTypes)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessGuestsOrExternalUsers) GetOdataType()(*string) {
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
func (m *ConditionalAccessGuestsOrExternalUsers) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("externalTenants", m.GetExternalTenants())
        if err != nil {
            return err
        }
    }
    if m.GetGuestOrExternalUserTypes() != nil {
        cast := (*m.GetGuestOrExternalUserTypes()).String()
        err := writer.WriteStringValue("guestOrExternalUserTypes", &cast)
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
func (m *ConditionalAccessGuestsOrExternalUsers) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessGuestsOrExternalUsers) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetExternalTenants sets the externalTenants property value. The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property guestOrExternalUserTypes isn't null or an empty String.
func (m *ConditionalAccessGuestsOrExternalUsers) SetExternalTenants(value ConditionalAccessExternalTenantsable)() {
    err := m.GetBackingStore().Set("externalTenants", value)
    if err != nil {
        panic(err)
    }
}
// SetGuestOrExternalUserTypes sets the guestOrExternalUserTypes property value. The guestOrExternalUserTypes property
func (m *ConditionalAccessGuestsOrExternalUsers) SetGuestOrExternalUserTypes(value *ConditionalAccessGuestOrExternalUserTypes)() {
    err := m.GetBackingStore().Set("guestOrExternalUserTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessGuestsOrExternalUsers) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessGuestsOrExternalUsersable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetExternalTenants()(ConditionalAccessExternalTenantsable)
    GetGuestOrExternalUserTypes()(*ConditionalAccessGuestOrExternalUserTypes)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetExternalTenants(value ConditionalAccessExternalTenantsable)()
    SetGuestOrExternalUserTypes(value *ConditionalAccessGuestOrExternalUserTypes)()
    SetOdataType(value *string)()
}
