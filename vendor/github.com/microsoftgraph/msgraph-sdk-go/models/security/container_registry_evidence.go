package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ContainerRegistryEvidence struct {
    AlertEvidence
}
// NewContainerRegistryEvidence instantiates a new ContainerRegistryEvidence and sets the default values.
func NewContainerRegistryEvidence()(*ContainerRegistryEvidence) {
    m := &ContainerRegistryEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.containerRegistryEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateContainerRegistryEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateContainerRegistryEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContainerRegistryEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ContainerRegistryEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["registry"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistry(val)
        }
        return nil
    }
    return res
}
// GetRegistry gets the registry property value. The registry URI.
// returns a *string when successful
func (m *ContainerRegistryEvidence) GetRegistry()(*string) {
    val, err := m.GetBackingStore().Get("registry")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ContainerRegistryEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("registry", m.GetRegistry())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRegistry sets the registry property value. The registry URI.
func (m *ContainerRegistryEvidence) SetRegistry(value *string)() {
    err := m.GetBackingStore().Set("registry", value)
    if err != nil {
        panic(err)
    }
}
type ContainerRegistryEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRegistry()(*string)
    SetRegistry(value *string)()
}
