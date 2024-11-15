package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventsRoot struct {
    Entity
}
// NewVirtualEventsRoot instantiates a new VirtualEventsRoot and sets the default values.
func NewVirtualEventsRoot()(*VirtualEventsRoot) {
    m := &VirtualEventsRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateVirtualEventsRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventsRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventsRoot(), nil
}
// GetEvents gets the events property value. The events property
// returns a []VirtualEventable when successful
func (m *VirtualEventsRoot) GetEvents()([]VirtualEventable) {
    val, err := m.GetBackingStore().Get("events")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]VirtualEventable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventsRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["events"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateVirtualEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]VirtualEventable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(VirtualEventable)
                }
            }
            m.SetEvents(res)
        }
        return nil
    }
    res["townhalls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateVirtualEventTownhallFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]VirtualEventTownhallable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(VirtualEventTownhallable)
                }
            }
            m.SetTownhalls(res)
        }
        return nil
    }
    res["webinars"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateVirtualEventWebinarFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]VirtualEventWebinarable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(VirtualEventWebinarable)
                }
            }
            m.SetWebinars(res)
        }
        return nil
    }
    return res
}
// GetTownhalls gets the townhalls property value. A collection of town halls. Nullable.
// returns a []VirtualEventTownhallable when successful
func (m *VirtualEventsRoot) GetTownhalls()([]VirtualEventTownhallable) {
    val, err := m.GetBackingStore().Get("townhalls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]VirtualEventTownhallable)
    }
    return nil
}
// GetWebinars gets the webinars property value. A collection of webinars. Nullable.
// returns a []VirtualEventWebinarable when successful
func (m *VirtualEventsRoot) GetWebinars()([]VirtualEventWebinarable) {
    val, err := m.GetBackingStore().Get("webinars")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]VirtualEventWebinarable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventsRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetEvents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEvents()))
        for i, v := range m.GetEvents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("events", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTownhalls() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTownhalls()))
        for i, v := range m.GetTownhalls() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("townhalls", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWebinars() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWebinars()))
        for i, v := range m.GetWebinars() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("webinars", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEvents sets the events property value. The events property
func (m *VirtualEventsRoot) SetEvents(value []VirtualEventable)() {
    err := m.GetBackingStore().Set("events", value)
    if err != nil {
        panic(err)
    }
}
// SetTownhalls sets the townhalls property value. A collection of town halls. Nullable.
func (m *VirtualEventsRoot) SetTownhalls(value []VirtualEventTownhallable)() {
    err := m.GetBackingStore().Set("townhalls", value)
    if err != nil {
        panic(err)
    }
}
// SetWebinars sets the webinars property value. A collection of webinars. Nullable.
func (m *VirtualEventsRoot) SetWebinars(value []VirtualEventWebinarable)() {
    err := m.GetBackingStore().Set("webinars", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventsRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEvents()([]VirtualEventable)
    GetTownhalls()([]VirtualEventTownhallable)
    GetWebinars()([]VirtualEventWebinarable)
    SetEvents(value []VirtualEventable)()
    SetTownhalls(value []VirtualEventTownhallable)()
    SetWebinars(value []VirtualEventWebinarable)()
}
