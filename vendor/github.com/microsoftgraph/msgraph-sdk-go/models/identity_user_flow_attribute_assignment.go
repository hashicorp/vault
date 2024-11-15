package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type IdentityUserFlowAttributeAssignment struct {
    Entity
}
// NewIdentityUserFlowAttributeAssignment instantiates a new IdentityUserFlowAttributeAssignment and sets the default values.
func NewIdentityUserFlowAttributeAssignment()(*IdentityUserFlowAttributeAssignment) {
    m := &IdentityUserFlowAttributeAssignment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateIdentityUserFlowAttributeAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityUserFlowAttributeAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityUserFlowAttributeAssignment(), nil
}
// GetDisplayName gets the displayName property value. The display name of the identityUserFlowAttribute within a user flow.
// returns a *string when successful
func (m *IdentityUserFlowAttributeAssignment) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IdentityUserFlowAttributeAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["isOptional"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOptional(val)
        }
        return nil
    }
    res["requiresVerification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequiresVerification(val)
        }
        return nil
    }
    res["userAttribute"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentityUserFlowAttributeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAttribute(val.(IdentityUserFlowAttributeable))
        }
        return nil
    }
    res["userAttributeValues"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserAttributeValuesItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserAttributeValuesItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserAttributeValuesItemable)
                }
            }
            m.SetUserAttributeValues(res)
        }
        return nil
    }
    res["userInputType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseIdentityUserFlowAttributeInputType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserInputType(val.(*IdentityUserFlowAttributeInputType))
        }
        return nil
    }
    return res
}
// GetIsOptional gets the isOptional property value. Determines whether the identityUserFlowAttribute is optional. true means the user doesn't have to provide a value. false means the user can't complete sign-up without providing a value.
// returns a *bool when successful
func (m *IdentityUserFlowAttributeAssignment) GetIsOptional()(*bool) {
    val, err := m.GetBackingStore().Get("isOptional")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRequiresVerification gets the requiresVerification property value. Determines whether the identityUserFlowAttribute requires verification, and is only used for verifying the user's phone number or email address.
// returns a *bool when successful
func (m *IdentityUserFlowAttributeAssignment) GetRequiresVerification()(*bool) {
    val, err := m.GetBackingStore().Get("requiresVerification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUserAttribute gets the userAttribute property value. The user attribute that you want to add to your user flow.
// returns a IdentityUserFlowAttributeable when successful
func (m *IdentityUserFlowAttributeAssignment) GetUserAttribute()(IdentityUserFlowAttributeable) {
    val, err := m.GetBackingStore().Get("userAttribute")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentityUserFlowAttributeable)
    }
    return nil
}
// GetUserAttributeValues gets the userAttributeValues property value. The input options for the user flow attribute. Only applicable when the userInputType is radioSingleSelect, dropdownSingleSelect, or checkboxMultiSelect.
// returns a []UserAttributeValuesItemable when successful
func (m *IdentityUserFlowAttributeAssignment) GetUserAttributeValues()([]UserAttributeValuesItemable) {
    val, err := m.GetBackingStore().Get("userAttributeValues")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserAttributeValuesItemable)
    }
    return nil
}
// GetUserInputType gets the userInputType property value. The userInputType property
// returns a *IdentityUserFlowAttributeInputType when successful
func (m *IdentityUserFlowAttributeAssignment) GetUserInputType()(*IdentityUserFlowAttributeInputType) {
    val, err := m.GetBackingStore().Get("userInputType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*IdentityUserFlowAttributeInputType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityUserFlowAttributeAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOptional", m.GetIsOptional())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("requiresVerification", m.GetRequiresVerification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("userAttribute", m.GetUserAttribute())
        if err != nil {
            return err
        }
    }
    if m.GetUserAttributeValues() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserAttributeValues()))
        for i, v := range m.GetUserAttributeValues() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userAttributeValues", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserInputType() != nil {
        cast := (*m.GetUserInputType()).String()
        err = writer.WriteStringValue("userInputType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The display name of the identityUserFlowAttribute within a user flow.
func (m *IdentityUserFlowAttributeAssignment) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOptional sets the isOptional property value. Determines whether the identityUserFlowAttribute is optional. true means the user doesn't have to provide a value. false means the user can't complete sign-up without providing a value.
func (m *IdentityUserFlowAttributeAssignment) SetIsOptional(value *bool)() {
    err := m.GetBackingStore().Set("isOptional", value)
    if err != nil {
        panic(err)
    }
}
// SetRequiresVerification sets the requiresVerification property value. Determines whether the identityUserFlowAttribute requires verification, and is only used for verifying the user's phone number or email address.
func (m *IdentityUserFlowAttributeAssignment) SetRequiresVerification(value *bool)() {
    err := m.GetBackingStore().Set("requiresVerification", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAttribute sets the userAttribute property value. The user attribute that you want to add to your user flow.
func (m *IdentityUserFlowAttributeAssignment) SetUserAttribute(value IdentityUserFlowAttributeable)() {
    err := m.GetBackingStore().Set("userAttribute", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAttributeValues sets the userAttributeValues property value. The input options for the user flow attribute. Only applicable when the userInputType is radioSingleSelect, dropdownSingleSelect, or checkboxMultiSelect.
func (m *IdentityUserFlowAttributeAssignment) SetUserAttributeValues(value []UserAttributeValuesItemable)() {
    err := m.GetBackingStore().Set("userAttributeValues", value)
    if err != nil {
        panic(err)
    }
}
// SetUserInputType sets the userInputType property value. The userInputType property
func (m *IdentityUserFlowAttributeAssignment) SetUserInputType(value *IdentityUserFlowAttributeInputType)() {
    err := m.GetBackingStore().Set("userInputType", value)
    if err != nil {
        panic(err)
    }
}
type IdentityUserFlowAttributeAssignmentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetIsOptional()(*bool)
    GetRequiresVerification()(*bool)
    GetUserAttribute()(IdentityUserFlowAttributeable)
    GetUserAttributeValues()([]UserAttributeValuesItemable)
    GetUserInputType()(*IdentityUserFlowAttributeInputType)
    SetDisplayName(value *string)()
    SetIsOptional(value *bool)()
    SetRequiresVerification(value *bool)()
    SetUserAttribute(value IdentityUserFlowAttributeable)()
    SetUserAttributeValues(value []UserAttributeValuesItemable)()
    SetUserInputType(value *IdentityUserFlowAttributeInputType)()
}
