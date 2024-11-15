package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DomainDnsSrvRecord struct {
    DomainDnsRecord
}
// NewDomainDnsSrvRecord instantiates a new DomainDnsSrvRecord and sets the default values.
func NewDomainDnsSrvRecord()(*DomainDnsSrvRecord) {
    m := &DomainDnsSrvRecord{
        DomainDnsRecord: *NewDomainDnsRecord(),
    }
    return m
}
// CreateDomainDnsSrvRecordFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDomainDnsSrvRecordFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDomainDnsSrvRecord(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DomainDnsSrvRecord) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DomainDnsRecord.GetFieldDeserializers()
    res["nameTarget"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNameTarget(val)
        }
        return nil
    }
    res["port"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPort(val)
        }
        return nil
    }
    res["priority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPriority(val)
        }
        return nil
    }
    res["protocol"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtocol(val)
        }
        return nil
    }
    res["service"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetService(val)
        }
        return nil
    }
    res["weight"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWeight(val)
        }
        return nil
    }
    return res
}
// GetNameTarget gets the nameTarget property value. Value to use when configuring the Target property of the SRV record at the DNS host.
// returns a *string when successful
func (m *DomainDnsSrvRecord) GetNameTarget()(*string) {
    val, err := m.GetBackingStore().Get("nameTarget")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPort gets the port property value. Value to use when configuring the port property of the SRV record at the DNS host.
// returns a *int32 when successful
func (m *DomainDnsSrvRecord) GetPort()(*int32) {
    val, err := m.GetBackingStore().Get("port")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPriority gets the priority property value. Value to use when configuring the priority property of the SRV record at the DNS host.
// returns a *int32 when successful
func (m *DomainDnsSrvRecord) GetPriority()(*int32) {
    val, err := m.GetBackingStore().Get("priority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetProtocol gets the protocol property value. Value to use when configuring the protocol property of the SRV record at the DNS host.
// returns a *string when successful
func (m *DomainDnsSrvRecord) GetProtocol()(*string) {
    val, err := m.GetBackingStore().Get("protocol")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetService gets the service property value. Value to use when configuring the service property of the SRV record at the DNS host.
// returns a *string when successful
func (m *DomainDnsSrvRecord) GetService()(*string) {
    val, err := m.GetBackingStore().Get("service")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWeight gets the weight property value. Value to use when configuring the weight property of the SRV record at the DNS host.
// returns a *int32 when successful
func (m *DomainDnsSrvRecord) GetWeight()(*int32) {
    val, err := m.GetBackingStore().Get("weight")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DomainDnsSrvRecord) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DomainDnsRecord.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("nameTarget", m.GetNameTarget())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("port", m.GetPort())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("priority", m.GetPriority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("protocol", m.GetProtocol())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("service", m.GetService())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("weight", m.GetWeight())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetNameTarget sets the nameTarget property value. Value to use when configuring the Target property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetNameTarget(value *string)() {
    err := m.GetBackingStore().Set("nameTarget", value)
    if err != nil {
        panic(err)
    }
}
// SetPort sets the port property value. Value to use when configuring the port property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetPort(value *int32)() {
    err := m.GetBackingStore().Set("port", value)
    if err != nil {
        panic(err)
    }
}
// SetPriority sets the priority property value. Value to use when configuring the priority property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetPriority(value *int32)() {
    err := m.GetBackingStore().Set("priority", value)
    if err != nil {
        panic(err)
    }
}
// SetProtocol sets the protocol property value. Value to use when configuring the protocol property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetProtocol(value *string)() {
    err := m.GetBackingStore().Set("protocol", value)
    if err != nil {
        panic(err)
    }
}
// SetService sets the service property value. Value to use when configuring the service property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetService(value *string)() {
    err := m.GetBackingStore().Set("service", value)
    if err != nil {
        panic(err)
    }
}
// SetWeight sets the weight property value. Value to use when configuring the weight property of the SRV record at the DNS host.
func (m *DomainDnsSrvRecord) SetWeight(value *int32)() {
    err := m.GetBackingStore().Set("weight", value)
    if err != nil {
        panic(err)
    }
}
type DomainDnsSrvRecordable interface {
    DomainDnsRecordable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetNameTarget()(*string)
    GetPort()(*int32)
    GetPriority()(*int32)
    GetProtocol()(*string)
    GetService()(*string)
    GetWeight()(*int32)
    SetNameTarget(value *string)()
    SetPort(value *int32)()
    SetPriority(value *int32)()
    SetProtocol(value *string)()
    SetService(value *string)()
    SetWeight(value *int32)()
}
