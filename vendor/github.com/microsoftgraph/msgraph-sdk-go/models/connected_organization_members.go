package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConnectedOrganizationMembers struct {
    SubjectSet
}
// NewConnectedOrganizationMembers instantiates a new ConnectedOrganizationMembers and sets the default values.
func NewConnectedOrganizationMembers()(*ConnectedOrganizationMembers) {
    m := &ConnectedOrganizationMembers{
        SubjectSet: *NewSubjectSet(),
    }
    odataTypeValue := "#microsoft.graph.connectedOrganizationMembers"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateConnectedOrganizationMembersFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConnectedOrganizationMembersFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConnectedOrganizationMembers(), nil
}
// GetConnectedOrganizationId gets the connectedOrganizationId property value. The ID of the connected organization in entitlement management.
// returns a *string when successful
func (m *ConnectedOrganizationMembers) GetConnectedOrganizationId()(*string) {
    val, err := m.GetBackingStore().Get("connectedOrganizationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescription gets the description property value. The name of the connected organization.
// returns a *string when successful
func (m *ConnectedOrganizationMembers) GetDescription()(*string) {
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
func (m *ConnectedOrganizationMembers) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SubjectSet.GetFieldDeserializers()
    res["connectedOrganizationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectedOrganizationId(val)
        }
        return nil
    }
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
    return res
}
// Serialize serializes information the current object
func (m *ConnectedOrganizationMembers) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SubjectSet.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("connectedOrganizationId", m.GetConnectedOrganizationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConnectedOrganizationId sets the connectedOrganizationId property value. The ID of the connected organization in entitlement management.
func (m *ConnectedOrganizationMembers) SetConnectedOrganizationId(value *string)() {
    err := m.GetBackingStore().Set("connectedOrganizationId", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. The name of the connected organization.
func (m *ConnectedOrganizationMembers) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
type ConnectedOrganizationMembersable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SubjectSetable
    GetConnectedOrganizationId()(*string)
    GetDescription()(*string)
    SetConnectedOrganizationId(value *string)()
    SetDescription(value *string)()
}
