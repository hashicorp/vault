package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsAppInstallation struct {
    Entity
}
// NewTeamsAppInstallation instantiates a new TeamsAppInstallation and sets the default values.
func NewTeamsAppInstallation()(*TeamsAppInstallation) {
    m := &TeamsAppInstallation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsAppInstallationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsAppInstallationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.userScopeTeamsAppInstallation":
                        return NewUserScopeTeamsAppInstallation(), nil
                }
            }
        }
    }
    return NewTeamsAppInstallation(), nil
}
// GetConsentedPermissionSet gets the consentedPermissionSet property value. The set of resource-specific permissions consented to while installing or upgrading the teamsApp.
// returns a TeamsAppPermissionSetable when successful
func (m *TeamsAppInstallation) GetConsentedPermissionSet()(TeamsAppPermissionSetable) {
    val, err := m.GetBackingStore().Get("consentedPermissionSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamsAppPermissionSetable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsAppInstallation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["consentedPermissionSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamsAppPermissionSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConsentedPermissionSet(val.(TeamsAppPermissionSetable))
        }
        return nil
    }
    res["teamsApp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamsAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamsApp(val.(TeamsAppable))
        }
        return nil
    }
    res["teamsAppDefinition"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTeamsAppDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamsAppDefinition(val.(TeamsAppDefinitionable))
        }
        return nil
    }
    return res
}
// GetTeamsApp gets the teamsApp property value. The app that is installed.
// returns a TeamsAppable when successful
func (m *TeamsAppInstallation) GetTeamsApp()(TeamsAppable) {
    val, err := m.GetBackingStore().Get("teamsApp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamsAppable)
    }
    return nil
}
// GetTeamsAppDefinition gets the teamsAppDefinition property value. The details of this version of the app.
// returns a TeamsAppDefinitionable when successful
func (m *TeamsAppInstallation) GetTeamsAppDefinition()(TeamsAppDefinitionable) {
    val, err := m.GetBackingStore().Get("teamsAppDefinition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TeamsAppDefinitionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamsAppInstallation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("consentedPermissionSet", m.GetConsentedPermissionSet())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("teamsApp", m.GetTeamsApp())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("teamsAppDefinition", m.GetTeamsAppDefinition())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConsentedPermissionSet sets the consentedPermissionSet property value. The set of resource-specific permissions consented to while installing or upgrading the teamsApp.
func (m *TeamsAppInstallation) SetConsentedPermissionSet(value TeamsAppPermissionSetable)() {
    err := m.GetBackingStore().Set("consentedPermissionSet", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamsApp sets the teamsApp property value. The app that is installed.
func (m *TeamsAppInstallation) SetTeamsApp(value TeamsAppable)() {
    err := m.GetBackingStore().Set("teamsApp", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamsAppDefinition sets the teamsAppDefinition property value. The details of this version of the app.
func (m *TeamsAppInstallation) SetTeamsAppDefinition(value TeamsAppDefinitionable)() {
    err := m.GetBackingStore().Set("teamsAppDefinition", value)
    if err != nil {
        panic(err)
    }
}
type TeamsAppInstallationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConsentedPermissionSet()(TeamsAppPermissionSetable)
    GetTeamsApp()(TeamsAppable)
    GetTeamsAppDefinition()(TeamsAppDefinitionable)
    SetConsentedPermissionSet(value TeamsAppPermissionSetable)()
    SetTeamsApp(value TeamsAppable)()
    SetTeamsAppDefinition(value TeamsAppDefinitionable)()
}
