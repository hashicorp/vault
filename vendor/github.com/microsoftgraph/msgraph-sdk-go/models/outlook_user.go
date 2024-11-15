package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OutlookUser struct {
    Entity
}
// NewOutlookUser instantiates a new OutlookUser and sets the default values.
func NewOutlookUser()(*OutlookUser) {
    m := &OutlookUser{
        Entity: *NewEntity(),
    }
    return m
}
// CreateOutlookUserFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOutlookUserFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOutlookUser(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OutlookUser) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["masterCategories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOutlookCategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OutlookCategoryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OutlookCategoryable)
                }
            }
            m.SetMasterCategories(res)
        }
        return nil
    }
    return res
}
// GetMasterCategories gets the masterCategories property value. A list of categories defined for the user.
// returns a []OutlookCategoryable when successful
func (m *OutlookUser) GetMasterCategories()([]OutlookCategoryable) {
    val, err := m.GetBackingStore().Get("masterCategories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OutlookCategoryable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OutlookUser) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMasterCategories() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMasterCategories()))
        for i, v := range m.GetMasterCategories() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("masterCategories", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMasterCategories sets the masterCategories property value. A list of categories defined for the user.
func (m *OutlookUser) SetMasterCategories(value []OutlookCategoryable)() {
    err := m.GetBackingStore().Set("masterCategories", value)
    if err != nil {
        panic(err)
    }
}
type OutlookUserable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMasterCategories()([]OutlookCategoryable)
    SetMasterCategories(value []OutlookCategoryable)()
}
