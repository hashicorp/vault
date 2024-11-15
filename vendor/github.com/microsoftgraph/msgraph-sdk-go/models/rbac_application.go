package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RbacApplication struct {
    Entity
}
// NewRbacApplication instantiates a new RbacApplication and sets the default values.
func NewRbacApplication()(*RbacApplication) {
    m := &RbacApplication{
        Entity: *NewEntity(),
    }
    return m
}
// CreateRbacApplicationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRbacApplicationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRbacApplication(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RbacApplication) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["resourceNamespaces"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRbacResourceNamespaceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRbacResourceNamespaceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRbacResourceNamespaceable)
                }
            }
            m.SetResourceNamespaces(res)
        }
        return nil
    }
    res["roleAssignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleAssignmentable)
                }
            }
            m.SetRoleAssignments(res)
        }
        return nil
    }
    res["roleAssignmentScheduleInstances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleAssignmentScheduleInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleAssignmentScheduleInstanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleAssignmentScheduleInstanceable)
                }
            }
            m.SetRoleAssignmentScheduleInstances(res)
        }
        return nil
    }
    res["roleAssignmentScheduleRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleAssignmentScheduleRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleAssignmentScheduleRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleAssignmentScheduleRequestable)
                }
            }
            m.SetRoleAssignmentScheduleRequests(res)
        }
        return nil
    }
    res["roleAssignmentSchedules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleAssignmentScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleAssignmentScheduleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleAssignmentScheduleable)
                }
            }
            m.SetRoleAssignmentSchedules(res)
        }
        return nil
    }
    res["roleDefinitions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleDefinitionable)
                }
            }
            m.SetRoleDefinitions(res)
        }
        return nil
    }
    res["roleEligibilityScheduleInstances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleEligibilityScheduleInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleEligibilityScheduleInstanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleEligibilityScheduleInstanceable)
                }
            }
            m.SetRoleEligibilityScheduleInstances(res)
        }
        return nil
    }
    res["roleEligibilityScheduleRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleEligibilityScheduleRequestFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleEligibilityScheduleRequestable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleEligibilityScheduleRequestable)
                }
            }
            m.SetRoleEligibilityScheduleRequests(res)
        }
        return nil
    }
    res["roleEligibilitySchedules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleEligibilityScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleEligibilityScheduleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleEligibilityScheduleable)
                }
            }
            m.SetRoleEligibilitySchedules(res)
        }
        return nil
    }
    return res
}
// GetResourceNamespaces gets the resourceNamespaces property value. The resourceNamespaces property
// returns a []UnifiedRbacResourceNamespaceable when successful
func (m *RbacApplication) GetResourceNamespaces()([]UnifiedRbacResourceNamespaceable) {
    val, err := m.GetBackingStore().Get("resourceNamespaces")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRbacResourceNamespaceable)
    }
    return nil
}
// GetRoleAssignments gets the roleAssignments property value. Resource to grant access to users or groups.
// returns a []UnifiedRoleAssignmentable when successful
func (m *RbacApplication) GetRoleAssignments()([]UnifiedRoleAssignmentable) {
    val, err := m.GetBackingStore().Get("roleAssignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleAssignmentable)
    }
    return nil
}
// GetRoleAssignmentScheduleInstances gets the roleAssignmentScheduleInstances property value. Instances for active role assignments.
// returns a []UnifiedRoleAssignmentScheduleInstanceable when successful
func (m *RbacApplication) GetRoleAssignmentScheduleInstances()([]UnifiedRoleAssignmentScheduleInstanceable) {
    val, err := m.GetBackingStore().Get("roleAssignmentScheduleInstances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleAssignmentScheduleInstanceable)
    }
    return nil
}
// GetRoleAssignmentScheduleRequests gets the roleAssignmentScheduleRequests property value. Requests for active role assignments to principals through PIM.
// returns a []UnifiedRoleAssignmentScheduleRequestable when successful
func (m *RbacApplication) GetRoleAssignmentScheduleRequests()([]UnifiedRoleAssignmentScheduleRequestable) {
    val, err := m.GetBackingStore().Get("roleAssignmentScheduleRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleAssignmentScheduleRequestable)
    }
    return nil
}
// GetRoleAssignmentSchedules gets the roleAssignmentSchedules property value. Schedules for active role assignment operations.
// returns a []UnifiedRoleAssignmentScheduleable when successful
func (m *RbacApplication) GetRoleAssignmentSchedules()([]UnifiedRoleAssignmentScheduleable) {
    val, err := m.GetBackingStore().Get("roleAssignmentSchedules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleAssignmentScheduleable)
    }
    return nil
}
// GetRoleDefinitions gets the roleDefinitions property value. Resource representing the roles allowed by RBAC providers and the permissions assigned to the roles.
// returns a []UnifiedRoleDefinitionable when successful
func (m *RbacApplication) GetRoleDefinitions()([]UnifiedRoleDefinitionable) {
    val, err := m.GetBackingStore().Get("roleDefinitions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleDefinitionable)
    }
    return nil
}
// GetRoleEligibilityScheduleInstances gets the roleEligibilityScheduleInstances property value. Instances for role eligibility requests.
// returns a []UnifiedRoleEligibilityScheduleInstanceable when successful
func (m *RbacApplication) GetRoleEligibilityScheduleInstances()([]UnifiedRoleEligibilityScheduleInstanceable) {
    val, err := m.GetBackingStore().Get("roleEligibilityScheduleInstances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleEligibilityScheduleInstanceable)
    }
    return nil
}
// GetRoleEligibilityScheduleRequests gets the roleEligibilityScheduleRequests property value. Requests for role eligibilities for principals through PIM.
// returns a []UnifiedRoleEligibilityScheduleRequestable when successful
func (m *RbacApplication) GetRoleEligibilityScheduleRequests()([]UnifiedRoleEligibilityScheduleRequestable) {
    val, err := m.GetBackingStore().Get("roleEligibilityScheduleRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleEligibilityScheduleRequestable)
    }
    return nil
}
// GetRoleEligibilitySchedules gets the roleEligibilitySchedules property value. Schedules for role eligibility operations.
// returns a []UnifiedRoleEligibilityScheduleable when successful
func (m *RbacApplication) GetRoleEligibilitySchedules()([]UnifiedRoleEligibilityScheduleable) {
    val, err := m.GetBackingStore().Get("roleEligibilitySchedules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleEligibilityScheduleable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RbacApplication) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetResourceNamespaces() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetResourceNamespaces()))
        for i, v := range m.GetResourceNamespaces() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("resourceNamespaces", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleAssignments()))
        for i, v := range m.GetRoleAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleAssignments", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignmentScheduleInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleAssignmentScheduleInstances()))
        for i, v := range m.GetRoleAssignmentScheduleInstances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleAssignmentScheduleInstances", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignmentScheduleRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleAssignmentScheduleRequests()))
        for i, v := range m.GetRoleAssignmentScheduleRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleAssignmentScheduleRequests", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleAssignmentSchedules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleAssignmentSchedules()))
        for i, v := range m.GetRoleAssignmentSchedules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleAssignmentSchedules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleDefinitions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleDefinitions()))
        for i, v := range m.GetRoleDefinitions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleDefinitions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleEligibilityScheduleInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleEligibilityScheduleInstances()))
        for i, v := range m.GetRoleEligibilityScheduleInstances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleEligibilityScheduleInstances", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleEligibilityScheduleRequests() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleEligibilityScheduleRequests()))
        for i, v := range m.GetRoleEligibilityScheduleRequests() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleEligibilityScheduleRequests", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRoleEligibilitySchedules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRoleEligibilitySchedules()))
        for i, v := range m.GetRoleEligibilitySchedules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("roleEligibilitySchedules", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetResourceNamespaces sets the resourceNamespaces property value. The resourceNamespaces property
func (m *RbacApplication) SetResourceNamespaces(value []UnifiedRbacResourceNamespaceable)() {
    err := m.GetBackingStore().Set("resourceNamespaces", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignments sets the roleAssignments property value. Resource to grant access to users or groups.
func (m *RbacApplication) SetRoleAssignments(value []UnifiedRoleAssignmentable)() {
    err := m.GetBackingStore().Set("roleAssignments", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignmentScheduleInstances sets the roleAssignmentScheduleInstances property value. Instances for active role assignments.
func (m *RbacApplication) SetRoleAssignmentScheduleInstances(value []UnifiedRoleAssignmentScheduleInstanceable)() {
    err := m.GetBackingStore().Set("roleAssignmentScheduleInstances", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignmentScheduleRequests sets the roleAssignmentScheduleRequests property value. Requests for active role assignments to principals through PIM.
func (m *RbacApplication) SetRoleAssignmentScheduleRequests(value []UnifiedRoleAssignmentScheduleRequestable)() {
    err := m.GetBackingStore().Set("roleAssignmentScheduleRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignmentSchedules sets the roleAssignmentSchedules property value. Schedules for active role assignment operations.
func (m *RbacApplication) SetRoleAssignmentSchedules(value []UnifiedRoleAssignmentScheduleable)() {
    err := m.GetBackingStore().Set("roleAssignmentSchedules", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleDefinitions sets the roleDefinitions property value. Resource representing the roles allowed by RBAC providers and the permissions assigned to the roles.
func (m *RbacApplication) SetRoleDefinitions(value []UnifiedRoleDefinitionable)() {
    err := m.GetBackingStore().Set("roleDefinitions", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleEligibilityScheduleInstances sets the roleEligibilityScheduleInstances property value. Instances for role eligibility requests.
func (m *RbacApplication) SetRoleEligibilityScheduleInstances(value []UnifiedRoleEligibilityScheduleInstanceable)() {
    err := m.GetBackingStore().Set("roleEligibilityScheduleInstances", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleEligibilityScheduleRequests sets the roleEligibilityScheduleRequests property value. Requests for role eligibilities for principals through PIM.
func (m *RbacApplication) SetRoleEligibilityScheduleRequests(value []UnifiedRoleEligibilityScheduleRequestable)() {
    err := m.GetBackingStore().Set("roleEligibilityScheduleRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleEligibilitySchedules sets the roleEligibilitySchedules property value. Schedules for role eligibility operations.
func (m *RbacApplication) SetRoleEligibilitySchedules(value []UnifiedRoleEligibilityScheduleable)() {
    err := m.GetBackingStore().Set("roleEligibilitySchedules", value)
    if err != nil {
        panic(err)
    }
}
type RbacApplicationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetResourceNamespaces()([]UnifiedRbacResourceNamespaceable)
    GetRoleAssignments()([]UnifiedRoleAssignmentable)
    GetRoleAssignmentScheduleInstances()([]UnifiedRoleAssignmentScheduleInstanceable)
    GetRoleAssignmentScheduleRequests()([]UnifiedRoleAssignmentScheduleRequestable)
    GetRoleAssignmentSchedules()([]UnifiedRoleAssignmentScheduleable)
    GetRoleDefinitions()([]UnifiedRoleDefinitionable)
    GetRoleEligibilityScheduleInstances()([]UnifiedRoleEligibilityScheduleInstanceable)
    GetRoleEligibilityScheduleRequests()([]UnifiedRoleEligibilityScheduleRequestable)
    GetRoleEligibilitySchedules()([]UnifiedRoleEligibilityScheduleable)
    SetResourceNamespaces(value []UnifiedRbacResourceNamespaceable)()
    SetRoleAssignments(value []UnifiedRoleAssignmentable)()
    SetRoleAssignmentScheduleInstances(value []UnifiedRoleAssignmentScheduleInstanceable)()
    SetRoleAssignmentScheduleRequests(value []UnifiedRoleAssignmentScheduleRequestable)()
    SetRoleAssignmentSchedules(value []UnifiedRoleAssignmentScheduleable)()
    SetRoleDefinitions(value []UnifiedRoleDefinitionable)()
    SetRoleEligibilityScheduleInstances(value []UnifiedRoleEligibilityScheduleInstanceable)()
    SetRoleEligibilityScheduleRequests(value []UnifiedRoleEligibilityScheduleRequestable)()
    SetRoleEligibilitySchedules(value []UnifiedRoleEligibilityScheduleable)()
}
