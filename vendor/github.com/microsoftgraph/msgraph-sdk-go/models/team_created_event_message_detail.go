package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamCreatedEventMessageDetail struct {
    EventMessageDetail
}
// NewTeamCreatedEventMessageDetail instantiates a new TeamCreatedEventMessageDetail and sets the default values.
func NewTeamCreatedEventMessageDetail()(*TeamCreatedEventMessageDetail) {
    m := &TeamCreatedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.teamCreatedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTeamCreatedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamCreatedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamCreatedEventMessageDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamCreatedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["teamDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamDescription(val)
        }
        return nil
    }
    res["teamDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamDisplayName(val)
        }
        return nil
    }
    res["teamId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTeamId(val)
        }
        return nil
    }
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *TeamCreatedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetTeamDescription gets the teamDescription property value. Description for the team.
// returns a *string when successful
func (m *TeamCreatedEventMessageDetail) GetTeamDescription()(*string) {
    val, err := m.GetBackingStore().Get("teamDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTeamDisplayName gets the teamDisplayName property value. Display name of the team.
// returns a *string when successful
func (m *TeamCreatedEventMessageDetail) GetTeamDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("teamDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTeamId gets the teamId property value. Unique identifier of the team.
// returns a *string when successful
func (m *TeamCreatedEventMessageDetail) GetTeamId()(*string) {
    val, err := m.GetBackingStore().Get("teamId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamCreatedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("teamDescription", m.GetTeamDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamDisplayName", m.GetTeamDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("teamId", m.GetTeamId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *TeamCreatedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamDescription sets the teamDescription property value. Description for the team.
func (m *TeamCreatedEventMessageDetail) SetTeamDescription(value *string)() {
    err := m.GetBackingStore().Set("teamDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamDisplayName sets the teamDisplayName property value. Display name of the team.
func (m *TeamCreatedEventMessageDetail) SetTeamDisplayName(value *string)() {
    err := m.GetBackingStore().Set("teamDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetTeamId sets the teamId property value. Unique identifier of the team.
func (m *TeamCreatedEventMessageDetail) SetTeamId(value *string)() {
    err := m.GetBackingStore().Set("teamId", value)
    if err != nil {
        panic(err)
    }
}
type TeamCreatedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInitiator()(IdentitySetable)
    GetTeamDescription()(*string)
    GetTeamDisplayName()(*string)
    GetTeamId()(*string)
    SetInitiator(value IdentitySetable)()
    SetTeamDescription(value *string)()
    SetTeamDisplayName(value *string)()
    SetTeamId(value *string)()
}
