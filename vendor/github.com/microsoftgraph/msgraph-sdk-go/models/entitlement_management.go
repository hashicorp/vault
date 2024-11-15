package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EntitlementManagement struct {
    Entity
}
// NewEntitlementManagement instantiates a new EntitlementManagement and sets the default values.
func NewEntitlementManagement()(*EntitlementManagement) {
    m := &EntitlementManagement{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEntitlementManagementFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEntitlementManagementFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEntitlementManagement(), nil
}
// GetAccessPackageAssignmentApprovals gets the accessPackageAssignmentApprovals property value. Approval stages for decisions associated with access package assignment requests.
// returns a []Approvalable when successful
func (m *EntitlementManagement) GetAccessPackageAssignmentApprovals()([]Approvalable) {
    val, err := m.GetBackingStore().Get("accessPackageAssignmentApprovals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Approvalable)
    }
    return nil
}
// GetAccessPackages gets the accessPackages property value. Access packages define the collection of resource roles and the policies for which subjects can request or be assigned access to those resources.
// returns a []AccessPackageable when successful
func (m *EntitlementManagement) GetAccessPackages()([]AccessPackageable) {
    val, err := m.GetBackingStore().Get("accessPackages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageable)
    }
    return nil
}
// GetAssignmentPolicies gets the assignmentPolicies property value. Access package assignment policies govern which subjects can request or be assigned an access package via an access package assignment.
// returns a []AccessPackageAssignmentPolicyable when successful
func (m *EntitlementManagement) GetAssignmentPolicies()([]AccessPackageAssignmentPolicyable) {
    val, err := m.GetBackingStore().Get("assignmentPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageAssignmentPolicyable)
    }
    return nil
}
// GetAssignmentRequests gets the assignmentRequests property value. Access package assignment requests created by or on behalf of a subject.
// returns a []AccessPackageAssignmentRequestable when successful
func (m *EntitlementManagement) GetAssignmentRequests()([]AccessPackageAssignmentRequestable) {
    val, err := m.GetBackingStore().Get("assignmentRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageAssignmentRequestable)
    }
    return nil
}
// GetAssignments gets the assignments property value. The assignment of an access package to a subject for a period of time.
// returns a []AccessPackageAssignmentable when successful
func (m *EntitlementManagement) GetAssignments()([]AccessPackageAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageAssignmentable)
    }
    return nil
}
// GetCatalogs gets the catalogs property value. A container for access packages.
// returns a []AccessPackageCatalogable when successful
func (m *EntitlementManagement) GetCatalogs()([]AccessPackageCatalogable) {
    val, err := m.GetBackingStore().Get("catalogs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageCatalogable)
    }
    return nil
}
// GetConnectedOrganizations gets the connectedOrganizations property value. References to a directory or domain of another organization whose users can request access.
// returns a []ConnectedOrganizationable when successful
func (m *EntitlementManagement) GetConnectedOrganizations()([]ConnectedOrganizationable) {
    val, err := m.GetBackingStore().Get("connectedOrganizations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConnectedOrganizationable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EntitlementManagement) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessPackageAssignmentApprovals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateApprovalFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Approvalable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Approvalable)
                }
            }
            m.SetAccessPackageAssignmentApprovals(res)
        }
        return nil
    }
    res["accessPackages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAccessPackages(res)
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
    res["assignmentRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageAssignmentRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageAssignmentRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageAssignmentRequestable)
                }
            }
            m.SetAssignmentRequests(res)
        }
        return nil
    }
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageAssignmentable)
                }
            }
            m.SetAssignments(res)
        }
        return nil
    }
    res["catalogs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageCatalogFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageCatalogable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageCatalogable)
                }
            }
            m.SetCatalogs(res)
        }
        return nil
    }
    res["connectedOrganizations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateConnectedOrganizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConnectedOrganizationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ConnectedOrganizationable)
                }
            }
            m.SetConnectedOrganizations(res)
        }
        return nil
    }
    res["resourceEnvironments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageResourceEnvironmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageResourceEnvironmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageResourceEnvironmentable)
                }
            }
            m.SetResourceEnvironments(res)
        }
        return nil
    }
    res["resourceRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageResourceRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageResourceRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageResourceRequestable)
                }
            }
            m.SetResourceRequests(res)
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
    res["resources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageResourceable)
                }
            }
            m.SetResources(res)
        }
        return nil
    }
    res["settings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntitlementManagementSettingsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettings(val.(EntitlementManagementSettingsable))
        }
        return nil
    }
    return res
}
// GetResourceEnvironments gets the resourceEnvironments property value. A reference to the geolocation environments in which a resource is located.
// returns a []AccessPackageResourceEnvironmentable when successful
func (m *EntitlementManagement) GetResourceEnvironments()([]AccessPackageResourceEnvironmentable) {
    val, err := m.GetBackingStore().Get("resourceEnvironments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageResourceEnvironmentable)
    }
    return nil
}
// GetResourceRequests gets the resourceRequests property value. Represents a request to add or remove a resource to or from a catalog respectively.
// returns a []AccessPackageResourceRequestable when successful
func (m *EntitlementManagement) GetResourceRequests()([]AccessPackageResourceRequestable) {
    val, err := m.GetBackingStore().Get("resourceRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageResourceRequestable)
    }
    return nil
}
// GetResourceRoleScopes gets the resourceRoleScopes property value. The resourceRoleScopes property
// returns a []AccessPackageResourceRoleScopeable when successful
func (m *EntitlementManagement) GetResourceRoleScopes()([]AccessPackageResourceRoleScopeable) {
    val, err := m.GetBackingStore().Get("resourceRoleScopes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageResourceRoleScopeable)
    }
    return nil
}
// GetResources gets the resources property value. The resources associated with the catalogs.
// returns a []AccessPackageResourceable when successful
func (m *EntitlementManagement) GetResources()([]AccessPackageResourceable) {
    val, err := m.GetBackingStore().Get("resources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageResourceable)
    }
    return nil
}
// GetSettings gets the settings property value. The settings that control the behavior of Microsoft Entra entitlement management.
// returns a EntitlementManagementSettingsable when successful
func (m *EntitlementManagement) GetSettings()(EntitlementManagementSettingsable) {
    val, err := m.GetBackingStore().Get("settings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementSettingsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EntitlementManagement) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAccessPackageAssignmentApprovals() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAccessPackageAssignmentApprovals()))
        for i, v := range m.GetAccessPackageAssignmentApprovals() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("accessPackageAssignmentApprovals", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAccessPackages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAccessPackages()))
        for i, v := range m.GetAccessPackages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("accessPackages", cast)
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
    if m.GetAssignmentRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignmentRequests()))
        for i, v := range m.GetAssignmentRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignmentRequests", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetCatalogs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCatalogs()))
        for i, v := range m.GetCatalogs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("catalogs", cast)
        if err != nil {
            return err
        }
    }
    if m.GetConnectedOrganizations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetConnectedOrganizations()))
        for i, v := range m.GetConnectedOrganizations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("connectedOrganizations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResourceEnvironments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceEnvironments()))
        for i, v := range m.GetResourceEnvironments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resourceEnvironments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetResourceRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceRequests()))
        for i, v := range m.GetResourceRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resourceRequests", cast)
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
    if m.GetResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResources()))
        for i, v := range m.GetResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resources", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("settings", m.GetSettings())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackageAssignmentApprovals sets the accessPackageAssignmentApprovals property value. Approval stages for decisions associated with access package assignment requests.
