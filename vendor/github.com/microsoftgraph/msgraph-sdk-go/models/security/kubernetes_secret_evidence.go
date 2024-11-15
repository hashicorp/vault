package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesSecretEvidence struct {
    AlertEvidence
}
// NewKubernetesSecretEvidence instantiates a new KubernetesSecretEvidence and sets the default values.
func NewKubernetesSecretEvidence()(*KubernetesSecretEvidence) {
    m := &KubernetesSecretEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesSecretEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesSecretEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesSecretEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesSecretEvidence(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesSecretEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["secretType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecretType(val)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The secret name.
// returns a *string when successful
func (m *KubernetesSecretEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNamespace gets the namespace property value. The secret namespace.
// returns a KubernetesNamespaceEvidenceable when successful
func (m *KubernetesSecretEvidence) GetNamespace()(KubernetesNamespaceEvidenceable) {
    val, err := m.GetBackingStore().Get("namespace")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesNamespaceEvidenceable)
    }
    return nil
}
// GetSecretType gets the secretType property value. The secret type can include both built-in types and custom ones. Examples of built-in types are: Opaque, kubernetes.io/service-account-token, kubernetes.io/dockercfg, kubernetes.io/dockerconfigjson, kubernetes.io/basic-auth, kubernetes.io/ssh-auth, kubernetes.io/tls, bootstrap.kubernetes.io/token.
// returns a *string when successful
func (m *KubernetesSecretEvidence) GetSecretType()(*string) {
    val, err := m.GetBackingStore().Get("secretType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesSecretEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteStringValue("secretType", m.GetSecretType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetName sets the name property value. The secret name.
func (m *KubernetesSecretEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNamespace sets the namespace property value. The secret namespace.
func (m *KubernetesSecretEvidence) SetNamespace(value KubernetesNamespaceEvidenceable)() {
    err := m.GetBackingStore().Set("namespace", value)
    if err != nil {
        panic(err)
    }
}
// SetSecretType sets the secretType property value. The secret type can include both built-in types and custom ones. Examples of built-in types are: Opaque, kubernetes.io/service-account-token, kubernetes.io/dockercfg, kubernetes.io/dockerconfigjson, kubernetes.io/basic-auth, kubernetes.io/ssh-auth, kubernetes.io/tls, bootstrap.kubernetes.io/token.
func (m *KubernetesSecretEvidence) SetSecretType(value *string)() {
    err := m.GetBackingStore().Set("secretType", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesSecretEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetName()(*string)
    GetNamespace()(KubernetesNamespaceEvidenceable)
    GetSecretType()(*string)
    SetName(value *string)()
    SetNamespace(value KubernetesNamespaceEvidenceable)()
    SetSecretType(value *string)()
}
