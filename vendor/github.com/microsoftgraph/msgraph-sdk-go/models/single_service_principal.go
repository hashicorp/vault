package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SingleServicePrincipal struct {
    SubjectSet
}
// NewSingleServicePrincipal instantiates a new SingleServicePrincipal and sets the default values.
func NewSingleServicePrincipal()(*SingleServicePrincipal) {
    m := &SingleServicePrincipal{
        SubjectSet: *NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.singleServicePrincipal"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSingleServicePrincipalFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSingleServicePrincipalFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSingleServicePrincipal(), nil
}
// GetDescription gets the description property value. Description of this service principal.
// returns a *string when successful
func (m *SingleServicePrincipal) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
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
func (m *SingleServicePrincipal) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["servicePrincipalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePrincipalId(val)
        }
        return nil
    }
    return res
}
// GetServicePrincipalId gets the servicePrincipalId property value. ID of the servicePrincipal.
// returns a *string when successful
func (m *SingleServicePrincipal) GetServicePrincipalId()(*string) {
    val, err := m.GetBackingStore().Get("servicePrincipalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SingleServicePrincipal) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePrincipalId", m.GetServicePrincipalId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Description of this service principal.
func (m *SingleServicePrincipal) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalId sets the servicePrincipalId property value. ID of the servicePrincipal.
func (m *SingleServicePrincipal) SetServicePrincipalId(value *string)() {
    err := m.GetBackingStore().Set("servicePrincipalId", value)
    if err != nil {
        panic(err)
    }
}
type SingleServicePrincipalable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SubjectSetable
    GetDescription()(*string)
    GetServicePrincipalId()(*string)
    SetDescription(value *string)()
    SetServicePrincipalId(value *string)()
}
