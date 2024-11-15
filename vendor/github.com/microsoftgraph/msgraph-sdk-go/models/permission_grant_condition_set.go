package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PermissionGrantConditionSet struct {
    Entity
}
// NewPermissionGrantConditionSet instantiates a new PermissionGrantConditionSet and sets the default values.
func NewPermissionGrantConditionSet()(*PermissionGrantConditionSet) {
    m := &PermissionGrantConditionSet{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePermissionGrantConditionSetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePermissionGrantConditionSetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPermissionGrantConditionSet(), nil
}
// GetClientApplicationIds gets the clientApplicationIds property value. A list of appId values for the client applications to match with, or a list with the single value all to match any client application. Default is the single value all.
// returns a []string when successful
func (m *PermissionGrantConditionSet) GetClientApplicationIds()([]string) {
    val, err := m.GetBackingStore().Get("clientApplicationIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetClientApplicationPublisherIds gets the clientApplicationPublisherIds property value. A list of Microsoft Partner Network (MPN) IDs for verified publishers of the client application, or a list with the single value all to match with client apps from any publisher. Default is the single value all.
// returns a []string when successful
func (m *PermissionGrantConditionSet) GetClientApplicationPublisherIds()([]string) {
    val, err := m.GetBackingStore().Get("clientApplicationPublisherIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetClientApplicationsFromVerifiedPublisherOnly gets the clientApplicationsFromVerifiedPublisherOnly property value. Set to true to only match on client applications with a verified publisher. Set to false to match on any client app, even if it doesn't have a verified publisher. Default is false.
// returns a *bool when successful
func (m *PermissionGrantConditionSet) GetClientApplicationsFromVerifiedPublisherOnly()(*bool) {
    val, err := m.GetBackingStore().Get("clientApplicationsFromVerifiedPublisherOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetClientApplicationTenantIds gets the clientApplicationTenantIds property value. A list of Microsoft Entra tenant IDs in which the client application is registered, or a list with the single value all to match with client apps registered in any tenant. Default is the single value all.
// returns a []string when successful
func (m *PermissionGrantConditionSet) GetClientApplicationTenantIds()([]string) {
    val, err := m.GetBackingStore().Get("clientApplicationTenantIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PermissionGrantConditionSet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["clientApplicationIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetClientApplicationIds(res)
        }
        return nil
    }
    res["clientApplicationPublisherIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetClientApplicationPublisherIds(res)
        }
        return nil
    }
    res["clientApplicationsFromVerifiedPublisherOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClientApplicationsFromVerifiedPublisherOnly(val)
        }
        return nil
    }
    res["clientApplicationTenantIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetClientApplicationTenantIds(res)
        }
        return nil
    }
    res["permissionClassification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermissionClassification(val)
        }
        return nil
    }
    res["permissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetPermissions(res)
        }
        return nil
    }
    res["permissionType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePermissionType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPermissionType(val.(*PermissionType))
        }
        return nil
    }
    res["resourceApplication"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceApplication(val)
        }
        return nil
    }
    return res
}
// GetPermissionClassification gets the permissionClassification property value. The permission classification for the permission being granted, or all to match with any permission classification (including permissions that aren't classified). Default is all.
// returns a *string when successful
func (m *PermissionGrantConditionSet) GetPermissionClassification()(*string) {
    val, err := m.GetBackingStore().Get("permissionClassification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPermissions gets the permissions property value. The list of id values for the specific permissions to match with, or a list with the single value all to match with any permission. The id of delegated permissions can be found in the oauth2PermissionScopes property of the API's servicePrincipal object. The id of application permissions can be found in the appRoles property of the API's servicePrincipal object. The id of resource-specific application permissions can be found in the resourceSpecificApplicationPermissions property of the API's servicePrincipal object. Default is the single value all.
// returns a []string when successful
func (m *PermissionGrantConditionSet) GetPermissions()([]string) {
    val, err := m.GetBackingStore().Get("permissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPermissionType gets the permissionType property value. The permission type of the permission being granted. Possible values: application for application permissions (for example app roles), or delegated for delegated permissions. The value delegatedUserConsentable indicates delegated permissions that haven't been configured by the API publisher to require admin consent—this value may be used in built-in permission grant policies, but can't be used in custom permission grant policies. Required.
// returns a *PermissionType when successful
func (m *PermissionGrantConditionSet) GetPermissionType()(*PermissionType) {
    val, err := m.GetBackingStore().Get("permissionType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PermissionType)
    }
    return nil
}
// GetResourceApplication gets the resourceApplication property value. The appId of the resource application (for example the API) for which a permission is being granted, or any to match with any resource application or API. Default is any.
// returns a *string when successful
func (m *PermissionGrantConditionSet) GetResourceApplication()(*string) {
    val, err := m.GetBackingStore().Get("resourceApplication")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PermissionGrantConditionSet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetClientApplicationIds() != nil {
        err = writer.WriteCollectionOfStringValues("clientApplicationIds", m.GetClientApplicationIds())
        if err != nil {
            return err
        }
    }
    if m.GetClientApplicationPublisherIds() != nil {
        err = writer.WriteCollectionOfStringValues("clientApplicationPublisherIds", m.GetClientApplicationPublisherIds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("clientApplicationsFromVerifiedPublisherOnly", m.GetClientApplicationsFromVerifiedPublisherOnly())
        if err != nil {
            return err
        }
    }
    if m.GetClientApplicationTenantIds() != nil {
        err = writer.WriteCollectionOfStringValues("clientApplicationTenantIds", m.GetClientApplicationTenantIds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("permissionClassification", m.GetPermissionClassification())
        if err != nil {
            return err
        }
    }
    if m.GetPermissions() != nil {
        err = writer.WriteCollectionOfStringValues("permissions", m.GetPermissions())
        if err != nil {
            return err
        }
    }
    if m.GetPermissionType() != nil {
        cast := (*m.GetPermissionType()).String()
        err = writer.WriteStringValue("permissionType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceApplication", m.GetResourceApplication())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClientApplicationIds sets the clientApplicationIds property value. A list of appId values for the client applications to match with, or a list with the single value all to match any client application. Default is the single value all.
func (m *PermissionGrantConditionSet) SetClientApplicationIds(value []string)() {
    err := m.GetBackingStore().Set("clientApplicationIds", value)
    if err != nil {
        panic(err)
    }
}
// SetClientApplicationPublisherIds sets the clientApplicationPublisherIds property value. A list of Microsoft Partner Network (MPN) IDs for verified publishers of the client application, or a list with the single value all to match with client apps from any publisher. Default is the single value all.
func (m *PermissionGrantConditionSet) SetClientApplicationPublisherIds(value []string)() {
    err := m.GetBackingStore().Set("clientApplicationPublisherIds", value)
    if err != nil {
        panic(err)
    }
}
// SetClientApplicationsFromVerifiedPublisherOnly sets the clientApplicationsFromVerifiedPublisherOnly property value. Set to true to only match on client applications with a verified publisher. Set to false to match on any client app, even if it doesn't have a verified publisher. Default is false.
func (m *PermissionGrantConditionSet) SetClientApplicationsFromVerifiedPublisherOnly(value *bool)() {
    err := m.GetBackingStore().Set("clientApplicationsFromVerifiedPublisherOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetClientApplicationTenantIds sets the clientApplicationTenantIds property value. A list of Microsoft Entra tenant IDs in which the client application is registered, or a list with the single value all to match with client apps registered in any tenant. Default is the single value all.
func (m *PermissionGrantConditionSet) SetClientApplicationTenantIds(value []string)() {
    err := m.GetBackingStore().Set("clientApplicationTenantIds", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionClassification sets the permissionClassification property value. The permission classification for the permission being granted, or all to match with any permission classification (including permissions that aren't classified). Default is all.
func (m *PermissionGrantConditionSet) SetPermissionClassification(value *string)() {
    err := m.GetBackingStore().Set("permissionClassification", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissions sets the permissions property value. The list of id values for the specific permissions to match with, or a list with the single value all to match with any permission. The id of delegated permissions can be found in the oauth2PermissionScopes property of the API's servicePrincipal object. The id of application permissions can be found in the appRoles property of the API's servicePrincipal object. The id of resource-specific application permissions can be found in the resourceSpecificApplicationPermissions property of the API's servicePrincipal object. Default is the single value all.
func (m *PermissionGrantConditionSet) SetPermissions(value []string)() {
    err := m.GetBackingStore().Set("permissions", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissionType sets the permissionType property value. The permission type of the permission being granted. Possible values: application for application permissions (for example app roles), or delegated for delegated permissions. The value delegatedUserConsentable indicates delegated permissions that haven't been configured by the API publisher to require admin consent—this value may be used in built-in permission grant policies, but can't be used in custom permission grant policies. Required.
func (m *PermissionGrantConditionSet) SetPermissionType(value *PermissionType)() {
    err := m.GetBackingStore().Set("permissionType", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceApplication sets the resourceApplication property value. The appId of the resource application (for example the API) for which a permission is being granted, or any to match with any resource application or API. Default is any.
func (m *PermissionGrantConditionSet) SetResourceApplication(value *string)() {
    err := m.GetBackingStore().Set("resourceApplication", value)
    if err != nil {
        panic(err)
    }
}
type PermissionGrantConditionSetable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClientApplicationIds()([]string)
    GetClientApplicationPublisherIds()([]string)
    GetClientApplicationsFromVerifiedPublisherOnly()(*bool)
    GetClientApplicationTenantIds()([]string)
    GetPermissionClassification()(*string)
    GetPermissions()([]string)
    GetPermissionType()(*PermissionType)
    GetResourceApplication()(*string)
    SetClientApplicationIds(value []string)()
    SetClientApplicationPublisherIds(value []string)()
    SetClientApplicationsFromVerifiedPublisherOnly(value *bool)()
    SetClientApplicationTenantIds(value []string)()
    SetPermissionClassification(value *string)()
    SetPermissions(value []string)()
    SetPermissionType(value *PermissionType)()
    SetResourceApplication(value *string)()
}
