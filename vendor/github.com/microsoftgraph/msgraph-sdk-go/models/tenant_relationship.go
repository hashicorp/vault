package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type TenantRelationship struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewTenantRelationship instantiates a new TenantRelationship and sets the default values.
func NewTenantRelationship()(*TenantRelationship) {
    m := &TenantRelationship{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateTenantRelationshipFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTenantRelationshipFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTenantRelationship(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *TenantRelationship) GetAdditionalData()(map[string]any) {
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
func (m *TenantRelationship) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDelegatedAdminCustomers gets the delegatedAdminCustomers property value. The customer who has a delegated admin relationship with a Microsoft partner.
// returns a []DelegatedAdminCustomerable when successful
func (m *TenantRelationship) GetDelegatedAdminCustomers()([]DelegatedAdminCustomerable) {
    val, err := m.GetBackingStore().Get("delegatedAdminCustomers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DelegatedAdminCustomerable)
    }
    return nil
}
// GetDelegatedAdminRelationships gets the delegatedAdminRelationships property value. The details of the delegated administrative privileges that a Microsoft partner has in a customer tenant.
// returns a []DelegatedAdminRelationshipable when successful
func (m *TenantRelationship) GetDelegatedAdminRelationships()([]DelegatedAdminRelationshipable) {
    val, err := m.GetBackingStore().Get("delegatedAdminRelationships")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DelegatedAdminRelationshipable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TenantRelationship) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["delegatedAdminCustomers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDelegatedAdminCustomerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DelegatedAdminCustomerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DelegatedAdminCustomerable)
                }
            }
            m.SetDelegatedAdminCustomers(res)
        }
        return nil
    }
    res["delegatedAdminRelationships"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDelegatedAdminRelationshipFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DelegatedAdminRelationshipable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DelegatedAdminRelationshipable)
                }
            }
            m.SetDelegatedAdminRelationships(res)
        }
        return nil
    }
    res["multiTenantOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMultiTenantOrganizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMultiTenantOrganization(val.(MultiTenantOrganizationable))
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
// GetMultiTenantOrganization gets the multiTenantOrganization property value. Defines an organization with more than one instance of Microsoft Entra ID.
// returns a MultiTenantOrganizationable when successful
func (m *TenantRelationship) GetMultiTenantOrganization()(MultiTenantOrganizationable) {
    val, err := m.GetBackingStore().Get("multiTenantOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MultiTenantOrganizationable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *TenantRelationship) GetOdataType()(*string) {
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
func (m *TenantRelationship) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDelegatedAdminCustomers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDelegatedAdminCustomers()))
        for i, v := range m.GetDelegatedAdminCustomers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("delegatedAdminCustomers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDelegatedAdminRelationships() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDelegatedAdminRelationships()))
        for i, v := range m.GetDelegatedAdminRelationships() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("delegatedAdminRelationships", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("multiTenantOrganization", m.GetMultiTenantOrganization())
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
func (m *TenantRelationship) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *TenantRelationship) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDelegatedAdminCustomers sets the delegatedAdminCustomers property value. The customer who has a delegated admin relationship with a Microsoft partner.
func (m *TenantRelationship) SetDelegatedAdminCustomers(value []DelegatedAdminCustomerable)() {
    err := m.GetBackingStore().Set("delegatedAdminCustomers", value)
    if err != nil {
        panic(err)
    }
}
// SetDelegatedAdminRelationships sets the delegatedAdminRelationships property value. The details of the delegated administrative privileges that a Microsoft partner has in a customer tenant.
func (m *TenantRelationship) SetDelegatedAdminRelationships(value []DelegatedAdminRelationshipable)() {
    err := m.GetBackingStore().Set("delegatedAdminRelationships", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiTenantOrganization sets the multiTenantOrganization property value. Defines an organization with more than one instance of Microsoft Entra ID.
func (m *TenantRelationship) SetMultiTenantOrganization(value MultiTenantOrganizationable)() {
    err := m.GetBackingStore().Set("multiTenantOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *TenantRelationship) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type TenantRelationshipable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDelegatedAdminCustomers()([]DelegatedAdminCustomerable)
    GetDelegatedAdminRelationships()([]DelegatedAdminRelationshipable)
    GetMultiTenantOrganization()(MultiTenantOrganizationable)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDelegatedAdminCustomers(value []DelegatedAdminCustomerable)()
    SetDelegatedAdminRelationships(value []DelegatedAdminRelationshipable)()
    SetMultiTenantOrganization(value MultiTenantOrganizationable)()
    SetOdataType(value *string)()
}
