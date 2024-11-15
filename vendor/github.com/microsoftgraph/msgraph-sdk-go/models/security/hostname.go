package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Hostname struct {
    Host
}
// NewHostname instantiates a new Hostname and sets the default values.
func NewHostname()(*Hostname) {
    m := &Hostname{
        Host: *NewHost(),
    }
    odataTypeValue := "#microsoft.graph.security.hostname"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateHostnameFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateHostnameFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewHostname(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Hostname) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Host.GetFieldDeserializers()
    res["registrant"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistrant(val)
        }
        return nil
    }
    res["registrar"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistrar(val)
        }
        return nil
    }
    return res
}
// GetRegistrant gets the registrant property value. The company or individual who registered this hostname, from WHOIS data.
// returns a *string when successful
func (m *Hostname) GetRegistrant()(*string) {
    val, err := m.GetBackingStore().Get("registrant")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistrar gets the registrar property value. The registrar for this hostname, from WHOIS data.
// returns a *string when successful
func (m *Hostname) GetRegistrar()(*string) {
    val, err := m.GetBackingStore().Get("registrar")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Hostname) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Host.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("registrant", m.GetRegistrant())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("registrar", m.GetRegistrar())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRegistrant sets the registrant property value. The company or individual who registered this hostname, from WHOIS data.
func (m *Hostname) SetRegistrant(value *string)() {
    err := m.GetBackingStore().Set("registrant", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistrar sets the registrar property value. The registrar for this hostname, from WHOIS data.
func (m *Hostname) SetRegistrar(value *string)() {
    err := m.GetBackingStore().Set("registrar", value)
    if err != nil {
        panic(err)
    }
}
type Hostnameable interface {
    Hostable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRegistrant()(*string)
    GetRegistrar()(*string)
    SetRegistrant(value *string)()
    SetRegistrar(value *string)()
}
