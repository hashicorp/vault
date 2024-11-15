package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TenantAppManagementPolicy struct {
    PolicyBase
}
// NewTenantAppManagementPolicy instantiates a new TenantAppManagementPolicy and sets the default values.
func NewTenantAppManagementPolicy()(*TenantAppManagementPolicy) {
    m := &TenantAppManagementPolicy{
        PolicyBase: *NewPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.tenantAppManagementPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTenantAppManagementPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTenantAppManagementPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTenantAppManagementPolicy(), nil
}
// GetApplicationRestrictions gets the applicationRestrictions property value. Restrictions that apply as default to all application objects in the tenant.
// returns a AppManagementApplicationConfigurationable when successful
func (m *TenantAppManagementPolicy) GetApplicationRestrictions()(AppManagementApplicationConfigurationable) {
    val, err := m.GetBackingStore().Get("applicationRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AppManagementApplicationConfigurationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TenantAppManagementPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PolicyBase.GetFieldDeserializers()
    res["applicationRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAppManagementApplicationConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicationRestrictions(val.(AppManagementApplicationConfigurationable))
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
    res["servicePrincipalRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAppManagementServicePrincipalConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalRestrictions(val.(AppManagementServicePrincipalConfigurationable))
        }
        return nil
    }
    return res
}
// GetIsEnabled gets the isEnabled property value. Denotes whether the policy is enabled. Default value is false.
// returns a *bool when successful
func (m *TenantAppManagementPolicy) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetServicePrincipalRestrictions gets the servicePrincipalRestrictions property value. Restrictions that apply as default to all service principal objects in the tenant.
// returns a AppManagementServicePrincipalConfigurationable when successful
func (m *TenantAppManagementPolicy) GetServicePrincipalRestrictions()(AppManagementServicePrincipalConfigurationable) {
    val, err := m.GetBackingStore().Get("servicePrincipalRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AppManagementServicePrincipalConfigurationable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TenantAppManagementPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("applicationRestrictions", m.GetApplicationRestrictions())
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
        err = writer.WriteObjectValue("servicePrincipalRestrictions", m.GetServicePrincipalRestrictions())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicationRestrictions sets the applicationRestrictions property value. Restrictions that apply as default to all application objects in the tenant.
func (m *TenantAppManagementPolicy) SetApplicationRestrictions(value AppManagementApplicationConfigurationable)() {
    err := m.GetBackingStore().Set("applicationRestrictions", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Denotes whether the policy is enabled. Default value is false.
func (m *TenantAppManagementPolicy) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalRestrictions sets the servicePrincipalRestrictions property value. Restrictions that apply as default to all service principal objects in the tenant.
func (m *TenantAppManagementPolicy) SetServicePrincipalRestrictions(value AppManagementServicePrincipalConfigurationable)() {
    err := m.GetBackingStore().Set("servicePrincipalRestrictions", value)
    if err != nil {
        panic(err)
    }
}
type TenantAppManagementPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PolicyBaseable
    GetApplicationRestrictions()(AppManagementApplicationConfigurationable)
    GetIsEnabled()(*bool)
    GetServicePrincipalRestrictions()(AppManagementServicePrincipalConfigurationable)
    SetApplicationRestrictions(value AppManagementApplicationConfigurationable)()
    SetIsEnabled(value *bool)()
    SetServicePrincipalRestrictions(value AppManagementServicePrincipalConfigurationable)()
}
