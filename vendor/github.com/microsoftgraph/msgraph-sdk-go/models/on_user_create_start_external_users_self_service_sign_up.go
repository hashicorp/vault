package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OnUserCreateStartExternalUsersSelfServiceSignUp struct {
    OnUserCreateStartHandler
}
// NewOnUserCreateStartExternalUsersSelfServiceSignUp instantiates a new OnUserCreateStartExternalUsersSelfServiceSignUp and sets the default values.
func NewOnUserCreateStartExternalUsersSelfServiceSignUp()(*OnUserCreateStartExternalUsersSelfServiceSignUp) {
    m := &OnUserCreateStartExternalUsersSelfServiceSignUp{
        OnUserCreateStartHandler: *NewOnUserCreateStartHandler(),
    }
    odataTypeValue := "#microsoft.graph.onUserCreateStartExternalUsersSelfServiceSignUp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOnUserCreateStartExternalUsersSelfServiceSignUpFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnUserCreateStartExternalUsersSelfServiceSignUpFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnUserCreateStartExternalUsersSelfServiceSignUp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnUserCreateStartExternalUsersSelfServiceSignUp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OnUserCreateStartHandler.GetFieldDeserializers()
    res["userTypeToCreate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUserType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserTypeToCreate(val.(*UserType))
        }
        return nil
    }
    return res
}
// GetUserTypeToCreate gets the userTypeToCreate property value. The type of user to create. Maps to userType property of user object. The possible values are: member, guest, unknownFutureValue.
// returns a *UserType when successful
func (m *OnUserCreateStartExternalUsersSelfServiceSignUp) GetUserTypeToCreate()(*UserType) {
    val, err := m.GetBackingStore().Get("userTypeToCreate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UserType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OnUserCreateStartExternalUsersSelfServiceSignUp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OnUserCreateStartHandler.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetUserTypeToCreate() != nil {
        cast := (*m.GetUserTypeToCreate()).String()
        err = writer.WriteStringValue("userTypeToCreate", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUserTypeToCreate sets the userTypeToCreate property value. The type of user to create. Maps to userType property of user object. The possible values are: member, guest, unknownFutureValue.
func (m *OnUserCreateStartExternalUsersSelfServiceSignUp) SetUserTypeToCreate(value *UserType)() {
    err := m.GetBackingStore().Set("userTypeToCreate", value)
    if err != nil {
        panic(err)
    }
}
type OnUserCreateStartExternalUsersSelfServiceSignUpable interface {
    OnUserCreateStartHandlerable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUserTypeToCreate()(*UserType)
    SetUserTypeToCreate(value *UserType)()
}
