package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharePointRestoreSession struct {
    RestoreSessionBase
}
// NewSharePointRestoreSession instantiates a new SharePointRestoreSession and sets the default values.
func NewSharePointRestoreSession()(*SharePointRestoreSession) {
    m := &SharePointRestoreSession{
        RestoreSessionBase: *NewRestoreSessionBase(),
    }
    odataTypeValue := "#microsoft.graph.sharePointRestoreSession"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSharePointRestoreSessionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharePointRestoreSessionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharePointRestoreSession(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharePointRestoreSession) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RestoreSessionBase.GetFieldDeserializers()
    res["siteRestoreArtifacts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteRestoreArtifactFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SiteRestoreArtifactable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SiteRestoreArtifactable)
                }
            }
            m.SetSiteRestoreArtifacts(res)
        }
        return nil
    }
    return res
}
// GetSiteRestoreArtifacts gets the siteRestoreArtifacts property value. A collection of restore points and destination details that can be used to restore SharePoint sites.
// returns a []SiteRestoreArtifactable when successful
func (m *SharePointRestoreSession) GetSiteRestoreArtifacts()([]SiteRestoreArtifactable) {
    val, err := m.GetBackingStore().Get("siteRestoreArtifacts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SiteRestoreArtifactable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharePointRestoreSession) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RestoreSessionBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetSiteRestoreArtifacts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSiteRestoreArtifacts()))
        for i, v := range m.GetSiteRestoreArtifacts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("siteRestoreArtifacts", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSiteRestoreArtifacts sets the siteRestoreArtifacts property value. A collection of restore points and destination details that can be used to restore SharePoint sites.
func (m *SharePointRestoreSession) SetSiteRestoreArtifacts(value []SiteRestoreArtifactable)() {
    err := m.GetBackingStore().Set("siteRestoreArtifacts", value)
    if err != nil {
        panic(err)
    }
}
type SharePointRestoreSessionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RestoreSessionBaseable
    GetSiteRestoreArtifacts()([]SiteRestoreArtifactable)
    SetSiteRestoreArtifacts(value []SiteRestoreArtifactable)()
}
