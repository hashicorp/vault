package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type NetworkConnectionEvidence struct {
    AlertEvidence
}
// NewNetworkConnectionEvidence instantiates a new NetworkConnectionEvidence and sets the default values.
func NewNetworkConnectionEvidence()(*NetworkConnectionEvidence) {
    m := &NetworkConnectionEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.networkConnectionEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateNetworkConnectionEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateNetworkConnectionEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewNetworkConnectionEvidence(), nil
}
// GetDestinationAddress gets the destinationAddress property value. The destinationAddress property
// returns a IpEvidenceable when successful
func (m *NetworkConnectionEvidence) GetDestinationAddress()(IpEvidenceable) {
    val, err := m.GetBackingStore().Get("destinationAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IpEvidenceable)
    }
    return nil
}
// GetDestinationPort gets the destinationPort property value. The destinationPort property
// returns a *int32 when successful
func (m *NetworkConnectionEvidence) GetDestinationPort()(*int32) {
    val, err := m.GetBackingStore().Get("destinationPort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *NetworkConnectionEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["destinationAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationAddress(val.(IpEvidenceable))
        }
        return nil
    }
    res["destinationPort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDestinationPort(val)
        }
        return nil
    }
    res["protocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseProtocolType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtocol(val.(*ProtocolType))
        }
        return nil
    }
    res["sourceAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIpEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourceAddress(val.(IpEvidenceable))
        }
        return nil
    }
    res["sourcePort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSourcePort(val)
        }
        return nil
    }
    return res
}
// GetProtocol gets the protocol property value. The protocol property
// returns a *ProtocolType when successful
func (m *NetworkConnectionEvidence) GetProtocol()(*ProtocolType) {
    val, err := m.GetBackingStore().Get("protocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ProtocolType)
    }
    return nil
}
// GetSourceAddress gets the sourceAddress property value. The sourceAddress property
// returns a IpEvidenceable when successful
func (m *NetworkConnectionEvidence) GetSourceAddress()(IpEvidenceable) {
    val, err := m.GetBackingStore().Get("sourceAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IpEvidenceable)
    }
    return nil
}
// GetSourcePort gets the sourcePort property value. The sourcePort property
// returns a *int32 when successful
func (m *NetworkConnectionEvidence) GetSourcePort()(*int32) {
    val, err := m.GetBackingStore().Get("sourcePort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *NetworkConnectionEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("destinationAddress", m.GetDestinationAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("destinationPort", m.GetDestinationPort())
        if err != nil {
            return err
        }
    }
    if m.GetProtocol() != nil {
        cast := (*m.GetProtocol()).String()
        err = writer.WriteStringValue("protocol", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sourceAddress", m.GetSourceAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("sourcePort", m.GetSourcePort())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDestinationAddress sets the destinationAddress property value. The destinationAddress property
func (m *NetworkConnectionEvidence) SetDestinationAddress(value IpEvidenceable)() {
    err := m.GetBackingStore().Set("destinationAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetDestinationPort sets the destinationPort property value. The destinationPort property
func (m *NetworkConnectionEvidence) SetDestinationPort(value *int32)() {
    err := m.GetBackingStore().Set("destinationPort", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocol sets the protocol property value. The protocol property
func (m *NetworkConnectionEvidence) SetProtocol(value *ProtocolType)() {
    err := m.GetBackingStore().Set("protocol", value)
    if err != nil {
        panic(err)
    }
}
// SetSourceAddress sets the sourceAddress property value. The sourceAddress property
func (m *NetworkConnectionEvidence) SetSourceAddress(value IpEvidenceable)() {
    err := m.GetBackingStore().Set("sourceAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetSourcePort sets the sourcePort property value. The sourcePort property
func (m *NetworkConnectionEvidence) SetSourcePort(value *int32)() {
    err := m.GetBackingStore().Set("sourcePort", value)
    if err != nil {
        panic(err)
    }
}
type NetworkConnectionEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDestinationAddress()(IpEvidenceable)
    GetDestinationPort()(*int32)
    GetProtocol()(*ProtocolType)
    GetSourceAddress()(IpEvidenceable)
    GetSourcePort()(*int32)
    SetDestinationAddress(value IpEvidenceable)()
    SetDestinationPort(value *int32)()
    SetProtocol(value *ProtocolType)()
    SetSourceAddress(value IpEvidenceable)()
    SetSourcePort(value *int32)()
}
