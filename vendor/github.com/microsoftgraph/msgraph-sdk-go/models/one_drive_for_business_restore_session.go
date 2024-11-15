package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OneDriveForBusinessRestoreSession struct {
    RestoreSessionBase
}
// NewOneDriveForBusinessRestoreSession instantiates a new OneDriveForBusinessRestoreSession and sets the default values.
func NewOneDriveForBusinessRestoreSession()(*OneDriveForBusinessRestoreSession) {
    m := &OneDriveForBusinessRestoreSession{
        RestoreSessionBase: *NewRestoreSessionBase(),
    }
    odataTypeValue := "#microsoft.graph.oneDriveForBusinessRestoreSession"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOneDriveForBusinessRestoreSessionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOneDriveForBusinessRestoreSessionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOneDriveForBusinessRestoreSession(), nil
}
// GetDriveRestoreArtifacts gets the driveRestoreArtifacts property value. A collection of restore points and destination details that can be used to restore a OneDrive for Business drive.
// returns a []DriveRestoreArtifactable when successful
func (m *OneDriveForBusinessRestoreSession) GetDriveRestoreArtifacts()([]DriveRestoreArtifactable) {
    val, err := m.GetBackingStore().Get("driveRestoreArtifacts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveRestoreArtifactable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OneDriveForBusinessRestoreSession) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RestoreSessionBase.GetFieldDeserializers()
    res["driveRestoreArtifacts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveRestoreArtifactFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveRestoreArtifactable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveRestoreArtifactable)
                }
            }
            m.SetDriveRestoreArtifacts(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *OneDriveForBusinessRestoreSession) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RestoreSessionBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDriveRestoreArtifacts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDriveRestoreArtifacts()))
        for i, v := range m.GetDriveRestoreArtifacts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("driveRestoreArtifacts", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDriveRestoreArtifacts sets the driveRestoreArtifacts property value. A collection of restore points and destination details that can be used to restore a OneDrive for Business drive.
func (m *OneDriveForBusinessRestoreSession) SetDriveRestoreArtifacts(value []DriveRestoreArtifactable)() {
    err := m.GetBackingStore().Set("driveRestoreArtifacts", value)
    if err != nil {
        panic(err)
    }
}
type OneDriveForBusinessRestoreSessionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RestoreSessionBaseable
    GetDriveRestoreArtifacts()([]DriveRestoreArtifactable)
    SetDriveRestoreArtifacts(value []DriveRestoreArtifactable)()
}
