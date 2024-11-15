package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserSolutionRoot struct {
    Entity
}
// NewUserSolutionRoot instantiates a new UserSolutionRoot and sets the default values.
func NewUserSolutionRoot()(*UserSolutionRoot) {
    m := &UserSolutionRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserSolutionRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserSolutionRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserSolutionRoot(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserSolutionRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["workingTimeSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkingTimeScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkingTimeSchedule(val.(WorkingTimeScheduleable))
        }
        return nil
    }
    return res
}
// GetWorkingTimeSchedule gets the workingTimeSchedule property value. The working time schedule entity associated with the solution.
// returns a WorkingTimeScheduleable when successful
func (m *UserSolutionRoot) GetWorkingTimeSchedule()(WorkingTimeScheduleable) {
    val, err := m.GetBackingStore().Get("workingTimeSchedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkingTimeScheduleable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserSolutionRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("workingTimeSchedule", m.GetWorkingTimeSchedule())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetWorkingTimeSchedule sets the workingTimeSchedule property value. The working time schedule entity associated with the solution.
func (m *UserSolutionRoot) SetWorkingTimeSchedule(value WorkingTimeScheduleable)() {
    err := m.GetBackingStore().Set("workingTimeSchedule", value)
    if err != nil {
        panic(err)
    }
}
type UserSolutionRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetWorkingTimeSchedule()(WorkingTimeScheduleable)
    SetWorkingTimeSchedule(value WorkingTimeScheduleable)()
}
