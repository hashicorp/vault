package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DefaultUserRolePermissions struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDefaultUserRolePermissions instantiates a new DefaultUserRolePermissions and sets the default values.
func NewDefaultUserRolePermissions()(*DefaultUserRolePermissions) {
    m := &DefaultUserRolePermissions{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDefaultUserRolePermissionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDefaultUserRolePermissionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDefaultUserRolePermissions(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DefaultUserRolePermissions) GetAdditionalData()(map[string]any) {
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
// GetAllowedToCreateApps gets the allowedToCreateApps property value. Indicates whether the default user role can create applications. This setting corresponds to the Users can register applications setting in the User settings menu in the Microsoft Entra admin center.
// returns a *bool when successful
func (m *DefaultUserRolePermissions) GetAllowedToCreateApps()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToCreateApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedToCreateSecurityGroups gets the allowedToCreateSecurityGroups property value. Indicates whether the default user role can create security groups. This setting corresponds to the following menus in the Microsoft Entra admin center:  The Users can create security groups in Microsoft Entra admin centers, API or PowerShell setting in the Group settings menu.  Users can create security groups setting in the User settings menu.
// returns a *bool when successful
func (m *DefaultUserRolePermissions) GetAllowedToCreateSecurityGroups()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToCreateSecurityGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedToCreateTenants gets the allowedToCreateTenants property value. Indicates whether the default user role can create tenants. This setting corresponds to the Restrict non-admin users from creating tenants setting in the User settings menu in the Microsoft Entra admin center.  When this setting is false, users assigned the Tenant Creator role can still create tenants.
// returns a *bool when successful
func (m *DefaultUserRolePermissions) GetAllowedToCreateTenants()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToCreateTenants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedToReadBitlockerKeysForOwnedDevice gets the allowedToReadBitlockerKeysForOwnedDevice property value. Indicates whether the registered owners of a device can read their own BitLocker recovery keys with default user role.
// returns a *bool when successful
func (m *DefaultUserRolePermissions) GetAllowedToReadBitlockerKeysForOwnedDevice()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToReadBitlockerKeysForOwnedDevice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowedToReadOtherUsers gets the allowedToReadOtherUsers property value. Indicates whether the default user role can read other users. DO NOT SET THIS VALUE TO false.
// returns a *bool when successful
func (m *DefaultUserRolePermissions) GetAllowedToReadOtherUsers()(*bool) {
    val, err := m.GetBackingStore().Get("allowedToReadOtherUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DefaultUserRolePermissions) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DefaultUserRolePermissions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedToCreateApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToCreateApps(val)
        }
        return nil
    }
    res["allowedToCreateSecurityGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToCreateSecurityGroups(val)
        }
        return nil
    }
    res["allowedToCreateTenants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToCreateTenants(val)
        }
        return nil
    }
    res["allowedToReadBitlockerKeysForOwnedDevice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToReadBitlockerKeysForOwnedDevice(val)
        }
        return nil
    }
    res["allowedToReadOtherUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedToReadOtherUsers(val)
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
    res["permissionGrantPoliciesAssigned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetPermissionGrantPoliciesAssigned(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DefaultUserRolePermissions) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPermissionGrantPoliciesAssigned gets the permissionGrantPoliciesAssigned property value. Indicates if user consent to apps is allowed, and if it is, which permission to grant consent and which app consent policy (permissionGrantPolicy) govern the permission for users to grant consent. Value should be in the format managePermissionGrantsForSelf.{id}, where {id} is the id of a built-in or custom app consent policy. An empty list indicates user consent to apps is disabled.
// returns a []string when successful
func (m *DefaultUserRolePermissions) GetPermissionGrantPoliciesAssigned()([]string) {
    val, err := m.GetBackingStore().Get("permissionGrantPoliciesAssigned")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DefaultUserRolePermissions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowedToCreateApps", m.GetAllowedToCreateApps())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowedToCreateSecurityGroups", m.GetAllowedToCreateSecurityGroups())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowedToCreateTenants", m.GetAllowedToCreateTenants())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowedToReadBitlockerKeysForOwnedDevice", m.GetAllowedToReadBitlockerKeysForOwnedDevice())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowedToReadOtherUsers", m.GetAllowedToReadOtherUsers())
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
    if m.GetPermissionGrantPoliciesAssigned() != nil {
        err := writer.WriteCollectionOfStringValues("permissionGrantPoliciesAssigned", m.GetPermissionGrantPoliciesAssigned())
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
func (m *DefaultUserRolePermissions) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToCreateApps sets the allowedToCreateApps property value. Indicates whether the default user role can create applications. This setting corresponds to the Users can register applications setting in the User settings menu in the Microsoft Entra admin center.
func (m *DefaultUserRolePermissions) SetAllowedToCreateApps(value *bool)() {
    err := m.GetBackingStore().Set("allowedToCreateApps", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToCreateSecurityGroups sets the allowedToCreateSecurityGroups property value. Indicates whether the default user role can create security groups. This setting corresponds to the following menus in the Microsoft Entra admin center:  The Users can create security groups in Microsoft Entra admin centers, API or PowerShell setting in the Group settings menu.  Users can create security groups setting in the User settings menu.
func (m *DefaultUserRolePermissions) SetAllowedToCreateSecurityGroups(value *bool)() {
    err := m.GetBackingStore().Set("allowedToCreateSecurityGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToCreateTenants sets the allowedToCreateTenants property value. Indicates whether the default user role can create tenants. This setting corresponds to the Restrict non-admin users from creating tenants setting in the User settings menu in the Microsoft Entra admin center.  When this setting is false, users assigned the Tenant Creator role can still create tenants.
func (m *DefaultUserRolePermissions) SetAllowedToCreateTenants(value *bool)() {
    err := m.GetBackingStore().Set("allowedToCreateTenants", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToReadBitlockerKeysForOwnedDevice sets the allowedToReadBitlockerKeysForOwnedDevice property value. Indicates whether the registered owners of a device can read their own BitLocker recovery keys with default user role.
func (m *DefaultUserRolePermissions) SetAllowedToReadBitlockerKeysForOwnedDevice(value *bool)() {
    err := m.GetBackingStore().Set("allowedToReadBitlockerKeysForOwnedDevice", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedToReadOtherUsers sets the allowedToReadOtherUsers property value. Indicates whether the default user role can read other users. DO NOT SET THIS VALUE TO false.
func (m *DefaultUserRolePermissions) SetAllowedToReadOtherUsers(value *bool)() {
    err := m.GetBackingStore().Set("allowedToReadOtherUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DefaultUserRolePermissions) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DefaultUserRolePermissions) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionGrantPoliciesAssigned sets the permissionGrantPoliciesAssigned property value. Indicates if user consent to apps is allowed, and if it is, which permission to grant consent and which app consent policy (permissionGrantPolicy) govern the permission for users to grant consent. Value should be in the format managePermissionGrantsForSelf.{id}, where {id} is the id of a built-in or custom app consent policy. An empty list indicates user consent to apps is disabled.
func (m *DefaultUserRolePermissions) SetPermissionGrantPoliciesAssigned(value []string)() {
    err := m.GetBackingStore().Set("permissionGrantPoliciesAssigned", value)
    if err != nil {
        panic(err)
    }
}
type DefaultUserRolePermissionsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedToCreateApps()(*bool)
    GetAllowedToCreateSecurityGroups()(*bool)
    GetAllowedToCreateTenants()(*bool)
    GetAllowedToReadBitlockerKeysForOwnedDevice()(*bool)
    GetAllowedToReadOtherUsers()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetPermissionGrantPoliciesAssigned()([]string)
    SetAllowedToCreateApps(value *bool)()
    SetAllowedToCreateSecurityGroups(value *bool)()
    SetAllowedToCreateTenants(value *bool)()
    SetAllowedToReadBitlockerKeysForOwnedDevice(value *bool)()
    SetAllowedToReadOtherUsers(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetPermissionGrantPoliciesAssigned(value []string)()
}
