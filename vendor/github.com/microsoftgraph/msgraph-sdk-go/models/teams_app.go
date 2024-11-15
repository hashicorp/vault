package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsApp struct {
    Entity
}
// NewTeamsApp instantiates a new TeamsApp and sets the default values.
func NewTeamsApp()(*TeamsApp) {
    m := &TeamsApp{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsApp(), nil
}
// GetAppDefinitions gets the appDefinitions property value. The details for each version of the app.
// returns a []TeamsAppDefinitionable when successful
func (m *TeamsApp) GetAppDefinitions()([]TeamsAppDefinitionable) {
    val, err := m.GetBackingStore().Get("appDefinitions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TeamsAppDefinitionable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the catalog app provided by the app developer in the Microsoft Teams zip app package.
// returns a *string when successful
func (m *TeamsApp) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDistributionMethod gets the distributionMethod property value. The method of distribution for the app. Read-only.
// returns a *TeamsAppDistributionMethod when successful
func (m *TeamsApp) GetDistributionMethod()(*TeamsAppDistributionMethod) {
    val, err := m.GetBackingStore().Get("distributionMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamsAppDistributionMethod)
    }
    return nil
}
// GetExternalId gets the externalId property value. The ID of the catalog provided by the app developer in the Microsoft Teams zip app package.
// returns a *string when successful
func (m *TeamsApp) GetExternalId()(*string) {
    val, err := m.GetBackingStore().Get("externalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["appDefinitions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTeamsAppDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TeamsAppDefinitionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TeamsAppDefinitionable)
                }
            }
            m.SetAppDefinitions(res)
        }
        return nil
    }
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
    res["distributionMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamsAppDistributionMethod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDistributionMethod(val.(*TeamsAppDistributionMethod))
        }
        return nil
    }
    res["externalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *TeamsApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAppDefinitions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppDefinitions()))
        for i, v := range m.GetAppDefinitions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appDefinitions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetDistributionMethod() != nil {
        cast := (*m.GetDistributionMethod()).String()
        err = writer.WriteStringValue("distributionMethod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("externalId", m.GetExternalId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppDefinitions sets the appDefinitions property value. The details for each version of the app.
func (m *TeamsApp) SetAppDefinitions(value []TeamsAppDefinitionable)() {
    err := m.GetBackingStore().Set("appDefinitions", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the catalog app provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDistributionMethod sets the distributionMethod property value. The method of distribution for the app. Read-only.
func (m *TeamsApp) SetDistributionMethod(value *TeamsAppDistributionMethod)() {
    err := m.GetBackingStore().Set("distributionMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalId sets the externalId property value. The ID of the catalog provided by the app developer in the Microsoft Teams zip app package.
func (m *TeamsApp) SetExternalId(value *string)() {
    err := m.GetBackingStore().Set("externalId", value)
    if err != nil {
        panic(err)
    }
}
type TeamsAppable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppDefinitions()([]TeamsAppDefinitionable)
    GetDisplayName()(*string)
    GetDistributionMethod()(*TeamsAppDistributionMethod)
    GetExternalId()(*string)
    SetAppDefinitions(value []TeamsAppDefinitionable)()
    SetDisplayName(value *string)()
    SetDistributionMethod(value *TeamsAppDistributionMethod)()
    SetExternalId(value *string)()
}
