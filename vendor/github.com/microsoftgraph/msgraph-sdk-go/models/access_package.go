package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackage struct {
    Entity
}
// NewAccessPackage instantiates a new AccessPackage and sets the default values.
func NewAccessPackage()(*AccessPackage) {
    m := &AccessPackage{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackage(), nil
}
// GetAccessPackagesIncompatibleWith gets the accessPackagesIncompatibleWith property value. The access packages that are incompatible with this package. Read-only.
// returns a []AccessPackageable when successful
func (m *AccessPackage) GetAccessPackagesIncompatibleWith()([]AccessPackageable) {
    val, err := m.GetBackingStore().Get("accessPackagesIncompatibleWith")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageable)
    }
    return nil
}
// GetAssignmentPolicies gets the assignmentPolicies property value. Read-only. Nullable. Supports $expand.
// returns a []AccessPackageAssignmentPolicyable when successful
func (m *AccessPackage) GetAssignmentPolicies()([]AccessPackageAssignmentPolicyable) {
    val, err := m.GetBackingStore().Get("assignmentPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageAssignmentPolicyable)
    }
    return nil
}
// GetCatalog gets the catalog property value. Required when creating the access package. Read-only. Nullable.
// returns a AccessPackageCatalogable when successful
func (m *AccessPackage) GetCatalog()(AccessPackageCatalogable) {
    val, err := m.GetBackingStore().Get("catalog")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageCatalogable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *AccessPackage) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. The description of the access package.
