package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SiteRestoreArtifact struct {
    RestoreArtifactBase
}
// NewSiteRestoreArtifact instantiates a new SiteRestoreArtifact and sets the default values.
func NewSiteRestoreArtifact()(*SiteRestoreArtifact) {
    m := &SiteRestoreArtifact{
        RestoreArtifactBase: *NewRestoreArtifactBase(),
    }
    return m
}
// CreateSiteRestoreArtifactFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSiteRestoreArtifactFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSiteRestoreArtifact(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SiteRestoreArtifact) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RestoreArtifactBase.GetFieldDeserializers()
    res["restoredSiteId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestoredSiteId(val)
        }
        return nil
    }
    res["restoredSiteName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestoredSiteName(val)
        }
        return nil
    }
    res["restoredSiteWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestoredSiteWebUrl(val)
        }
        return nil
    }
    return res
}
// GetRestoredSiteId gets the restoredSiteId property value. The new site identifier if the value of the destinationType property is new, and the existing site ID if the value is inPlace.
// returns a *string when successful
func (m *SiteRestoreArtifact) GetRestoredSiteId()(*string) {
    val, err := m.GetBackingStore().Get("restoredSiteId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRestoredSiteName gets the restoredSiteName property value. The name of the restored site.
// returns a *string when successful
func (m *SiteRestoreArtifact) GetRestoredSiteName()(*string) {
    val, err := m.GetBackingStore().Get("restoredSiteName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRestoredSiteWebUrl gets the restoredSiteWebUrl property value. The web URL of the restored site.
// returns a *string when successful
func (m *SiteRestoreArtifact) GetRestoredSiteWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("restoredSiteWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SiteRestoreArtifact) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RestoreArtifactBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("restoredSiteId", m.GetRestoredSiteId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRestoredSiteId sets the restoredSiteId property value. The new site identifier if the value of the destinationType property is new, and the existing site ID if the value is inPlace.
func (m *SiteRestoreArtifact) SetRestoredSiteId(value *string)() {
    err := m.GetBackingStore().Set("restoredSiteId", value)
    if err != nil {
        panic(err)
    }
}
// SetRestoredSiteName sets the restoredSiteName property value. The name of the restored site.
func (m *SiteRestoreArtifact) SetRestoredSiteName(value *string)() {
    err := m.GetBackingStore().Set("restoredSiteName", value)
    if err != nil {
        panic(err)
    }
}
// SetRestoredSiteWebUrl sets the restoredSiteWebUrl property value. The web URL of the restored site.
func (m *SiteRestoreArtifact) SetRestoredSiteWebUrl(value *string)() {
    err := m.GetBackingStore().Set("restoredSiteWebUrl", value)
    if err != nil {
        panic(err)
    }
}
type SiteRestoreArtifactable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RestoreArtifactBaseable
    GetRestoredSiteId()(*string)
    GetRestoredSiteName()(*string)
    GetRestoredSiteWebUrl()(*string)
    SetRestoredSiteId(value *string)()
    SetRestoredSiteName(value *string)()
    SetRestoredSiteWebUrl(value *string)()
}
