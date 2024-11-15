package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemActivityStat struct {
    Entity
}
// NewItemActivityStat instantiates a new ItemActivityStat and sets the default values.
func NewItemActivityStat()(*ItemActivityStat) {
    m := &ItemActivityStat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateItemActivityStatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemActivityStatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemActivityStat(), nil
}
// GetAccess gets the access property value. Statistics about the access actions in this interval. Read-only.
// returns a ItemActionStatable when successful
func (m *ItemActivityStat) GetAccess()(ItemActionStatable) {
    val, err := m.GetBackingStore().Get("access")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActionStatable)
    }
    return nil
}
// GetActivities gets the activities property value. Exposes the itemActivities represented in this itemActivityStat resource.
// returns a []ItemActivityable when successful
func (m *ItemActivityStat) GetActivities()([]ItemActivityable) {
    val, err := m.GetBackingStore().Get("activities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ItemActivityable)
    }
    return nil
}
// GetCreate gets the create property value. Statistics about the create actions in this interval. Read-only.
// returns a ItemActionStatable when successful
func (m *ItemActivityStat) GetCreate()(ItemActionStatable) {
    val, err := m.GetBackingStore().Get("create")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActionStatable)
    }
    return nil
}
// GetDelete gets the delete property value. Statistics about the delete actions in this interval. Read-only.
// returns a ItemActionStatable when successful
func (m *ItemActivityStat) GetDelete()(ItemActionStatable) {
    val, err := m.GetBackingStore().Get("delete")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActionStatable)
    }
    return nil
}
// GetEdit gets the edit property value. Statistics about the edit actions in this interval. Read-only.
// returns a ItemActionStatable when successful
func (m *ItemActivityStat) GetEdit()(ItemActionStatable) {
    val, err := m.GetBackingStore().Get("edit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActionStatable)
    }
    return nil
}
// GetEndDateTime gets the endDateTime property value. When the interval ends. Read-only.
// returns a *Time when successful
func (m *ItemActivityStat) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemActivityStat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["access"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActionStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccess(val.(ItemActionStatable))
        }
        return nil
    }
    res["activities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateItemActivityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ItemActivityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ItemActivityable)
                }
            }
            m.SetActivities(res)
        }
        return nil
    }
    res["create"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActionStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreate(val.(ItemActionStatable))
        }
        return nil
    }
    res["delete"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActionStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDelete(val.(ItemActionStatable))
        }
        return nil
    }
    res["edit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActionStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdit(val.(ItemActionStatable))
        }
        return nil
    }
    res["endDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDateTime(val)
        }
        return nil
    }
    res["incompleteData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIncompleteDataFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncompleteData(val.(IncompleteDataable))
        }
        return nil
    }
    res["isTrending"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsTrending(val)
        }
        return nil
    }
    res["move"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActionStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMove(val.(ItemActionStatable))
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val)
        }
        return nil
    }
    return res
}
// GetIncompleteData gets the incompleteData property value. Indicates that the statistics in this interval are based on incomplete data. Read-only.
// returns a IncompleteDataable when successful
func (m *ItemActivityStat) GetIncompleteData()(IncompleteDataable) {
    val, err := m.GetBackingStore().Get("incompleteData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IncompleteDataable)
    }
    return nil
}
// GetIsTrending gets the isTrending property value. Indicates whether the item is 'trending.' Read-only.
// returns a *bool when successful
func (m *ItemActivityStat) GetIsTrending()(*bool) {
    val, err := m.GetBackingStore().Get("isTrending")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMove gets the move property value. Statistics about the move actions in this interval. Read-only.
// returns a ItemActionStatable when successful
func (m *ItemActivityStat) GetMove()(ItemActionStatable) {
    val, err := m.GetBackingStore().Get("move")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActionStatable)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. When the interval starts. Read-only.
// returns a *Time when successful
func (m *ItemActivityStat) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemActivityStat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("access", m.GetAccess())
        if err != nil {
            return err
        }
    }
    if m.GetActivities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetActivities()))
        for i, v := range m.GetActivities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("activities", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("create", m.GetCreate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("delete", m.GetDelete())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("edit", m.GetEdit())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("incompleteData", m.GetIncompleteData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isTrending", m.GetIsTrending())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("move", m.GetMove())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccess sets the access property value. Statistics about the access actions in this interval. Read-only.
func (m *ItemActivityStat) SetAccess(value ItemActionStatable)() {
    err := m.GetBackingStore().Set("access", value)
    if err != nil {
        panic(err)
    }
}
// SetActivities sets the activities property value. Exposes the itemActivities represented in this itemActivityStat resource.
func (m *ItemActivityStat) SetActivities(value []ItemActivityable)() {
    err := m.GetBackingStore().Set("activities", value)
    if err != nil {
        panic(err)
    }
}
// SetCreate sets the create property value. Statistics about the create actions in this interval. Read-only.
func (m *ItemActivityStat) SetCreate(value ItemActionStatable)() {
    err := m.GetBackingStore().Set("create", value)
    if err != nil {
        panic(err)
    }
}
// SetDelete sets the delete property value. Statistics about the delete actions in this interval. Read-only.
func (m *ItemActivityStat) SetDelete(value ItemActionStatable)() {
    err := m.GetBackingStore().Set("delete", value)
    if err != nil {
        panic(err)
    }
}
// SetEdit sets the edit property value. Statistics about the edit actions in this interval. Read-only.
func (m *ItemActivityStat) SetEdit(value ItemActionStatable)() {
    err := m.GetBackingStore().Set("edit", value)
    if err != nil {
        panic(err)
    }
}
// SetEndDateTime sets the endDateTime property value. When the interval ends. Read-only.
func (m *ItemActivityStat) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIncompleteData sets the incompleteData property value. Indicates that the statistics in this interval are based on incomplete data. Read-only.
func (m *ItemActivityStat) SetIncompleteData(value IncompleteDataable)() {
    err := m.GetBackingStore().Set("incompleteData", value)
    if err != nil {
        panic(err)
    }
}
// SetIsTrending sets the isTrending property value. Indicates whether the item is 'trending.' Read-only.
func (m *ItemActivityStat) SetIsTrending(value *bool)() {
    err := m.GetBackingStore().Set("isTrending", value)
    if err != nil {
        panic(err)
    }
}
// SetMove sets the move property value. Statistics about the move actions in this interval. Read-only.
func (m *ItemActivityStat) SetMove(value ItemActionStatable)() {
    err := m.GetBackingStore().Set("move", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. When the interval starts. Read-only.
func (m *ItemActivityStat) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
type ItemActivityStatable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccess()(ItemActionStatable)
    GetActivities()([]ItemActivityable)
    GetCreate()(ItemActionStatable)
    GetDelete()(ItemActionStatable)
    GetEdit()(ItemActionStatable)
    GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIncompleteData()(IncompleteDataable)
    GetIsTrending()(*bool)
    GetMove()(ItemActionStatable)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetAccess(value ItemActionStatable)()
    SetActivities(value []ItemActivityable)()
    SetCreate(value ItemActionStatable)()
    SetDelete(value ItemActionStatable)()
    SetEdit(value ItemActionStatable)()
    SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIncompleteData(value IncompleteDataable)()
    SetIsTrending(value *bool)()
    SetMove(value ItemActionStatable)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