func (m *EntitlementManagement) SetAccessPackageAssignmentApprovals(value []Approvalable)() {
    err := m.GetBackingStore().Set("accessPackageAssignmentApprovals", value)
    if err != nil {
        panic(err)
    }
}
// SetAccessPackages sets the accessPackages property value. Access packages define the collection of resource roles and the policies for which subjects can request or be assigned access to those resources.
func (m *EntitlementManagement) SetAccessPackages(value []AccessPackageable)() {
    err := m.GetBackingStore().Set("accessPackages", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentPolicies sets the assignmentPolicies property value. Access package assignment policies govern which subjects can request or be assigned an access package via an access package assignment.
func (m *EntitlementManagement) SetAssignmentPolicies(value []AccessPackageAssignmentPolicyable)() {
    err := m.GetBackingStore().Set("assignmentPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentRequests sets the assignmentRequests property value. Access package assignment requests created by or on behalf of a subject.
func (m *EntitlementManagement) SetAssignmentRequests(value []AccessPackageAssignmentRequestable)() {
    err := m.GetBackingStore().Set("assignmentRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignments sets the assignments property value. The assignment of an access package to a subject for a period of time.
func (m *EntitlementManagement) SetAssignments(value []AccessPackageAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetCatalogs sets the catalogs property value. A container for access packages.
func (m *EntitlementManagement) SetCatalogs(value []AccessPackageCatalogable)() {
    err := m.GetBackingStore().Set("catalogs", value)
    if err != nil {
        panic(err)
    }
}
// SetConnectedOrganizations sets the connectedOrganizations property value. References to a directory or domain of another organization whose users can request access.
func (m *EntitlementManagement) SetConnectedOrganizations(value []ConnectedOrganizationable)() {
    err := m.GetBackingStore().Set("connectedOrganizations", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceEnvironments sets the resourceEnvironments property value. A reference to the geolocation environments in which a resource is located.
func (m *EntitlementManagement) SetResourceEnvironments(value []AccessPackageResourceEnvironmentable)() {
    err := m.GetBackingStore().Set("resourceEnvironments", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceRequests sets the resourceRequests property value. Represents a request to add or remove a resource to or from a catalog respectively.
func (m *EntitlementManagement) SetResourceRequests(value []AccessPackageResourceRequestable)() {
    err := m.GetBackingStore().Set("resourceRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceRoleScopes sets the resourceRoleScopes property value. The resourceRoleScopes property
func (m *EntitlementManagement) SetResourceRoleScopes(value []AccessPackageResourceRoleScopeable)() {
    err := m.GetBackingStore().Set("resourceRoleScopes", value)
    if err != nil {
        panic(err)
    }
}
// SetResources sets the resources property value. The resources associated with the catalogs.
func (m *EntitlementManagement) SetResources(value []AccessPackageResourceable)() {
    err := m.GetBackingStore().Set("resources", value)
    if err != nil {
        panic(err)
    }
}
// SetSettings sets the settings property value. The settings that control the behavior of Microsoft Entra entitlement management.
func (m *EntitlementManagement) SetSettings(value EntitlementManagementSettingsable)() {
    err := m.GetBackingStore().Set("settings", value)
    if err != nil {
        panic(err)
    }
}
type EntitlementManagementable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackageAssignmentApprovals()([]Approvalable)
    GetAccessPackages()([]AccessPackageable)
    GetAssignmentPolicies()([]AccessPackageAssignmentPolicyable)
    GetAssignmentRequests()([]AccessPackageAssignmentRequestable)
    GetAssignments()([]AccessPackageAssignmentable)
    GetCatalogs()([]AccessPackageCatalogable)
    GetConnectedOrganizations()([]ConnectedOrganizationable)
    GetResourceEnvironments()([]AccessPackageResourceEnvironmentable)
    GetResourceRequests()([]AccessPackageResourceRequestable)
    GetResourceRoleScopes()([]AccessPackageResourceRoleScopeable)
    GetResources()([]AccessPackageResourceable)
    GetSettings()(EntitlementManagementSettingsable)
    SetAccessPackageAssignmentApprovals(value []Approvalable)()
    SetAccessPackages(value []AccessPackageable)()
    SetAssignmentPolicies(value []AccessPackageAssignmentPolicyable)()
    SetAssignmentRequests(value []AccessPackageAssignmentRequestable)()
    SetAssignments(value []AccessPackageAssignmentable)()
    SetCatalogs(value []AccessPackageCatalogable)()
    SetConnectedOrganizations(value []ConnectedOrganizationable)()
    SetResourceEnvironments(value []AccessPackageResourceEnvironmentable)()
    SetResourceRequests(value []AccessPackageResourceRequestable)()
    SetResourceRoleScopes(value []AccessPackageResourceRoleScopeable)()
    SetResources(value []AccessPackageResourceable)()
    SetSettings(value EntitlementManagementSettingsable)()
}
