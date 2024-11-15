package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DelegatedAdminCustomer struct {
    Entity
}
// NewDelegatedAdminCustomer instantiates a new DelegatedAdminCustomer and sets the default values.
func NewDelegatedAdminCustomer()(*DelegatedAdminCustomer) {
    m := &DelegatedAdminCustomer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDelegatedAdminCustomerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDelegatedAdminCustomerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDelegatedAdminCustomer(), nil
}
// GetDisplayName gets the displayName property value. The Microsoft Entra ID display name of the customer tenant. Read-only. Supports $orderby.
// returns a *string when successful
func (m *DelegatedAdminCustomer) GetDisplayName()(*string) {
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
func (m *DelegatedAdminCustomer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["serviceManagementDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDelegatedAdminServiceManagementDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DelegatedAdminServiceManagementDetailable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DelegatedAdminServiceManagementDetailable)
                }
            }
            m.SetServiceManagementDetails(res)
        }
        return nil
    }
    res["tenantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantId(val)
        }
        return nil
    }
    return res
}
// GetServiceManagementDetails gets the serviceManagementDetails property value. Contains the management details of a service in the customer tenant that's managed by delegated administration.
// returns a []DelegatedAdminServiceManagementDetailable when successful
func (m *DelegatedAdminCustomer) GetServiceManagementDetails()([]DelegatedAdminServiceManagementDetailable) {
    val, err := m.GetBackingStore().Get("serviceManagementDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DelegatedAdminServiceManagementDetailable)
    }
    return nil
}
// GetTenantId gets the tenantId property value. The Microsoft Entra ID-assigned tenant ID of the customer. Read-only.
// returns a *string when successful
func (m *DelegatedAdminCustomer) GetTenantId()(*string) {
    val, err := m.GetBackingStore().Get("tenantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DelegatedAdminCustomer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetServiceManagementDetails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServiceManagementDetails()))
        for i, v := range m.GetServiceManagementDetails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("serviceManagementDetails", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("tenantId", m.GetTenantId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The Microsoft Entra ID display name of the customer tenant. Read-only. Supports $orderby.
func (m *DelegatedAdminCustomer) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceManagementDetails sets the serviceManagementDetails property value. Contains the management details of a service in the customer tenant that's managed by delegated administration.
func (m *DelegatedAdminCustomer) SetServiceManagementDetails(value []DelegatedAdminServiceManagementDetailable)() {
    err := m.GetBackingStore().Set("serviceManagementDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantId sets the tenantId property value. The Microsoft Entra ID-assigned tenant ID of the customer. Read-only.
func (m *DelegatedAdminCustomer) SetTenantId(value *string)() {
    err := m.GetBackingStore().Set("tenantId", value)
    if err != nil {
        panic(err)
    }
}
type DelegatedAdminCustomerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetServiceManagementDetails()([]DelegatedAdminServiceManagementDetailable)
    GetTenantId()(*string)
    SetDisplayName(value *string)()
    SetServiceManagementDetails(value []DelegatedAdminServiceManagementDetailable)()
    SetTenantId(value *string)()
}
