package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MultiTenantOrganizationIdentitySyncPolicyTemplate struct {
    Entity
}
// NewMultiTenantOrganizationIdentitySyncPolicyTemplate instantiates a new MultiTenantOrganizationIdentitySyncPolicyTemplate and sets the default values.
func NewMultiTenantOrganizationIdentitySyncPolicyTemplate()(*MultiTenantOrganizationIdentitySyncPolicyTemplate) {
    m := &MultiTenantOrganizationIdentitySyncPolicyTemplate{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMultiTenantOrganizationIdentitySyncPolicyTemplateFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMultiTenantOrganizationIdentitySyncPolicyTemplateFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMultiTenantOrganizationIdentitySyncPolicyTemplate(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["templateApplicationLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTemplateApplicationLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTemplateApplicationLevel(val.(*TemplateApplicationLevel))
        }
        return nil
    }
    res["userSyncInbound"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCrossTenantUserSyncInboundFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserSyncInbound(val.(CrossTenantUserSyncInboundable))
        }
        return nil
    }
    return res
}
// GetTemplateApplicationLevel gets the templateApplicationLevel property value. The templateApplicationLevel property
// returns a *TemplateApplicationLevel when successful
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) GetTemplateApplicationLevel()(*TemplateApplicationLevel) {
    val, err := m.GetBackingStore().Get("templateApplicationLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TemplateApplicationLevel)
    }
    return nil
}
// GetUserSyncInbound gets the userSyncInbound property value. Defines whether users can be synchronized from the partner tenant.
// returns a CrossTenantUserSyncInboundable when successful
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) GetUserSyncInbound()(CrossTenantUserSyncInboundable) {
    val, err := m.GetBackingStore().Get("userSyncInbound")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CrossTenantUserSyncInboundable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetTemplateApplicationLevel() != nil {
        cast := (*m.GetTemplateApplicationLevel()).String()
        err = writer.WriteStringValue("templateApplicationLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("userSyncInbound", m.GetUserSyncInbound())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTemplateApplicationLevel sets the templateApplicationLevel property value. The templateApplicationLevel property
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) SetTemplateApplicationLevel(value *TemplateApplicationLevel)() {
    err := m.GetBackingStore().Set("templateApplicationLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetUserSyncInbound sets the userSyncInbound property value. Defines whether users can be synchronized from the partner tenant.
func (m *MultiTenantOrganizationIdentitySyncPolicyTemplate) SetUserSyncInbound(value CrossTenantUserSyncInboundable)() {
    err := m.GetBackingStore().Set("userSyncInbound", value)
    if err != nil {
        panic(err)
    }
}
type MultiTenantOrganizationIdentitySyncPolicyTemplateable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetTemplateApplicationLevel()(*TemplateApplicationLevel)
    GetUserSyncInbound()(CrossTenantUserSyncInboundable)
    SetTemplateApplicationLevel(value *TemplateApplicationLevel)()
    SetUserSyncInbound(value CrossTenantUserSyncInboundable)()
}
