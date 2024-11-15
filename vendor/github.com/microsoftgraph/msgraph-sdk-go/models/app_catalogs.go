package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AppCatalogs struct {
    Entity
}
// NewAppCatalogs instantiates a new AppCatalogs and sets the default values.
func NewAppCatalogs()(*AppCatalogs) {
    m := &AppCatalogs{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppCatalogsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppCatalogsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppCatalogs(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AppCatalogs) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["teamsApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamsAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamsAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamsAppable)
                }
            }
            m.SetTeamsApps(res)
        }
        return nil
    }
    return res
}
// GetTeamsApps gets the teamsApps property value. The teamsApps property
// returns a []TeamsAppable when successful
func (m *AppCatalogs) GetTeamsApps()([]TeamsAppable) {
    val, err := m.GetBackingStore().Get("teamsApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsAppable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AppCatalogs) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetTeamsApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTeamsApps()))
        for i, v := range m.GetTeamsApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("teamsApps", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetTeamsApps sets the teamsApps property value. The teamsApps property
func (m *AppCatalogs) SetTeamsApps(value []TeamsAppable)() {
    err := m.GetBackingStore().Set("teamsApps", value)
    if err != nil {
        panic(err)
    }
}
type AppCatalogsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetTeamsApps()([]TeamsAppable)
    SetTeamsApps(value []TeamsAppable)()
}
