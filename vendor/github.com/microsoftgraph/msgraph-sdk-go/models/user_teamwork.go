package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UserTeamwork struct {
    Entity
}
// NewUserTeamwork instantiates a new UserTeamwork and sets the default values.
func NewUserTeamwork()(*UserTeamwork) {
    m := &UserTeamwork{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserTeamworkFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUserTeamworkFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserTeamwork(), nil
}
// GetAssociatedTeams gets the associatedTeams property value. The list of associatedTeamInfo objects that a user is associated with.
// returns a []AssociatedTeamInfoable when successful
func (m *UserTeamwork) GetAssociatedTeams()([]AssociatedTeamInfoable) {
    val, err := m.GetBackingStore().Get("associatedTeams")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AssociatedTeamInfoable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UserTeamwork) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["associatedTeams"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAssociatedTeamInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AssociatedTeamInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AssociatedTeamInfoable)
                }
            }
            m.SetAssociatedTeams(res)
        }
        return nil
    }
    res["installedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUserScopeTeamsAppInstallationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UserScopeTeamsAppInstallationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UserScopeTeamsAppInstallationable)
                }
            }
            m.SetInstalledApps(res)
        }
        return nil
    }
    res["locale"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocale(val)
        }
        return nil
    }
    res["region"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegion(val)
        }
        return nil
    }
    return res
}
// GetInstalledApps gets the installedApps property value. The apps installed in the personal scope of this user.
// returns a []UserScopeTeamsAppInstallationable when successful
func (m *UserTeamwork) GetInstalledApps()([]UserScopeTeamsAppInstallationable) {
    val, err := m.GetBackingStore().Get("installedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UserScopeTeamsAppInstallationable)
    }
    return nil
}
// GetLocale gets the locale property value. Represents the location that a user selected in Microsoft Teams and doesn't follow the Office's locale setting. A user’s locale is represented by their preferred language and country or region. For example, en-us. The language component follows two-letter codes as defined in ISO 639-1, and the country component follows two-letter codes as defined in ISO 3166-1 alpha-2.
// returns a *string when successful
func (m *UserTeamwork) GetLocale()(*string) {
    val, err := m.GetBackingStore().Get("locale")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegion gets the region property value. Represents the region of the organization or the user. For users with multigeo licenses, the property contains the user's region (if available). For users without multigeo licenses, the property contains the organization's region.The region value can be any region supported by the Teams payload. The possible values are: Americas, Europe and MiddleEast, Asia Pacific, UAE, Australia, Brazil, Canada, Switzerland, Germany, France, India, Japan, South Korea, Norway, Singapore, United Kingdom, South Africa, Sweden, Qatar, Poland, Italy, Israel, Spain, Mexico, USGov Community Cloud, USGov Community Cloud High, USGov Department of Defense, and China.
// returns a *string when successful
func (m *UserTeamwork) GetRegion()(*string) {
    val, err := m.GetBackingStore().Get("region")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UserTeamwork) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssociatedTeams() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssociatedTeams()))
        for i, v := range m.GetAssociatedTeams() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("associatedTeams", cast)
        if err != nil {
            return err
        }
    }
    if m.GetInstalledApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInstalledApps()))
        for i, v := range m.GetInstalledApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("installedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("locale", m.GetLocale())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("region", m.GetRegion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssociatedTeams sets the associatedTeams property value. The list of associatedTeamInfo objects that a user is associated with.
func (m *UserTeamwork) SetAssociatedTeams(value []AssociatedTeamInfoable)() {
    err := m.GetBackingStore().Set("associatedTeams", value)
    if err != nil {
        panic(err)
    }
}
// SetInstalledApps sets the installedApps property value. The apps installed in the personal scope of this user.
func (m *UserTeamwork) SetInstalledApps(value []UserScopeTeamsAppInstallationable)() {
    err := m.GetBackingStore().Set("installedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetLocale sets the locale property value. Represents the location that a user selected in Microsoft Teams and doesn't follow the Office's locale setting. A user’s locale is represented by their preferred language and country or region. For example, en-us. The language component follows two-letter codes as defined in ISO 639-1, and the country component follows two-letter codes as defined in ISO 3166-1 alpha-2.
func (m *UserTeamwork) SetLocale(value *string)() {
    err := m.GetBackingStore().Set("locale", value)
    if err != nil {
        panic(err)
    }
}
// SetRegion sets the region property value. Represents the region of the organization or the user. For users with multigeo licenses, the property contains the user's region (if available). For users without multigeo licenses, the property contains the organization's region.The region value can be any region supported by the Teams payload. The possible values are: Americas, Europe and MiddleEast, Asia Pacific, UAE, Australia, Brazil, Canada, Switzerland, Germany, France, India, Japan, South Korea, Norway, Singapore, United Kingdom, South Africa, Sweden, Qatar, Poland, Italy, Israel, Spain, Mexico, USGov Community Cloud, USGov Community Cloud High, USGov Department of Defense, and China.
func (m *UserTeamwork) SetRegion(value *string)() {
    err := m.GetBackingStore().Set("region", value)
    if err != nil {
        panic(err)
    }
}
type UserTeamworkable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssociatedTeams()([]AssociatedTeamInfoable)
    GetInstalledApps()([]UserScopeTeamsAppInstallationable)
    GetLocale()(*string)
    GetRegion()(*string)
    SetAssociatedTeams(value []AssociatedTeamInfoable)()
    SetInstalledApps(value []UserScopeTeamsAppInstallationable)()
    SetLocale(value *string)()
    SetRegion(value *string)()
}
