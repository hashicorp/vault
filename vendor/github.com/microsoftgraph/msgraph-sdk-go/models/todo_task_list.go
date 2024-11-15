package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TodoTaskList struct {
    Entity
}
// NewTodoTaskList instantiates a new TodoTaskList and sets the default values.
func NewTodoTaskList()(*TodoTaskList) {
    m := &TodoTaskList{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTodoTaskListFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTodoTaskListFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTodoTaskList(), nil
}
// GetDisplayName gets the displayName property value. The name of the task list.
// returns a *string when successful
func (m *TodoTaskList) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the task list. Nullable.
// returns a []Extensionable when successful
func (m *TodoTaskList) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TodoTaskList) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["extensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Extensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Extensionable)
                }
            }
            m.SetExtensions(res)
        }
        return nil
    }
    res["isOwner"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOwner(val)
        }
        return nil
    }
    res["isShared"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsShared(val)
        }
        return nil
    }
    res["tasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTodoTaskFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TodoTaskable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TodoTaskable)
                }
            }
            m.SetTasks(res)
        }
        return nil
    }
    res["wellknownListName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWellknownListName)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWellknownListName(val.(*WellknownListName))
        }
        return nil
    }
    return res
}
// GetIsOwner gets the isOwner property value. True if the user is owner of the given task list.
// returns a *bool when successful
func (m *TodoTaskList) GetIsOwner()(*bool) {
    val, err := m.GetBackingStore().Get("isOwner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsShared gets the isShared property value. True if the task list is shared with other users
// returns a *bool when successful
func (m *TodoTaskList) GetIsShared()(*bool) {
    val, err := m.GetBackingStore().Get("isShared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTasks gets the tasks property value. The tasks in this task list. Read-only. Nullable.
// returns a []TodoTaskable when successful
func (m *TodoTaskList) GetTasks()([]TodoTaskable) {
    val, err := m.GetBackingStore().Get("tasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TodoTaskable)
    }
    return nil
}
// GetWellknownListName gets the wellknownListName property value. The wellknownListName property
// returns a *WellknownListName when successful
func (m *TodoTaskList) GetWellknownListName()(*WellknownListName) {
    val, err := m.GetBackingStore().Get("wellknownListName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WellknownListName)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TodoTaskList) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensions()))
        for i, v := range m.GetExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("extensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOwner", m.GetIsOwner())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isShared", m.GetIsShared())
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
    if m.GetWellknownListName() != nil {
        cast := (*m.GetWellknownListName()).String()
        err = writer.WriteStringValue("wellknownListName", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The name of the task list.
func (m *TodoTaskList) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the task list. Nullable.
func (m *TodoTaskList) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOwner sets the isOwner property value. True if the user is owner of the given task list.
func (m *TodoTaskList) SetIsOwner(value *bool)() {
    err := m.GetBackingStore().Set("isOwner", value)
    if err != nil {
        panic(err)
    }
}
// SetIsShared sets the isShared property value. True if the task list is shared with other users
func (m *TodoTaskList) SetIsShared(value *bool)() {
    err := m.GetBackingStore().Set("isShared", value)
    if err != nil {
        panic(err)
    }
}
// SetTasks sets the tasks property value. The tasks in this task list. Read-only. Nullable.
func (m *TodoTaskList) SetTasks(value []TodoTaskable)() {
    err := m.GetBackingStore().Set("tasks", value)
    if err != nil {
        panic(err)
    }
}
// SetWellknownListName sets the wellknownListName property value. The wellknownListName property
func (m *TodoTaskList) SetWellknownListName(value *WellknownListName)() {
    err := m.GetBackingStore().Set("wellknownListName", value)
    if err != nil {
        panic(err)
    }
}
type TodoTaskListable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetExtensions()([]Extensionable)
    GetIsOwner()(*bool)
    GetIsShared()(*bool)
    GetTasks()([]TodoTaskable)
    GetWellknownListName()(*WellknownListName)
    SetDisplayName(value *string)()
    SetExtensions(value []Extensionable)()
    SetIsOwner(value *bool)()
    SetIsShared(value *bool)()
    SetTasks(value []TodoTaskable)()
    SetWellknownListName(value *WellknownListName)()
}
