package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ConditionalAccessUsers struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConditionalAccessUsers instantiates a new ConditionalAccessUsers and sets the default values.
func NewConditionalAccessUsers()(*ConditionalAccessUsers) {
    m := &ConditionalAccessUsers{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConditionalAccessUsersFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConditionalAccessUsersFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConditionalAccessUsers(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConditionalAccessUsers) GetAdditionalData()(map[string]any) {
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
func (m *ConditionalAccessUsers) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetExcludeGroups gets the excludeGroups property value. Group IDs excluded from scope of policy.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetExcludeGroups()([]string) {
    val, err := m.GetBackingStore().Get("excludeGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetExcludeGuestsOrExternalUsers gets the excludeGuestsOrExternalUsers property value. Internal guests or external users excluded from the policy scope. Optionally populated.
// returns a ConditionalAccessGuestsOrExternalUsersable when successful
func (m *ConditionalAccessUsers) GetExcludeGuestsOrExternalUsers()(ConditionalAccessGuestsOrExternalUsersable) {
    val, err := m.GetBackingStore().Get("excludeGuestsOrExternalUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessGuestsOrExternalUsersable)
    }
    return nil
}
// GetExcludeRoles gets the excludeRoles property value. Role IDs excluded from scope of policy.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetExcludeRoles()([]string) {
    val, err := m.GetBackingStore().Get("excludeRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetExcludeUsers gets the excludeUsers property value. User IDs excluded from scope of policy and/or GuestsOrExternalUsers.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetExcludeUsers()([]string) {
    val, err := m.GetBackingStore().Get("excludeUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConditionalAccessUsers) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["excludeGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetExcludeGroups(res)
        }
        return nil
    }
    res["excludeGuestsOrExternalUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessGuestsOrExternalUsersFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExcludeGuestsOrExternalUsers(val.(ConditionalAccessGuestsOrExternalUsersable))
        }
        return nil
    }
    res["excludeRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetExcludeRoles(res)
        }
        return nil
    }
    res["excludeUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetExcludeUsers(res)
        }
        return nil
    }
    res["includeGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetIncludeGroups(res)
        }
        return nil
    }
    res["includeGuestsOrExternalUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateConditionalAccessGuestsOrExternalUsersFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncludeGuestsOrExternalUsers(val.(ConditionalAccessGuestsOrExternalUsersable))
        }
        return nil
    }
    res["includeRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetIncludeRoles(res)
        }
        return nil
    }
    res["includeUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetIncludeUsers(res)
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
// GetIncludeGroups gets the includeGroups property value. Group IDs in scope of policy unless explicitly excluded.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetIncludeGroups()([]string) {
    val, err := m.GetBackingStore().Get("includeGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetIncludeGuestsOrExternalUsers gets the includeGuestsOrExternalUsers property value. Internal guests or external users included in the policy scope. Optionally populated.
// returns a ConditionalAccessGuestsOrExternalUsersable when successful
func (m *ConditionalAccessUsers) GetIncludeGuestsOrExternalUsers()(ConditionalAccessGuestsOrExternalUsersable) {
    val, err := m.GetBackingStore().Get("includeGuestsOrExternalUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ConditionalAccessGuestsOrExternalUsersable)
    }
    return nil
}
// GetIncludeRoles gets the includeRoles property value. Role IDs in scope of policy unless explicitly excluded.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetIncludeRoles()([]string) {
    val, err := m.GetBackingStore().Get("includeRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetIncludeUsers gets the includeUsers property value. User IDs in scope of policy unless explicitly excluded, None, All, or GuestsOrExternalUsers.
// returns a []string when successful
func (m *ConditionalAccessUsers) GetIncludeUsers()([]string) {
    val, err := m.GetBackingStore().Get("includeUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConditionalAccessUsers) GetOdataType()(*string) {
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
func (m *ConditionalAccessUsers) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetExcludeGroups() != nil {
        err := writer.WriteCollectionOfStringValues("excludeGroups", m.GetExcludeGroups())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("excludeGuestsOrExternalUsers", m.GetExcludeGuestsOrExternalUsers())
        if err != nil {
            return err
        }
    }
    if m.GetExcludeRoles() != nil {
        err := writer.WriteCollectionOfStringValues("excludeRoles", m.GetExcludeRoles())
        if err != nil {
            return err
        }
    }
    if m.GetExcludeUsers() != nil {
        err := writer.WriteCollectionOfStringValues("excludeUsers", m.GetExcludeUsers())
        if err != nil {
            return err
        }
    }
    if m.GetIncludeGroups() != nil {
        err := writer.WriteCollectionOfStringValues("includeGroups", m.GetIncludeGroups())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("includeGuestsOrExternalUsers", m.GetIncludeGuestsOrExternalUsers())
        if err != nil {
            return err
        }
    }
    if m.GetIncludeRoles() != nil {
        err := writer.WriteCollectionOfStringValues("includeRoles", m.GetIncludeRoles())
        if err != nil {
            return err
        }
    }
    if m.GetIncludeUsers() != nil {
        err := writer.WriteCollectionOfStringValues("includeUsers", m.GetIncludeUsers())
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
func (m *ConditionalAccessUsers) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConditionalAccessUsers) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetExcludeGroups sets the excludeGroups property value. Group IDs excluded from scope of policy.
func (m *ConditionalAccessUsers) SetExcludeGroups(value []string)() {
    err := m.GetBackingStore().Set("excludeGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetExcludeGuestsOrExternalUsers sets the excludeGuestsOrExternalUsers property value. Internal guests or external users excluded from the policy scope. Optionally populated.
func (m *ConditionalAccessUsers) SetExcludeGuestsOrExternalUsers(value ConditionalAccessGuestsOrExternalUsersable)() {
    err := m.GetBackingStore().Set("excludeGuestsOrExternalUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetExcludeRoles sets the excludeRoles property value. Role IDs excluded from scope of policy.
func (m *ConditionalAccessUsers) SetExcludeRoles(value []string)() {
    err := m.GetBackingStore().Set("excludeRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetExcludeUsers sets the excludeUsers property value. User IDs excluded from scope of policy and/or GuestsOrExternalUsers.
func (m *ConditionalAccessUsers) SetExcludeUsers(value []string)() {
    err := m.GetBackingStore().Set("excludeUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeGroups sets the includeGroups property value. Group IDs in scope of policy unless explicitly excluded.
func (m *ConditionalAccessUsers) SetIncludeGroups(value []string)() {
    err := m.GetBackingStore().Set("includeGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeGuestsOrExternalUsers sets the includeGuestsOrExternalUsers property value. Internal guests or external users included in the policy scope. Optionally populated.
func (m *ConditionalAccessUsers) SetIncludeGuestsOrExternalUsers(value ConditionalAccessGuestsOrExternalUsersable)() {
    err := m.GetBackingStore().Set("includeGuestsOrExternalUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeRoles sets the includeRoles property value. Role IDs in scope of policy unless explicitly excluded.
func (m *ConditionalAccessUsers) SetIncludeRoles(value []string)() {
    err := m.GetBackingStore().Set("includeRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludeUsers sets the includeUsers property value. User IDs in scope of policy unless explicitly excluded, None, All, or GuestsOrExternalUsers.
func (m *ConditionalAccessUsers) SetIncludeUsers(value []string)() {
    err := m.GetBackingStore().Set("includeUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConditionalAccessUsers) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type ConditionalAccessUsersable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetExcludeGroups()([]string)
    GetExcludeGuestsOrExternalUsers()(ConditionalAccessGuestsOrExternalUsersable)
    GetExcludeRoles()([]string)
    GetExcludeUsers()([]string)
    GetIncludeGroups()([]string)
    GetIncludeGuestsOrExternalUsers()(ConditionalAccessGuestsOrExternalUsersable)
    GetIncludeRoles()([]string)
    GetIncludeUsers()([]string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetExcludeGroups(value []string)()
    SetExcludeGuestsOrExternalUsers(value ConditionalAccessGuestsOrExternalUsersable)()
    SetExcludeRoles(value []string)()
    SetExcludeUsers(value []string)()
    SetIncludeGroups(value []string)()
    SetIncludeGuestsOrExternalUsers(value ConditionalAccessGuestsOrExternalUsersable)()
    SetIncludeRoles(value []string)()
    SetIncludeUsers(value []string)()
    SetOdataType(value *string)()
}
