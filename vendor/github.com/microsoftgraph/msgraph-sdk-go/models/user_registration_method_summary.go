package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type UserRegistrationMethodSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewUserRegistrationMethodSummary instantiates a new UserRegistrationMethodSummary and sets the default values.
func NewUserRegistrationMethodSummary()(*UserRegistrationMethodSummary) {
    m := &UserRegistrationMethodSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateUserRegistrationMethodSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserRegistrationMethodSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserRegistrationMethodSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *UserRegistrationMethodSummary) GetAdditionalData()(map[string]any) {
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
func (m *UserRegistrationMethodSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserRegistrationMethodSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["totalUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalUserCount(val)
        }
        return nil
    }
    res["userRegistrationMethodCounts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserRegistrationMethodCountFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserRegistrationMethodCountable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserRegistrationMethodCountable)
                }
            }
            m.SetUserRegistrationMethodCounts(res)
        }
        return nil
    }
    res["userRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseIncludedUserRoles)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserRoles(val.(*IncludedUserRoles))
        }
        return nil
    }
    res["userTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseIncludedUserTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserTypes(val.(*IncludedUserTypes))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *UserRegistrationMethodSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalUserCount gets the totalUserCount property value. Total number of users in the tenant.
// returns a *int64 when successful
func (m *UserRegistrationMethodSummary) GetTotalUserCount()(*int64) {
    val, err := m.GetBackingStore().Get("totalUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUserRegistrationMethodCounts gets the userRegistrationMethodCounts property value. Number of users registered for each authentication method.
// returns a []UserRegistrationMethodCountable when successful
func (m *UserRegistrationMethodSummary) GetUserRegistrationMethodCounts()([]UserRegistrationMethodCountable) {
    val, err := m.GetBackingStore().Get("userRegistrationMethodCounts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserRegistrationMethodCountable)
    }
    return nil
}
// GetUserRoles gets the userRoles property value. The role type of the user. Possible values are: all, privilegedAdmin, admin, user, unknownFutureValue.
// returns a *IncludedUserRoles when successful
func (m *UserRegistrationMethodSummary) GetUserRoles()(*IncludedUserRoles) {
    val, err := m.GetBackingStore().Get("userRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IncludedUserRoles)
    }
    return nil
}
// GetUserTypes gets the userTypes property value. User type. Possible values are: all, member, guest, unknownFutureValue.
// returns a *IncludedUserTypes when successful
func (m *UserRegistrationMethodSummary) GetUserTypes()(*IncludedUserTypes) {
    val, err := m.GetBackingStore().Get("userTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IncludedUserTypes)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserRegistrationMethodSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("totalUserCount", m.GetTotalUserCount())
        if err != nil {
            return err
        }
    }
    if m.GetUserRegistrationMethodCounts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserRegistrationMethodCounts()))
        for i, v := range m.GetUserRegistrationMethodCounts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("userRegistrationMethodCounts", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserRoles() != nil {
        cast := (*m.GetUserRoles()).String()
        err := writer.WriteStringValue("userRoles", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserTypes() != nil {
        cast := (*m.GetUserTypes()).String()
        err := writer.WriteStringValue("userTypes", &cast)
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
func (m *UserRegistrationMethodSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *UserRegistrationMethodSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *UserRegistrationMethodSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalUserCount sets the totalUserCount property value. Total number of users in the tenant.
func (m *UserRegistrationMethodSummary) SetTotalUserCount(value *int64)() {
    err := m.GetBackingStore().Set("totalUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUserRegistrationMethodCounts sets the userRegistrationMethodCounts property value. Number of users registered for each authentication method.
func (m *UserRegistrationMethodSummary) SetUserRegistrationMethodCounts(value []UserRegistrationMethodCountable)() {
    err := m.GetBackingStore().Set("userRegistrationMethodCounts", value)
    if err != nil {
        panic(err)
    }
}
// SetUserRoles sets the userRoles property value. The role type of the user. Possible values are: all, privilegedAdmin, admin, user, unknownFutureValue.
func (m *UserRegistrationMethodSummary) SetUserRoles(value *IncludedUserRoles)() {
    err := m.GetBackingStore().Set("userRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetUserTypes sets the userTypes property value. User type. Possible values are: all, member, guest, unknownFutureValue.
func (m *UserRegistrationMethodSummary) SetUserTypes(value *IncludedUserTypes)() {
    err := m.GetBackingStore().Set("userTypes", value)
    if err != nil {
        panic(err)
    }
}
type UserRegistrationMethodSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetTotalUserCount()(*int64)
    GetUserRegistrationMethodCounts()([]UserRegistrationMethodCountable)
    GetUserRoles()(*IncludedUserRoles)
    GetUserTypes()(*IncludedUserTypes)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetTotalUserCount(value *int64)()
    SetUserRegistrationMethodCounts(value []UserRegistrationMethodCountable)()
    SetUserRoles(value *IncludedUserRoles)()
    SetUserTypes(value *IncludedUserTypes)()
}
