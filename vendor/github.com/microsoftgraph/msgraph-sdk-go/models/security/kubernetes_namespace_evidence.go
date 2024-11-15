package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesNamespaceEvidence struct {
    AlertEvidence
}
// NewKubernetesNamespaceEvidence instantiates a new KubernetesNamespaceEvidence and sets the default values.
func NewKubernetesNamespaceEvidence()(*KubernetesNamespaceEvidence) {
    m := &KubernetesNamespaceEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesNamespaceEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesNamespaceEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesNamespaceEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesNamespaceEvidence(), nil
}
// GetCluster gets the cluster property value. The namespace cluster.
// returns a KubernetesClusterEvidenceable when successful
func (m *KubernetesNamespaceEvidence) GetCluster()(KubernetesClusterEvidenceable) {
    val, err := m.GetBackingStore().Get("cluster")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesClusterEvidenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesNamespaceEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["cluster"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateKubernetesClusterEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCluster(val.(KubernetesClusterEvidenceable))
        }
        return nil
    }
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
    return res
}
// GetLabels gets the labels property value. The labels for the Kubernetes pod.
// returns a Dictionaryable when successful
func (m *KubernetesNamespaceEvidence) GetLabels()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("labels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetName gets the name property value. The namespace name.
// returns a *string when successful
func (m *KubernetesNamespaceEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesNamespaceEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("cluster", m.GetCluster())
        if err != nil {
            return err
        }
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
    return nil
}
// SetCluster sets the cluster property value. The namespace cluster.
func (m *KubernetesNamespaceEvidence) SetCluster(value KubernetesClusterEvidenceable)() {
    err := m.GetBackingStore().Set("cluster", value)
    if err != nil {
        panic(err)
    }
}
// SetLabels sets the labels property value. The labels for the Kubernetes pod.
func (m *KubernetesNamespaceEvidence) SetLabels(value Dictionaryable)() {
    err := m.GetBackingStore().Set("labels", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The namespace name.
func (m *KubernetesNamespaceEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesNamespaceEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCluster()(KubernetesClusterEvidenceable)
    GetLabels()(Dictionaryable)
    GetName()(*string)
    SetCluster(value KubernetesClusterEvidenceable)()
    SetLabels(value Dictionaryable)()
    SetName(value *string)()
}
