package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AppManagementPolicy struct {
    PolicyBase
}
// NewAppManagementPolicy instantiates a new AppManagementPolicy and sets the default values.
func NewAppManagementPolicy()(*AppManagementPolicy) {
    m := &AppManagementPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.appManagementPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAppManagementPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppManagementPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppManagementPolicy(), nil
}
// GetAppliesTo gets the appliesTo property value. Collection of applications and service principals to which the policy is applied.
// returns a []DirectoryObjectable when successful
func (m *AppManagementPolicy) GetAppliesTo()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("appliesTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AppManagementPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["appliesTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetAppliesTo(res)
        }
        return nil
    }
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
        }
        return nil
    }
    res["restrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCustomAppManagementConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestrictions(val.(CustomAppManagementConfigurationable))
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. Denotes whether the policy is enabled.
// returns a *bool when successful
func (m *AppManagementPolicy) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRestrictions gets the restrictions property value. Restrictions that apply to an application or service principal object.
// returns a CustomAppManagementConfigurationable when successful
func (m *AppManagementPolicy) GetRestrictions()(CustomAppManagementConfigurationable) {
    val, err := m.GetBackingStore().Get("restrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CustomAppManagementConfigurationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AppManagementPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppliesTo() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppliesTo()))
        for i, v := range m.GetAppliesTo() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appliesTo", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("restrictions", m.GetRestrictions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppliesTo sets the appliesTo property value. Collection of applications and service principals to which the policy is applied.
func (m *AppManagementPolicy) SetAppliesTo(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("appliesTo", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Denotes whether the policy is enabled.
func (m *AppManagementPolicy) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetRestrictions sets the restrictions property value. Restrictions that apply to an application or service principal object.
func (m *AppManagementPolicy) SetRestrictions(value CustomAppManagementConfigurationable)() {
    err := m.GetBackingStore().Set("restrictions", value)
    if err != nil {
        panic(err)
    }
}
type AppManagementPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PolicyBaseable
    GetAppliesTo()([]DirectoryObjectable)
    GetIsEnabled()(*bool)
    GetRestrictions()(CustomAppManagementConfigurationable)
    SetAppliesTo(value []DirectoryObjectable)()
    SetIsEnabled(value *bool)()
    SetRestrictions(value CustomAppManagementConfigurationable)()
}
