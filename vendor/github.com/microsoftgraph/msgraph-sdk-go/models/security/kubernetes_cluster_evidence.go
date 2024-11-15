package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesClusterEvidence struct {
    AlertEvidence
}
// NewKubernetesClusterEvidence instantiates a new KubernetesClusterEvidence and sets the default values.
func NewKubernetesClusterEvidence()(*KubernetesClusterEvidence) {
    m := &KubernetesClusterEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesClusterEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesClusterEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesClusterEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesClusterEvidence(), nil
}
// GetCloudResource gets the cloudResource property value. The cloud identifier of the cluster. Can be either an amazonResourceEvidence, azureResourceEvidence, or googleCloudResourceEvidence object.
// returns a AlertEvidenceable when successful
func (m *KubernetesClusterEvidence) GetCloudResource()(AlertEvidenceable) {
    val, err := m.GetBackingStore().Get("cloudResource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AlertEvidenceable)
    }
    return nil
}
// GetDistribution gets the distribution property value. The distribution type of the cluster.
// returns a *string when successful
func (m *KubernetesClusterEvidence) GetDistribution()(*string) {
    val, err := m.GetBackingStore().Get("distribution")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesClusterEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["cloudResource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAlertEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCloudResource(val.(AlertEvidenceable))
        }
        return nil
    }
    res["distribution"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDistribution(val)
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
    res["platform"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseKubernetesPlatform)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlatform(val.(*KubernetesPlatform))
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The cluster name.
// returns a *string when successful
func (m *KubernetesClusterEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPlatform gets the platform property value. The platform the cluster runs on. Possible values are: unknown, aks, eks, gke, arc, unknownFutureValue.
// returns a *KubernetesPlatform when successful
func (m *KubernetesClusterEvidence) GetPlatform()(*KubernetesPlatform) {
    val, err := m.GetBackingStore().Get("platform")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*KubernetesPlatform)
    }
    return nil
}
// GetVersion gets the version property value. The kubernetes version of the cluster.
// returns a *string when successful
func (m *KubernetesClusterEvidence) GetVersion()(*string) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesClusterEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("cloudResource", m.GetCloudResource())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("distribution", m.GetDistribution())
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
    if m.GetPlatform() != nil {
        cast := (*m.GetPlatform()).String()
        err = writer.WriteStringValue("platform", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCloudResource sets the cloudResource property value. The cloud identifier of the cluster. Can be either an amazonResourceEvidence, azureResourceEvidence, or googleCloudResourceEvidence object.
func (m *KubernetesClusterEvidence) SetCloudResource(value AlertEvidenceable)() {
    err := m.GetBackingStore().Set("cloudResource", value)
    if err != nil {
        panic(err)
    }
}
// SetDistribution sets the distribution property value. The distribution type of the cluster.
func (m *KubernetesClusterEvidence) SetDistribution(value *string)() {
    err := m.GetBackingStore().Set("distribution", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The cluster name.
func (m *KubernetesClusterEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatform sets the platform property value. The platform the cluster runs on. Possible values are: unknown, aks, eks, gke, arc, unknownFutureValue.
func (m *KubernetesClusterEvidence) SetPlatform(value *KubernetesPlatform)() {
    err := m.GetBackingStore().Set("platform", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. The kubernetes version of the cluster.
func (m *KubernetesClusterEvidence) SetVersion(value *string)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesClusterEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCloudResource()(AlertEvidenceable)
    GetDistribution()(*string)
    GetName()(*string)
    GetPlatform()(*KubernetesPlatform)
    GetVersion()(*string)
    SetCloudResource(value AlertEvidenceable)()
    SetDistribution(value *string)()
    SetName(value *string)()
    SetPlatform(value *KubernetesPlatform)()
    SetVersion(value *string)()
}
