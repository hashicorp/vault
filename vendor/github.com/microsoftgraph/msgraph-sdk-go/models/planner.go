package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Planner struct {
    Entity
}
// NewPlanner instantiates a new Planner and sets the default values.
func NewPlanner()(*Planner) {
    m := &Planner{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePlannerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlannerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlanner(), nil
}
// GetBuckets gets the buckets property value. Read-only. Nullable. Returns a collection of the specified buckets
// returns a []PlannerBucketable when successful
func (m *Planner) GetBuckets()([]PlannerBucketable) {
    val, err := m.GetBackingStore().Get("buckets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PlannerBucketable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Planner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["buckets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePlannerBucketFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PlannerBucketable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PlannerBucketable)
                }
            }
            m.SetBuckets(res)
        }
        return nil
    }
    res["plans"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePlannerPlanFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PlannerPlanable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PlannerPlanable)
                }
            }
            m.SetPlans(res)
        }
        return nil
    }
    res["tasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePlannerTaskFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PlannerTaskable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PlannerTaskable)
                }
            }
            m.SetTasks(res)
        }
        return nil
    }
    return res
}
// GetPlans gets the plans property value. Read-only. Nullable. Returns a collection of the specified plans
// returns a []PlannerPlanable when successful
func (m *Planner) GetPlans()([]PlannerPlanable) {
    val, err := m.GetBackingStore().Get("plans")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PlannerPlanable)
    }
    return nil
}
// GetTasks gets the tasks property value. Read-only. Nullable. Returns a collection of the specified tasks
// returns a []PlannerTaskable when successful
func (m *Planner) GetTasks()([]PlannerTaskable) {
    val, err := m.GetBackingStore().Get("tasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PlannerTaskable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Planner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBuckets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBuckets()))
        for i, v := range m.GetBuckets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("buckets", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPlans() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPlans()))
        for i, v := range m.GetPlans() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("plans", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTasks() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTasks()))
        for i, v := range m.GetTasks() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tasks", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBuckets sets the buckets property value. Read-only. Nullable. Returns a collection of the specified buckets
func (m *Planner) SetBuckets(value []PlannerBucketable)() {
    err := m.GetBackingStore().Set("buckets", value)
    if err != nil {
        panic(err)
    }
}
// SetPlans sets the plans property value. Read-only. Nullable. Returns a collection of the specified plans
func (m *Planner) SetPlans(value []PlannerPlanable)() {
    err := m.GetBackingStore().Set("plans", value)
    if err != nil {
        panic(err)
    }
}
// SetTasks sets the tasks property value. Read-only. Nullable. Returns a collection of the specified tasks
func (m *Planner) SetTasks(value []PlannerTaskable)() {
    err := m.GetBackingStore().Set("tasks", value)
    if err != nil {
        panic(err)
    }
}
type Plannerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBuckets()([]PlannerBucketable)
    GetPlans()([]PlannerPlanable)
    GetTasks()([]PlannerTaskable)
    SetBuckets(value []PlannerBucketable)()
    SetPlans(value []PlannerPlanable)()
    SetTasks(value []PlannerTaskable)()
}
