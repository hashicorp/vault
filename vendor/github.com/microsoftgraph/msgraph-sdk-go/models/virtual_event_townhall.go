package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventTownhall struct {
    VirtualEvent
}
// NewVirtualEventTownhall instantiates a new VirtualEventTownhall and sets the default values.
func NewVirtualEventTownhall()(*VirtualEventTownhall) {
    m := &VirtualEventTownhall{
        VirtualEvent: *NewVirtualEvent(),
    }
    odataTypeValue := "#microsoft.graph.virtualEventTownhall"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateVirtualEventTownhallFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventTownhallFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventTownhall(), nil
}
// GetAudience gets the audience property value. The audience to whom the town hall is visible. Possible values are: everyone, organization, and unknownFutureValue.
// returns a *MeetingAudience when successful
func (m *VirtualEventTownhall) GetAudience()(*MeetingAudience) {
    val, err := m.GetBackingStore().Get("audience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MeetingAudience)
    }
    return nil
}
// GetCoOrganizers gets the coOrganizers property value. Identity information of the coorganizers of the town hall.
// returns a []CommunicationsUserIdentityable when successful
func (m *VirtualEventTownhall) GetCoOrganizers()([]CommunicationsUserIdentityable) {
    val, err := m.GetBackingStore().Get("coOrganizers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CommunicationsUserIdentityable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventTownhall) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.VirtualEvent.GetFieldDeserializers()
    res["audience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMeetingAudience)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudience(val.(*MeetingAudience))
        }
        return nil
    }
    res["coOrganizers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCommunicationsUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CommunicationsUserIdentityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CommunicationsUserIdentityable)
                }
            }
            m.SetCoOrganizers(res)
        }
        return nil
    }
    res["invitedAttendees"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Identityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Identityable)
                }
            }
            m.SetInvitedAttendees(res)
        }
        return nil
    }
    res["isInviteOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInviteOnly(val)
        }
        return nil
    }
    return res
}
// GetInvitedAttendees gets the invitedAttendees property value. The attendees invited to the town hall. The supported identities are: communicationsUserIdentity and communicationsGuestIdentity.
// returns a []Identityable when successful
func (m *VirtualEventTownhall) GetInvitedAttendees()([]Identityable) {
    val, err := m.GetBackingStore().Get("invitedAttendees")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Identityable)
    }
    return nil
}
// GetIsInviteOnly gets the isInviteOnly property value. Indicates whether the town hall is only open to invited people and groups within your organization. The isInviteOnly property can only be true if the value of the audience property is set to organization.
// returns a *bool when successful
func (m *VirtualEventTownhall) GetIsInviteOnly()(*bool) {
    val, err := m.GetBackingStore().Get("isInviteOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventTownhall) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.VirtualEvent.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAudience() != nil {
        cast := (*m.GetAudience()).String()
        err = writer.WriteStringValue("audience", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetCoOrganizers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCoOrganizers()))
        for i, v := range m.GetCoOrganizers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("coOrganizers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetInvitedAttendees() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInvitedAttendees()))
        for i, v := range m.GetInvitedAttendees() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("invitedAttendees", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isInviteOnly", m.GetIsInviteOnly())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAudience sets the audience property value. The audience to whom the town hall is visible. Possible values are: everyone, organization, and unknownFutureValue.
func (m *VirtualEventTownhall) SetAudience(value *MeetingAudience)() {
    err := m.GetBackingStore().Set("audience", value)
    if err != nil {
        panic(err)
    }
}
// SetCoOrganizers sets the coOrganizers property value. Identity information of the coorganizers of the town hall.
func (m *VirtualEventTownhall) SetCoOrganizers(value []CommunicationsUserIdentityable)() {
    err := m.GetBackingStore().Set("coOrganizers", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedAttendees sets the invitedAttendees property value. The attendees invited to the town hall. The supported identities are: communicationsUserIdentity and communicationsGuestIdentity.
func (m *VirtualEventTownhall) SetInvitedAttendees(value []Identityable)() {
    err := m.GetBackingStore().Set("invitedAttendees", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInviteOnly sets the isInviteOnly property value. Indicates whether the town hall is only open to invited people and groups within your organization. The isInviteOnly property can only be true if the value of the audience property is set to organization.
func (m *VirtualEventTownhall) SetIsInviteOnly(value *bool)() {
    err := m.GetBackingStore().Set("isInviteOnly", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventTownhallable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    VirtualEventable
    GetAudience()(*MeetingAudience)
    GetCoOrganizers()([]CommunicationsUserIdentityable)
    GetInvitedAttendees()([]Identityable)
    GetIsInviteOnly()(*bool)
    SetAudience(value *MeetingAudience)()
    SetCoOrganizers(value []CommunicationsUserIdentityable)()
    SetInvitedAttendees(value []Identityable)()
    SetIsInviteOnly(value *bool)()
}
