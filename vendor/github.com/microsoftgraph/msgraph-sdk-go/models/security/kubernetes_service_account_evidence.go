package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesServiceAccountEvidence struct {
    AlertEvidence
}
// NewKubernetesServiceAccountEvidence instantiates a new KubernetesServiceAccountEvidence and sets the default values.
func NewKubernetesServiceAccountEvidence()(*KubernetesServiceAccountEvidence) {
    m := &KubernetesServiceAccountEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesServiceAccountEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesServiceAccountEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesServiceAccountEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesServiceAccountEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesServiceAccountEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["namespace"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateKubernetesNamespaceEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNamespace(val.(KubernetesNamespaceEvidenceable))
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The service account name.
// returns a *string when successful
func (m *KubernetesServiceAccountEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNamespace gets the namespace property value. The service account namespace.
// returns a KubernetesNamespaceEvidenceable when successful
func (m *KubernetesServiceAccountEvidence) GetNamespace()(KubernetesNamespaceEvidenceable) {
    val, err := m.GetBackingStore().Get("namespace")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesNamespaceEvidenceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesServiceAccountEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("namespace", m.GetNamespace())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetName sets the name property value. The service account name.
func (m *KubernetesServiceAccountEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNamespace sets the namespace property value. The service account namespace.
func (m *KubernetesServiceAccountEvidence) SetNamespace(value KubernetesNamespaceEvidenceable)() {
    err := m.GetBackingStore().Set("namespace", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesServiceAccountEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetName()(*string)
    GetNamespace()(KubernetesNamespaceEvidenceable)
    SetName(value *string)()
    SetNamespace(value KubernetesNamespaceEvidenceable)()
}
