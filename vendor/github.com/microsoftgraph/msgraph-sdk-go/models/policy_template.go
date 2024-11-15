package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PolicyTemplate struct {
    Entity
}
// NewPolicyTemplate instantiates a new PolicyTemplate and sets the default values.
func NewPolicyTemplate()(*PolicyTemplate) {
    m := &PolicyTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePolicyTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePolicyTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPolicyTemplate(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PolicyTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["multiTenantOrganizationIdentitySynchronization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMultiTenantOrganizationIdentitySyncPolicyTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultiTenantOrganizationIdentitySynchronization(val.(MultiTenantOrganizationIdentitySyncPolicyTemplateable))
        }
        return nil
    }
    res["multiTenantOrganizationPartnerConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMultiTenantOrganizationPartnerConfigurationTemplateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultiTenantOrganizationPartnerConfiguration(val.(MultiTenantOrganizationPartnerConfigurationTemplateable))
        }
        return nil
    }
    return res
}
// GetMultiTenantOrganizationIdentitySynchronization gets the multiTenantOrganizationIdentitySynchronization property value. Defines an optional cross-tenant access policy template with user synchronization settings for a multitenant organization.
// returns a MultiTenantOrganizationIdentitySyncPolicyTemplateable when successful
func (m *PolicyTemplate) GetMultiTenantOrganizationIdentitySynchronization()(MultiTenantOrganizationIdentitySyncPolicyTemplateable) {
    val, err := m.GetBackingStore().Get("multiTenantOrganizationIdentitySynchronization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MultiTenantOrganizationIdentitySyncPolicyTemplateable)
    }
    return nil
}
// GetMultiTenantOrganizationPartnerConfiguration gets the multiTenantOrganizationPartnerConfiguration property value. Defines an optional cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
// returns a MultiTenantOrganizationPartnerConfigurationTemplateable when successful
func (m *PolicyTemplate) GetMultiTenantOrganizationPartnerConfiguration()(MultiTenantOrganizationPartnerConfigurationTemplateable) {
    val, err := m.GetBackingStore().Get("multiTenantOrganizationPartnerConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MultiTenantOrganizationPartnerConfigurationTemplateable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PolicyTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("multiTenantOrganizationIdentitySynchronization", m.GetMultiTenantOrganizationIdentitySynchronization())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("multiTenantOrganizationPartnerConfiguration", m.GetMultiTenantOrganizationPartnerConfiguration())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMultiTenantOrganizationIdentitySynchronization sets the multiTenantOrganizationIdentitySynchronization property value. Defines an optional cross-tenant access policy template with user synchronization settings for a multitenant organization.
func (m *PolicyTemplate) SetMultiTenantOrganizationIdentitySynchronization(value MultiTenantOrganizationIdentitySyncPolicyTemplateable)() {
    err := m.GetBackingStore().Set("multiTenantOrganizationIdentitySynchronization", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiTenantOrganizationPartnerConfiguration sets the multiTenantOrganizationPartnerConfiguration property value. Defines an optional cross-tenant access policy template with inbound and outbound partner configuration settings for a multitenant organization.
func (m *PolicyTemplate) SetMultiTenantOrganizationPartnerConfiguration(value MultiTenantOrganizationPartnerConfigurationTemplateable)() {
    err := m.GetBackingStore().Set("multiTenantOrganizationPartnerConfiguration", value)
    if err != nil {
        panic(err)
    }
}
type PolicyTemplateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMultiTenantOrganizationIdentitySynchronization()(MultiTenantOrganizationIdentitySyncPolicyTemplateable)
    GetMultiTenantOrganizationPartnerConfiguration()(MultiTenantOrganizationPartnerConfigurationTemplateable)
    SetMultiTenantOrganizationIdentitySynchronization(value MultiTenantOrganizationIdentitySyncPolicyTemplateable)()
    SetMultiTenantOrganizationPartnerConfiguration(value MultiTenantOrganizationPartnerConfigurationTemplateable)()
}
