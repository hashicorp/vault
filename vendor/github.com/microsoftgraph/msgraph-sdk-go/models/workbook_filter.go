package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookFilter struct {
    Entity
}
// NewWorkbookFilter instantiates a new WorkbookFilter and sets the default values.
func NewWorkbookFilter()(*WorkbookFilter) {
    m := &WorkbookFilter{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFilterFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookFilterFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookFilter(), nil
}
// GetCriteria gets the criteria property value. The currently applied filter on the given column. Read-only.
// returns a WorkbookFilterCriteriaable when successful
func (m *WorkbookFilter) GetCriteria()(WorkbookFilterCriteriaable) {
    val, err := m.GetBackingStore().Get("criteria")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookFilterCriteriaable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookFilter) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["criteria"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookFilterCriteriaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCriteria(val.(WorkbookFilterCriteriaable))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WorkbookFilter) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("criteria", m.GetCriteria())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCriteria sets the criteria property value. The currently applied filter on the given column. Read-only.
func (m *WorkbookFilter) SetCriteria(value WorkbookFilterCriteriaable)() {
    err := m.GetBackingStore().Set("criteria", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookFilterable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCriteria()(WorkbookFilterCriteriaable)
    SetCriteria(value WorkbookFilterCriteriaable)()
}
