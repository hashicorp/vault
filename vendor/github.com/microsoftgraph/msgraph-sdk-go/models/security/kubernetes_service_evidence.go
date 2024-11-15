package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type KubernetesServiceEvidence struct {
    AlertEvidence
}
// NewKubernetesServiceEvidence instantiates a new KubernetesServiceEvidence and sets the default values.
func NewKubernetesServiceEvidence()(*KubernetesServiceEvidence) {
    m := &KubernetesServiceEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.kubernetesServiceEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateKubernetesServiceEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateKubernetesServiceEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewKubernetesServiceEvidence(), nil
}
// GetClusterIP gets the clusterIP property value. The service cluster IP.
// returns a IpEvidenceable when successful
func (m *KubernetesServiceEvidence) GetClusterIP()(IpEvidenceable) {
    val, err := m.GetBackingStore().Get("clusterIP")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IpEvidenceable)
    }
    return nil
}
// GetExternalIPs gets the externalIPs property value. The service external IPs.
// returns a []IpEvidenceable when successful
func (m *KubernetesServiceEvidence) GetExternalIPs()([]IpEvidenceable) {
    val, err := m.GetBackingStore().Get("externalIPs")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IpEvidenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *KubernetesServiceEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["clusterIP"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClusterIP(val.(IpEvidenceable))
        }
        return nil
    }
    res["externalIPs"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IpEvidenceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IpEvidenceable)
                }
            }
            m.SetExternalIPs(res)
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
    res["selector"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDictionaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSelector(val.(Dictionaryable))
        }
        return nil
    }
    res["servicePorts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKubernetesServicePortFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KubernetesServicePortable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KubernetesServicePortable)
                }
            }
            m.SetServicePorts(res)
        }
        return nil
    }
    res["serviceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseKubernetesServiceType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceType(val.(*KubernetesServiceType))
        }
        return nil
    }
    return res
}
// GetLabels gets the labels property value. The service labels.
// returns a Dictionaryable when successful
func (m *KubernetesServiceEvidence) GetLabels()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("labels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetName gets the name property value. The service name.
// returns a *string when successful
func (m *KubernetesServiceEvidence) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNamespace gets the namespace property value. The service namespace.
// returns a KubernetesNamespaceEvidenceable when successful
func (m *KubernetesServiceEvidence) GetNamespace()(KubernetesNamespaceEvidenceable) {
    val, err := m.GetBackingStore().Get("namespace")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(KubernetesNamespaceEvidenceable)
    }
    return nil
}
// GetSelector gets the selector property value. The service selector.
// returns a Dictionaryable when successful
func (m *KubernetesServiceEvidence) GetSelector()(Dictionaryable) {
    val, err := m.GetBackingStore().Get("selector")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Dictionaryable)
    }
    return nil
}
// GetServicePorts gets the servicePorts property value. The list of service ports.
// returns a []KubernetesServicePortable when successful
func (m *KubernetesServiceEvidence) GetServicePorts()([]KubernetesServicePortable) {
    val, err := m.GetBackingStore().Get("servicePorts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KubernetesServicePortable)
    }
    return nil
}
// GetServiceType gets the serviceType property value. The serviceType property
// returns a *KubernetesServiceType when successful
func (m *KubernetesServiceEvidence) GetServiceType()(*KubernetesServiceType) {
    val, err := m.GetBackingStore().Get("serviceType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*KubernetesServiceType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *KubernetesServiceEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("clusterIP", m.GetClusterIP())
        if err != nil {
            return err
        }
    }
    if m.GetExternalIPs() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExternalIPs()))
        for i, v := range m.GetExternalIPs() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("externalIPs", cast)
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
        err = writer.WriteObjectValue("selector", m.GetSelector())
        if err != nil {
            return err
        }
    }
    if m.GetServicePorts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServicePorts()))
        for i, v := range m.GetServicePorts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("servicePorts", cast)
        if err != nil {
            return err
        }
    }
    if m.GetServiceType() != nil {
        cast := (*m.GetServiceType()).String()
        err = writer.WriteStringValue("serviceType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetClusterIP sets the clusterIP property value. The service cluster IP.
func (m *KubernetesServiceEvidence) SetClusterIP(value IpEvidenceable)() {
    err := m.GetBackingStore().Set("clusterIP", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalIPs sets the externalIPs property value. The service external IPs.
func (m *KubernetesServiceEvidence) SetExternalIPs(value []IpEvidenceable)() {
    err := m.GetBackingStore().Set("externalIPs", value)
    if err != nil {
        panic(err)
    }
}
// SetLabels sets the labels property value. The service labels.
func (m *KubernetesServiceEvidence) SetLabels(value Dictionaryable)() {
    err := m.GetBackingStore().Set("labels", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The service name.
func (m *KubernetesServiceEvidence) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNamespace sets the namespace property value. The service namespace.
func (m *KubernetesServiceEvidence) SetNamespace(value KubernetesNamespaceEvidenceable)() {
    err := m.GetBackingStore().Set("namespace", value)
    if err != nil {
        panic(err)
    }
}
// SetSelector sets the selector property value. The service selector.
func (m *KubernetesServiceEvidence) SetSelector(value Dictionaryable)() {
    err := m.GetBackingStore().Set("selector", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePorts sets the servicePorts property value. The list of service ports.
func (m *KubernetesServiceEvidence) SetServicePorts(value []KubernetesServicePortable)() {
    err := m.GetBackingStore().Set("servicePorts", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceType sets the serviceType property value. The serviceType property
func (m *KubernetesServiceEvidence) SetServiceType(value *KubernetesServiceType)() {
    err := m.GetBackingStore().Set("serviceType", value)
    if err != nil {
        panic(err)
    }
}
type KubernetesServiceEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetClusterIP()(IpEvidenceable)
    GetExternalIPs()([]IpEvidenceable)
    GetLabels()(Dictionaryable)
    GetName()(*string)
    GetNamespace()(KubernetesNamespaceEvidenceable)
    GetSelector()(Dictionaryable)
    GetServicePorts()([]KubernetesServicePortable)
    GetServiceType()(*KubernetesServiceType)
    SetClusterIP(value IpEvidenceable)()
    SetExternalIPs(value []IpEvidenceable)()
    SetLabels(value Dictionaryable)()
    SetName(value *string)()
    SetNamespace(value KubernetesNamespaceEvidenceable)()
    SetSelector(value Dictionaryable)()
    SetServicePorts(value []KubernetesServicePortable)()
    SetServiceType(value *KubernetesServiceType)()
}
