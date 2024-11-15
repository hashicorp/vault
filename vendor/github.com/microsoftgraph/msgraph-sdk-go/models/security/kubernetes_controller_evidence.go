package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesControllerEvidence struct {
    AlertEvidence
}
// NewKubernetesControllerEvidence instantiates a new KubernetesControllerEvidence and sets the default values.
func NewKubernetesControllerEvidence()(*KubernetesControllerEvidence) {
    m := &KubernetesControllerEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesControllerEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesControllerEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesControllerEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesControllerEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesControllerEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["labels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDictionaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabels(val.(Dictionaryable))
        }
        return nil
    }
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
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val)
        }
        return nil
    }
    return res
}
// GetLabels gets the labels property value. The labels for the Kubernetes pod.
// returns a Dictionaryable when successful
func (m *KubernetesControllerEvidence) GetLabels()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("labels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetName gets the name property value. The controller name.
// returns a *string when successful
func (m *KubernetesControllerEvidence) GetName()(*string) {
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
func (m *KubernetesControllerEvidence) GetNamespace()(KubernetesNamespaceEvidenceable) {
    val, err := m.GetBackingStore().Get("namespace")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesNamespaceEvidenceable)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The controller type.
// returns a *string when successful
func (m *KubernetesControllerEvidence) GetTypeEscaped()(*string) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesControllerEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("labels", m.GetLabels())
        if err != nil {
            return err
        }
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
    {
        err = writer.WriteStringValue("type", m.GetTypeEscaped())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLabels sets the labels property value. The labels for the Kubernetes pod.
func (m *KubernetesControllerEvidence) SetLabels(value Dictionaryable)() {
    err := m.GetBackingStore().Set("labels", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The controller name.
func (m *KubernetesControllerEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNamespace sets the namespace property value. The service account namespace.
func (m *KubernetesControllerEvidence) SetNamespace(value KubernetesNamespaceEvidenceable)() {
    err := m.GetBackingStore().Set("namespace", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The controller type.
func (m *KubernetesControllerEvidence) SetTypeEscaped(value *string)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesControllerEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLabels()(Dictionaryable)
    GetName()(*string)
    GetNamespace()(KubernetesNamespaceEvidenceable)
    GetTypeEscaped()(*string)
    SetLabels(value Dictionaryable)()
    SetName(value *string)()
    SetNamespace(value KubernetesNamespaceEvidenceable)()
    SetTypeEscaped(value *string)()
}
