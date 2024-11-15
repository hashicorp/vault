package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type RegistryKeyEvidence struct {
    AlertEvidence
}
// NewRegistryKeyEvidence instantiates a new RegistryKeyEvidence and sets the default values.
func NewRegistryKeyEvidence()(*RegistryKeyEvidence) {
    m := &RegistryKeyEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.registryKeyEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateRegistryKeyEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRegistryKeyEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRegistryKeyEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RegistryKeyEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["registryHive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryHive(val)
        }
        return nil
    }
    res["registryKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistryKey(val)
        }
        return nil
    }
    return res
}
// GetRegistryHive gets the registryHive property value. Registry hive of the key that the recorded action was applied to.
// returns a *string when successful
func (m *RegistryKeyEvidence) GetRegistryHive()(*string) {
    val, err := m.GetBackingStore().Get("registryHive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistryKey gets the registryKey property value. Registry key that the recorded action was applied to.
// returns a *string when successful
func (m *RegistryKeyEvidence) GetRegistryKey()(*string) {
    val, err := m.GetBackingStore().Get("registryKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RegistryKeyEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("registryHive", m.GetRegistryHive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registryKey", m.GetRegistryKey())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRegistryHive sets the registryHive property value. Registry hive of the key that the recorded action was applied to.
func (m *RegistryKeyEvidence) SetRegistryHive(value *string)() {
    err := m.GetBackingStore().Set("registryHive", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistryKey sets the registryKey property value. Registry key that the recorded action was applied to.
func (m *RegistryKeyEvidence) SetRegistryKey(value *string)() {
    err := m.GetBackingStore().Set("registryKey", value)
    if err != nil {
        panic(err)
    }
}
type RegistryKeyEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRegistryHive()(*string)
    GetRegistryKey()(*string)
    SetRegistryHive(value *string)()
    SetRegistryKey(value *string)()
}