// returns a *string when successful
func (m *AccessPackage) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Required. The display name of the access package. Supports $filter (eq, contains).
// returns a *string when successful
func (m *AccessPackage) GetDisplayName()(*string) {
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
func (m *AccessPackage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessPackagesIncompatibleWith"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageable)
                }
            }
            m.SetAccessPackagesIncompatibleWith(res)
        }
        return nil
    }
    res["assignmentPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageAssignmentPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageAssignmentPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageAssignmentPolicyable)
                }
            }
            m.SetAssignmentPolicies(res)
        }
        return nil
    }
    res["catalog"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageCatalogFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCatalog(val.(AccessPackageCatalogable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
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
    res["incompatibleAccessPackages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageable)
                }
            }
            m.SetIncompatibleAccessPackages(res)
        }
        return nil
    }
    res["incompatibleGroups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Groupable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Groupable)
                }
            }
            m.SetIncompatibleGroups(res)
        }
        return nil
    }
    res["isHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsHidden(val)
        }
        return nil
    }
    res["modifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModifiedDateTime(val)
        }
        return nil
    }
    res["resourceRoleScopes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageResourceRoleScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageResourceRoleScopeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageResourceRoleScopeable)
                }
            }
            m.SetResourceRoleScopes(res)
        }
        return nil
    }
    return res
}
// GetIncompatibleAccessPackages gets the incompatibleAccessPackages property value. The access packages whose assigned users are ineligible to be assigned this access package.
// returns a []AccessPackageable when successful
func (m *AccessPackage) GetIncompatibleAccessPackages()([]AccessPackageable) {
    val, err := m.GetBackingStore().Get("incompatibleAccessPackages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageable)
    }
    return nil
}
// GetIncompatibleGroups gets the incompatibleGroups property value. The groups whose members are ineligible to be assigned this access package.
// returns a []Groupable when successful
func (m *AccessPackage) GetIncompatibleGroups()([]Groupable) {
    val, err := m.GetBackingStore().Get("incompatibleGroups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Groupable)
    }
    return nil
}
// GetIsHidden gets the isHidden property value. Indicates whether the access package is hidden from the requestor.
// returns a *bool when successful
func (m *AccessPackage) GetIsHidden()(*bool) {
    val, err := m.GetBackingStore().Get("isHidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetModifiedDateTime gets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *AccessPackage) GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("modifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetResourceRoleScopes gets the resourceRoleScopes property value. The resource roles and scopes in this access package.
// returns a []AccessPackageResourceRoleScopeable when successful
func (m *AccessPackage) GetResourceRoleScopes()([]AccessPackageResourceRoleScopeable) {
    val, err := m.GetBackingStore().Get("resourceRoleScopes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageResourceRoleScopeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAccessPackagesIncompatibleWith() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAccessPackagesIncompatibleWith()))
        for i, v := range m.GetAccessPackagesIncompatibleWith() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("accessPackagesIncompatibleWith", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignmentPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignmentPolicies()))
        for i, v := range m.GetAssignmentPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignmentPolicies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("catalog", m.GetCatalog())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetIncompatibleAccessPackages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIncompatibleAccessPackages()))
        for i, v := range m.GetIncompatibleAccessPackages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("incompatibleAccessPackages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetIncompatibleGroups() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetIncompatibleGroups()))
        for i, v := range m.GetIncompatibleGroups() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("incompatibleGroups", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isHidden", m.GetIsHidden())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("modifiedDateTime", m.GetModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetResourceRoleScopes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceRoleScopes()))
        for i, v := range m.GetResourceRoleScopes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resourceRoleScopes", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackagesIncompatibleWith sets the accessPackagesIncompatibleWith property value. The access packages that are incompatible with this package. Read-only.
func (m *AccessPackage) SetAccessPackagesIncompatibleWith(value []AccessPackageable)() {
    err := m.GetBackingStore().Set("accessPackagesIncompatibleWith", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentPolicies sets the assignmentPolicies property value. Read-only. Nullable. Supports $expand.
func (m *AccessPackage) SetAssignmentPolicies(value []AccessPackageAssignmentPolicyable)() {
    err := m.GetBackingStore().Set("assignmentPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetCatalog sets the catalog property value. Required when creating the access package. Read-only. Nullable.
func (m *AccessPackage) SetCatalog(value AccessPackageCatalogable)() {
    err := m.GetBackingStore().Set("catalog", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AccessPackage) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The description of the access package.
func (m *AccessPackage) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Required. The display name of the access package. Supports $filter (eq, contains).
func (m *AccessPackage) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIncompatibleAccessPackages sets the incompatibleAccessPackages property value. The access packages whose assigned users are ineligible to be assigned this access package.
func (m *AccessPackage) SetIncompatibleAccessPackages(value []AccessPackageable)() {
    err := m.GetBackingStore().Set("incompatibleAccessPackages", value)
    if err != nil {
        panic(err)
    }
}
// SetIncompatibleGroups sets the incompatibleGroups property value. The groups whose members are ineligible to be assigned this access package.
func (m *AccessPackage) SetIncompatibleGroups(value []Groupable)() {
    err := m.GetBackingStore().Set("incompatibleGroups", value)
    if err != nil {
        panic(err)
    }
}
// SetIsHidden sets the isHidden property value. Indicates whether the access package is hidden from the requestor.
func (m *AccessPackage) SetIsHidden(value *bool)() {
    err := m.GetBackingStore().Set("isHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetModifiedDateTime sets the modifiedDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AccessPackage) SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("modifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceRoleScopes sets the resourceRoleScopes property value. The resource roles and scopes in this access package.
func (m *AccessPackage) SetResourceRoleScopes(value []AccessPackageResourceRoleScopeable)() {
    err := m.GetBackingStore().Set("resourceRoleScopes", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackagesIncompatibleWith()([]AccessPackageable)
    GetAssignmentPolicies()([]AccessPackageAssignmentPolicyable)
    GetCatalog()(AccessPackageCatalogable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetIncompatibleAccessPackages()([]AccessPackageable)
    GetIncompatibleGroups()([]Groupable)
    GetIsHidden()(*bool)
    GetModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetResourceRoleScopes()([]AccessPackageResourceRoleScopeable)
    SetAccessPackagesIncompatibleWith(value []AccessPackageable)()
    SetAssignmentPolicies(value []AccessPackageAssignmentPolicyable)()
    SetCatalog(value AccessPackageCatalogable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetIncompatibleAccessPackages(value []AccessPackageable)()
    SetIncompatibleGroups(value []Groupable)()
    SetIsHidden(value *bool)()
    SetModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetResourceRoleScopes(value []AccessPackageResourceRoleScopeable)()
}
