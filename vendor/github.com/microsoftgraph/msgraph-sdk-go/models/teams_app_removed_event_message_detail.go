package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsAppRemovedEventMessageDetail struct {
    EventMessageDetail
}
// NewTeamsAppRemovedEventMessageDetail instantiates a new TeamsAppRemovedEventMessageDetail and sets the default values.
func NewTeamsAppRemovedEventMessageDetail()(*TeamsAppRemovedEventMessageDetail) {
    m := &TeamsAppRemovedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.teamsAppRemovedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTeamsAppRemovedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsAppRemovedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsAppRemovedEventMessageDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsAppRemovedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["initiator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiator(val.(IdentitySetable))
        }
        return nil
    }
    res["teamsAppDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamsAppDisplayName(val)
        }
        return nil
    }
    res["teamsAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamsAppId(val)
        }
        return nil
    }
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *TeamsAppRemovedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetTeamsAppDisplayName gets the teamsAppDisplayName property value. Display name of the teamsApp.
// returns a *string when successful
func (m *TeamsAppRemovedEventMessageDetail) GetTeamsAppDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("teamsAppDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTeamsAppId gets the teamsAppId property value. Unique identifier of the teamsApp.
// returns a *string when successful
func (m *TeamsAppRemovedEventMessageDetail) GetTeamsAppId()(*string) {
    val, err := m.GetBackingStore().Get("teamsAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamsAppRemovedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamsAppDisplayName", m.GetTeamsAppDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamsAppId", m.GetTeamsAppId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *TeamsAppRemovedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamsAppDisplayName sets the teamsAppDisplayName property value. Display name of the teamsApp.
func (m *TeamsAppRemovedEventMessageDetail) SetTeamsAppDisplayName(value *string)() {
    err := m.GetBackingStore().Set("teamsAppDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamsAppId sets the teamsAppId property value. Unique identifier of the teamsApp.
func (m *TeamsAppRemovedEventMessageDetail) SetTeamsAppId(value *string)() {
    err := m.GetBackingStore().Set("teamsAppId", value)
    if err != nil {
        panic(err)
    }
}
type TeamsAppRemovedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInitiator()(IdentitySetable)
    GetTeamsAppDisplayName()(*string)
    GetTeamsAppId()(*string)
    SetInitiator(value IdentitySetable)()
    SetTeamsAppDisplayName(value *string)()
    SetTeamsAppId(value *string)()
}
