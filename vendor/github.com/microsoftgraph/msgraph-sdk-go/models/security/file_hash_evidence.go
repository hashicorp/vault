package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FileHashEvidence struct {
    AlertEvidence
}
// NewFileHashEvidence instantiates a new FileHashEvidence and sets the default values.
func NewFileHashEvidence()(*FileHashEvidence) {
    m := &FileHashEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.fileHashEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateFileHashEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileHashEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileHashEvidence(), nil
}
// GetAlgorithm gets the algorithm property value. The algorithm property
// returns a *FileHashAlgorithm when successful
func (m *FileHashEvidence) GetAlgorithm()(*FileHashAlgorithm) {
    val, err := m.GetBackingStore().Get("algorithm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FileHashAlgorithm)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FileHashEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["algorithm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFileHashAlgorithm)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlgorithm(val.(*FileHashAlgorithm))
        }
        return nil
    }
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValue(val)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
// returns a *string when successful
func (m *FileHashEvidence) GetValue()(*string) {
    val, err := m.GetBackingStore().Get("value")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FileHashEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAlgorithm() != nil {
        cast := (*m.GetAlgorithm()).String()
        err = writer.WriteStringValue("algorithm", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAlgorithm sets the algorithm property value. The algorithm property
func (m *FileHashEvidence) SetAlgorithm(value *FileHashAlgorithm)() {
    err := m.GetBackingStore().Set("algorithm", value)
    if err != nil {
        panic(err)
    }
}
// SetValue sets the value property value. The value property
func (m *FileHashEvidence) SetValue(value *string)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
type FileHashEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlgorithm()(*FileHashAlgorithm)
    GetValue()(*string)
    SetAlgorithm(value *FileHashAlgorithm)()
    SetValue(value *string)()
}
