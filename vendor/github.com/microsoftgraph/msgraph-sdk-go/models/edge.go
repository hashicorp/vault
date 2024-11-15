package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Edge struct {
    Entity
}
// NewEdge instantiates a new Edge and sets the default values.
func NewEdge()(*Edge) {
    m := &Edge{
        Entity: *NewEntity(),
    }
    return m
}
// CreateEdgeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEdgeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEdge(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Edge) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["internetExplorerMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInternetExplorerModeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternetExplorerMode(val.(InternetExplorerModeable))
        }
        return nil
    }
    return res
}
// GetInternetExplorerMode gets the internetExplorerMode property value. A container for Internet Explorer mode resources.
// returns a InternetExplorerModeable when successful
func (m *Edge) GetInternetExplorerMode()(InternetExplorerModeable) {
    val, err := m.GetBackingStore().Get("internetExplorerMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InternetExplorerModeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Edge) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("internetExplorerMode", m.GetInternetExplorerMode())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInternetExplorerMode sets the internetExplorerMode property value. A container for Internet Explorer mode resources.
func (m *Edge) SetInternetExplorerMode(value InternetExplorerModeable)() {
    err := m.GetBackingStore().Set("internetExplorerMode", value)
    if err != nil {
        panic(err)
    }
}
type Edgeable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInternetExplorerMode()(InternetExplorerModeable)
    SetInternetExplorerMode(value InternetExplorerModeable)()
}
