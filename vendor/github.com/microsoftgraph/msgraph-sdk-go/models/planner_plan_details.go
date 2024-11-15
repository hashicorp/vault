package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PlannerPlanDetails struct {
    Entity
}
// NewPlannerPlanDetails instantiates a new PlannerPlanDetails and sets the default values.
func NewPlannerPlanDetails()(*PlannerPlanDetails) {
    m := &PlannerPlanDetails{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerPlanDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlannerPlanDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerPlanDetails(), nil
}
// GetCategoryDescriptions gets the categoryDescriptions property value. An object that specifies the descriptions of the 25 categories that can be associated with tasks in the plan.
// returns a PlannerCategoryDescriptionsable when successful
func (m *PlannerPlanDetails) GetCategoryDescriptions()(PlannerCategoryDescriptionsable) {
    val, err := m.GetBackingStore().Get("categoryDescriptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerCategoryDescriptionsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PlannerPlanDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["categoryDescriptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerCategoryDescriptionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategoryDescriptions(val.(PlannerCategoryDescriptionsable))
        }
        return nil
    }
    res["sharedWith"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePlannerUserIdsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharedWith(val.(PlannerUserIdsable))
        }
        return nil
    }
    return res
}
// GetSharedWith gets the sharedWith property value. Set of user IDs that this plan is shared with. If you're using Microsoft 365 groups, use the Groups API to manage group membership to share the group's plan. You can also add existing members of the group to this collection, although it isn't required for them to access the plan owned by the group.
// returns a PlannerUserIdsable when successful
func (m *PlannerPlanDetails) GetSharedWith()(PlannerUserIdsable) {
    val, err := m.GetBackingStore().Get("sharedWith")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PlannerUserIdsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PlannerPlanDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("categoryDescriptions", m.GetCategoryDescriptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharedWith", m.GetSharedWith())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCategoryDescriptions sets the categoryDescriptions property value. An object that specifies the descriptions of the 25 categories that can be associated with tasks in the plan.
func (m *PlannerPlanDetails) SetCategoryDescriptions(value PlannerCategoryDescriptionsable)() {
    err := m.GetBackingStore().Set("categoryDescriptions", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedWith sets the sharedWith property value. Set of user IDs that this plan is shared with. If you're using Microsoft 365 groups, use the Groups API to manage group membership to share the group's plan. You can also add existing members of the group to this collection, although it isn't required for them to access the plan owned by the group.
func (m *PlannerPlanDetails) SetSharedWith(value PlannerUserIdsable)() {
    err := m.GetBackingStore().Set("sharedWith", value)
    if err != nil {
        panic(err)
    }
}
type PlannerPlanDetailsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCategoryDescriptions()(PlannerCategoryDescriptionsable)
    GetSharedWith()(PlannerUserIdsable)
    SetCategoryDescriptions(value PlannerCategoryDescriptionsable)()
    SetSharedWith(value PlannerUserIdsable)()
}
