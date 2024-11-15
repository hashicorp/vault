package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsAppSettings struct {
    Entity
}
// NewTeamsAppSettings instantiates a new TeamsAppSettings and sets the default values.
func NewTeamsAppSettings()(*TeamsAppSettings) {
    m := &TeamsAppSettings{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsAppSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsAppSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsAppSettings(), nil
}
// GetAllowUserRequestsForAppAccess gets the allowUserRequestsForAppAccess property value. Indicates whether users are allowed to request access to the unavailable Teams apps.
// returns a *bool when successful
func (m *TeamsAppSettings) GetAllowUserRequestsForAppAccess()(*bool) {
    val, err := m.GetBackingStore().Get("allowUserRequestsForAppAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsAppSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowUserRequestsForAppAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowUserRequestsForAppAccess(val)
        }
        return nil
    }
    res["isUserPersonalScopeResourceSpecificConsentEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsUserPersonalScopeResourceSpecificConsentEnabled(val)
        }
        return nil
    }
    return res
}
// GetIsUserPersonalScopeResourceSpecificConsentEnabled gets the isUserPersonalScopeResourceSpecificConsentEnabled property value. Indicates whether resource-specific consent for personal scope in Teams apps is enabled for the tenant. True indicates that Teams apps that are allowed in the tenant and require resource-specific permissions can be installed in the personal scope. False blocks the installation of any Teams app that requires resource-specific permissions in the personal scope.
// returns a *bool when successful
func (m *TeamsAppSettings) GetIsUserPersonalScopeResourceSpecificConsentEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isUserPersonalScopeResourceSpecificConsentEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamsAppSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowUserRequestsForAppAccess", m.GetAllowUserRequestsForAppAccess())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isUserPersonalScopeResourceSpecificConsentEnabled", m.GetIsUserPersonalScopeResourceSpecificConsentEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowUserRequestsForAppAccess sets the allowUserRequestsForAppAccess property value. Indicates whether users are allowed to request access to the unavailable Teams apps.
func (m *TeamsAppSettings) SetAllowUserRequestsForAppAccess(value *bool)() {
    err := m.GetBackingStore().Set("allowUserRequestsForAppAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetIsUserPersonalScopeResourceSpecificConsentEnabled sets the isUserPersonalScopeResourceSpecificConsentEnabled property value. Indicates whether resource-specific consent for personal scope in Teams apps is enabled for the tenant. True indicates that Teams apps that are allowed in the tenant and require resource-specific permissions can be installed in the personal scope. False blocks the installation of any Teams app that requires resource-specific permissions in the personal scope.
func (m *TeamsAppSettings) SetIsUserPersonalScopeResourceSpecificConsentEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isUserPersonalScopeResourceSpecificConsentEnabled", value)
    if err != nil {
        panic(err)
    }
}
type TeamsAppSettingsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowUserRequestsForAppAccess()(*bool)
    GetIsUserPersonalScopeResourceSpecificConsentEnabled()(*bool)
    SetAllowUserRequestsForAppAccess(value *bool)()
    SetIsUserPersonalScopeResourceSpecificConsentEnabled(value *bool)()
}
