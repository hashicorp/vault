package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnclassifiedArtifact struct {
    Artifact
}
// NewUnclassifiedArtifact instantiates a new UnclassifiedArtifact and sets the default values.
func NewUnclassifiedArtifact()(*UnclassifiedArtifact) {
    m := &UnclassifiedArtifact{
        Artifact: *NewArtifact(),
    }
    odataTypeValue := "#microsoft.graph.security.unclassifiedArtifact"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUnclassifiedArtifactFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnclassifiedArtifactFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnclassifiedArtifact(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnclassifiedArtifact) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Artifact.GetFieldDeserializers()
    res["kind"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKind(val)
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
// GetKind gets the kind property value. The kind for this unclassifiedArtifact resource, describing what this value means.
// returns a *string when successful
func (m *UnclassifiedArtifact) GetKind()(*string) {
    val, err := m.GetBackingStore().Get("kind")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetValue gets the value property value. The value for this unclassifiedArtifact.
// returns a *string when successful
func (m *UnclassifiedArtifact) GetValue()(*string) {
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
func (m *UnclassifiedArtifact) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Artifact.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("kind", m.GetKind())
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
// SetKind sets the kind property value. The kind for this unclassifiedArtifact resource, describing what this value means.
func (m *UnclassifiedArtifact) SetKind(value *string)() {
    err := m.GetBackingStore().Set("kind", value)
    if err != nil {
        panic(err)
    }
}
// SetValue sets the value property value. The value for this unclassifiedArtifact.
func (m *UnclassifiedArtifact) SetValue(value *string)() {
    err := m.GetBackingStore().Set("value", value)
    if err != nil {
        panic(err)
    }
}
type UnclassifiedArtifactable interface {
    Artifactable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetKind()(*string)
    GetValue()(*string)
    SetKind(value *string)()
    SetValue(value *string)()
}
