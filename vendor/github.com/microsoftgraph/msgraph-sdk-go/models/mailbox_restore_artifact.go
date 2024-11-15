package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailboxRestoreArtifact struct {
    RestoreArtifactBase
}
// NewMailboxRestoreArtifact instantiates a new MailboxRestoreArtifact and sets the default values.
func NewMailboxRestoreArtifact()(*MailboxRestoreArtifact) {
    m := &MailboxRestoreArtifact{
        RestoreArtifactBase: *NewRestoreArtifactBase(),
    }
    return m
}
// CreateMailboxRestoreArtifactFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailboxRestoreArtifactFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.granularMailboxRestoreArtifact":
                        return NewGranularMailboxRestoreArtifact(), nil
                }
            }
        }
    }
    return NewMailboxRestoreArtifact(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailboxRestoreArtifact) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RestoreArtifactBase.GetFieldDeserializers()
    res["restoredFolderId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestoredFolderId(val)
        }
        return nil
    }
    res["restoredFolderName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRestoredFolderName(val)
        }
        return nil
    }
    return res
}
// GetRestoredFolderId gets the restoredFolderId property value. The new restored folder identifier for the user.
// returns a *string when successful
func (m *MailboxRestoreArtifact) GetRestoredFolderId()(*string) {
    val, err := m.GetBackingStore().Get("restoredFolderId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRestoredFolderName gets the restoredFolderName property value. The new restored folder name.
// returns a *string when successful
func (m *MailboxRestoreArtifact) GetRestoredFolderName()(*string) {
    val, err := m.GetBackingStore().Get("restoredFolderName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailboxRestoreArtifact) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RestoreArtifactBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("restoredFolderId", m.GetRestoredFolderId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRestoredFolderId sets the restoredFolderId property value. The new restored folder identifier for the user.
func (m *MailboxRestoreArtifact) SetRestoredFolderId(value *string)() {
    err := m.GetBackingStore().Set("restoredFolderId", value)
    if err != nil {
        panic(err)
    }
}
// SetRestoredFolderName sets the restoredFolderName property value. The new restored folder name.
func (m *MailboxRestoreArtifact) SetRestoredFolderName(value *string)() {
    err := m.GetBackingStore().Set("restoredFolderName", value)
    if err != nil {
        panic(err)
    }
}
type MailboxRestoreArtifactable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RestoreArtifactBaseable
    GetRestoredFolderId()(*string)
    GetRestoredFolderName()(*string)
    SetRestoredFolderId(value *string)()
    SetRestoredFolderName(value *string)()
}
