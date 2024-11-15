package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesPodEvidence struct {
    AlertEvidence
}
// NewKubernetesPodEvidence instantiates a new KubernetesPodEvidence and sets the default values.
func NewKubernetesPodEvidence()(*KubernetesPodEvidence) {
    m := &KubernetesPodEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesPodEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesPodEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesPodEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesPodEvidence(), nil
}
// GetContainers gets the containers property value. The list of pod containers which are not init or ephemeral containers.
// returns a []ContainerEvidenceable when successful
func (m *KubernetesPodEvidence) GetContainers()([]ContainerEvidenceable) {
    val, err := m.GetBackingStore().Get("containers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContainerEvidenceable)
    }
    return nil
}
// GetController gets the controller property value. The pod controller.
// returns a KubernetesControllerEvidenceable when successful
func (m *KubernetesPodEvidence) GetController()(KubernetesControllerEvidenceable) {
    val, err := m.GetBackingStore().Get("controller")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesControllerEvidenceable)
    }
    return nil
}
// GetEphemeralContainers gets the ephemeralContainers property value. The list of pod ephemeral containers.
// returns a []ContainerEvidenceable when successful
func (m *KubernetesPodEvidence) GetEphemeralContainers()([]ContainerEvidenceable) {
    val, err := m.GetBackingStore().Get("ephemeralContainers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContainerEvidenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesPodEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["containers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContainerEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContainerEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContainerEvidenceable)
                }
            }
            m.SetContainers(res)
        }
        return nil
    }
    res["controller"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateKubernetesControllerEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetController(val.(KubernetesControllerEvidenceable))
        }
        return nil
    }
    res["ephemeralContainers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContainerEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContainerEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContainerEvidenceable)
                }
            }
            m.SetEphemeralContainers(res)
        }
        return nil
    }
    res["initContainers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateContainerEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ContainerEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ContainerEvidenceable)
                }
            }
            m.SetInitContainers(res)
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
    res["podIp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPodIp(val.(IpEvidenceable))
        }
        return nil
    }
    res["serviceAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateKubernetesServiceAccountEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceAccount(val.(KubernetesServiceAccountEvidenceable))
        }
        return nil
    }
    return res
}
// GetInitContainers gets the initContainers property value. The list of pod init containers.
// returns a []ContainerEvidenceable when successful
func (m *KubernetesPodEvidence) GetInitContainers()([]ContainerEvidenceable) {
    val, err := m.GetBackingStore().Get("initContainers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ContainerEvidenceable)
    }
    return nil
}
// GetLabels gets the labels property value. The pod labels.
// returns a Dictionaryable when successful
func (m *KubernetesPodEvidence) GetLabels()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("labels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetName gets the name property value. The pod name.
// returns a *string when successful
func (m *KubernetesPodEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNamespace gets the namespace property value. The pod namespace.
// returns a KubernetesNamespaceEvidenceable when successful
func (m *KubernetesPodEvidence) GetNamespace()(KubernetesNamespaceEvidenceable) {
    val, err := m.GetBackingStore().Get("namespace")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesNamespaceEvidenceable)
    }
    return nil
}
// GetPodIp gets the podIp property value. The pod IP.
// returns a IpEvidenceable when successful
func (m *KubernetesPodEvidence) GetPodIp()(IpEvidenceable) {
    val, err := m.GetBackingStore().Get("podIp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IpEvidenceable)
    }
    return nil
}
// GetServiceAccount gets the serviceAccount property value. The pod service account.
// returns a KubernetesServiceAccountEvidenceable when successful
func (m *KubernetesPodEvidence) GetServiceAccount()(KubernetesServiceAccountEvidenceable) {
    val, err := m.GetBackingStore().Get("serviceAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesServiceAccountEvidenceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesPodEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetContainers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetContainers()))
        for i, v := range m.GetContainers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("containers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("controller", m.GetController())
        if err != nil {
            return err
        }
    }
    if m.GetEphemeralContainers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEphemeralContainers()))
        for i, v := range m.GetEphemeralContainers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ephemeralContainers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetInitContainers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInitContainers()))
        for i, v := range m.GetInitContainers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("initContainers", cast)
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
    {
        err = writer.WriteObjectValue("namespace", m.GetNamespace())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("podIp", m.GetPodIp())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("serviceAccount", m.GetServiceAccount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContainers sets the containers property value. The list of pod containers which are not init or ephemeral containers.
func (m *KubernetesPodEvidence) SetContainers(value []ContainerEvidenceable)() {
    err := m.GetBackingStore().Set("containers", value)
    if err != nil {
        panic(err)
    }
}
// SetController sets the controller property value. The pod controller.
func (m *KubernetesPodEvidence) SetController(value KubernetesControllerEvidenceable)() {
    err := m.GetBackingStore().Set("controller", value)
    if err != nil {
        panic(err)
    }
}
// SetEphemeralContainers sets the ephemeralContainers property value. The list of pod ephemeral containers.
func (m *KubernetesPodEvidence) SetEphemeralContainers(value []ContainerEvidenceable)() {
    err := m.GetBackingStore().Set("ephemeralContainers", value)
    if err != nil {
        panic(err)
    }
}
// SetInitContainers sets the initContainers property value. The list of pod init containers.
func (m *KubernetesPodEvidence) SetInitContainers(value []ContainerEvidenceable)() {
    err := m.GetBackingStore().Set("initContainers", value)
    if err != nil {
        panic(err)
    }
}
// SetLabels sets the labels property value. The pod labels.
func (m *KubernetesPodEvidence) SetLabels(value Dictionaryable)() {
    err := m.GetBackingStore().Set("labels", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The pod name.
func (m *KubernetesPodEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNamespace sets the namespace property value. The pod namespace.
func (m *KubernetesPodEvidence) SetNamespace(value KubernetesNamespaceEvidenceable)() {
    err := m.GetBackingStore().Set("namespace", value)
    if err != nil {
        panic(err)
    }
}
// SetPodIp sets the podIp property value. The pod IP.
func (m *KubernetesPodEvidence) SetPodIp(value IpEvidenceable)() {
    err := m.GetBackingStore().Set("podIp", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceAccount sets the serviceAccount property value. The pod service account.
func (m *KubernetesPodEvidence) SetServiceAccount(value KubernetesServiceAccountEvidenceable)() {
    err := m.GetBackingStore().Set("serviceAccount", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesPodEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContainers()([]ContainerEvidenceable)
    GetController()(KubernetesControllerEvidenceable)
    GetEphemeralContainers()([]ContainerEvidenceable)
    GetInitContainers()([]ContainerEvidenceable)
    GetLabels()(Dictionaryable)
    GetName()(*string)
    GetNamespace()(KubernetesNamespaceEvidenceable)
    GetPodIp()(IpEvidenceable)
    GetServiceAccount()(KubernetesServiceAccountEvidenceable)
    SetContainers(value []ContainerEvidenceable)()
    SetController(value KubernetesControllerEvidenceable)()
    SetEphemeralContainers(value []ContainerEvidenceable)()
    SetInitContainers(value []ContainerEvidenceable)()
    SetLabels(value Dictionaryable)()
    SetName(value *string)()
    SetNamespace(value KubernetesNamespaceEvidenceable)()
    SetPodIp(value IpEvidenceable)()
    SetServiceAccount(value KubernetesServiceAccountEvidenceable)()
}
