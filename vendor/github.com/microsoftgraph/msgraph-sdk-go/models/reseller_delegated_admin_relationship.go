package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ResellerDelegatedAdminRelationship struct {
    DelegatedAdminRelationship
}
// NewResellerDelegatedAdminRelationship instantiates a new ResellerDelegatedAdminRelationship and sets the default values.
func NewResellerDelegatedAdminRelationship()(*ResellerDelegatedAdminRelationship) {
    m := &ResellerDelegatedAdminRelationship{
        DelegatedAdminRelationship: *NewDelegatedAdminRelationship(),
    }
    return m
}
// CreateResellerDelegatedAdminRelationshipFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateResellerDelegatedAdminRelationshipFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewResellerDelegatedAdminRelationship(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ResellerDelegatedAdminRelationship) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DelegatedAdminRelationship.GetFieldDeserializers()
    res["indirectProviderTenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIndirectProviderTenantId(val)
        }
        return nil
    }
    res["isPartnerConsentPending"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPartnerConsentPending(val)
        }
        return nil
    }
    return res
}
// GetIndirectProviderTenantId gets the indirectProviderTenantId property value. The tenant ID of the indirect provider partner who created the relationship for the indirect reseller partner.
// returns a *string when successful
func (m *ResellerDelegatedAdminRelationship) GetIndirectProviderTenantId()(*string) {
    val, err := m.GetBackingStore().Get("indirectProviderTenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsPartnerConsentPending gets the isPartnerConsentPending property value. Indicates the indirect reseller partner consent status. true indicates that the partner has yet to review the relationship; false indicates that the partner has already provided consent by approving or rejecting the relationship.
// returns a *bool when successful
func (m *ResellerDelegatedAdminRelationship) GetIsPartnerConsentPending()(*bool) {
    val, err := m.GetBackingStore().Get("isPartnerConsentPending")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ResellerDelegatedAdminRelationship) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DelegatedAdminRelationship.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("indirectProviderTenantId", m.GetIndirectProviderTenantId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isPartnerConsentPending", m.GetIsPartnerConsentPending())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIndirectProviderTenantId sets the indirectProviderTenantId property value. The tenant ID of the indirect provider partner who created the relationship for the indirect reseller partner.
func (m *ResellerDelegatedAdminRelationship) SetIndirectProviderTenantId(value *string)() {
    err := m.GetBackingStore().Set("indirectProviderTenantId", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPartnerConsentPending sets the isPartnerConsentPending property value. Indicates the indirect reseller partner consent status. true indicates that the partner has yet to review the relationship; false indicates that the partner has already provided consent by approving or rejecting the relationship.
func (m *ResellerDelegatedAdminRelationship) SetIsPartnerConsentPending(value *bool)() {
    err := m.GetBackingStore().Set("isPartnerConsentPending", value)
    if err != nil {
        panic(err)
    }
}
type ResellerDelegatedAdminRelationshipable interface {
    DelegatedAdminRelationshipable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIndirectProviderTenantId()(*string)
    GetIsPartnerConsentPending()(*bool)
    SetIndirectProviderTenantId(value *string)()
    SetIsPartnerConsentPending(value *bool)()
}
