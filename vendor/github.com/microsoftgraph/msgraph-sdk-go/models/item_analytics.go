package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemAnalytics struct {
    Entity
}
// NewItemAnalytics instantiates a new ItemAnalytics and sets the default values.
func NewItemAnalytics()(*ItemAnalytics) {
    m := &ItemAnalytics{
        Entity: *NewEntity(),
    }
    return m
}
// CreateItemAnalyticsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemAnalyticsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemAnalytics(), nil
}
// GetAllTime gets the allTime property value. The allTime property
// returns a ItemActivityStatable when successful
func (m *ItemAnalytics) GetAllTime()(ItemActivityStatable) {
    val, err := m.GetBackingStore().Get("allTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActivityStatable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemAnalytics) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActivityStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllTime(val.(ItemActivityStatable))
        }
        return nil
    }
    res["itemActivityStats"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateItemActivityStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ItemActivityStatable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ItemActivityStatable)
                }
            }
            m.SetItemActivityStats(res)
        }
        return nil
    }
    res["lastSevenDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemActivityStatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSevenDays(val.(ItemActivityStatable))
        }
        return nil
    }
    return res
}
// GetItemActivityStats gets the itemActivityStats property value. The itemActivityStats property
// returns a []ItemActivityStatable when successful
func (m *ItemAnalytics) GetItemActivityStats()([]ItemActivityStatable) {
    val, err := m.GetBackingStore().Get("itemActivityStats")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ItemActivityStatable)
    }
    return nil
}
// GetLastSevenDays gets the lastSevenDays property value. The lastSevenDays property
// returns a ItemActivityStatable when successful
func (m *ItemAnalytics) GetLastSevenDays()(ItemActivityStatable) {
    val, err := m.GetBackingStore().Get("lastSevenDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemActivityStatable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemAnalytics) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("allTime", m.GetAllTime())
        if err != nil {
            return err
        }
    }
    if m.GetItemActivityStats() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetItemActivityStats()))
        for i, v := range m.GetItemActivityStats() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("itemActivityStats", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastSevenDays", m.GetLastSevenDays())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllTime sets the allTime property value. The allTime property
func (m *ItemAnalytics) SetAllTime(value ItemActivityStatable)() {
    err := m.GetBackingStore().Set("allTime", value)
    if err != nil {
        panic(err)
    }
}
// SetItemActivityStats sets the itemActivityStats property value. The itemActivityStats property
func (m *ItemAnalytics) SetItemActivityStats(value []ItemActivityStatable)() {
    err := m.GetBackingStore().Set("itemActivityStats", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSevenDays sets the lastSevenDays property value. The lastSevenDays property
func (m *ItemAnalytics) SetLastSevenDays(value ItemActivityStatable)() {
    err := m.GetBackingStore().Set("lastSevenDays", value)
    if err != nil {
        panic(err)
    }
}
type ItemAnalyticsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllTime()(ItemActivityStatable)
    GetItemActivityStats()([]ItemActivityStatable)
    GetLastSevenDays()(ItemActivityStatable)
    SetAllTime(value ItemActivityStatable)()
    SetItemActivityStats(value []ItemActivityStatable)()
    SetLastSevenDays(value ItemActivityStatable)()
}
