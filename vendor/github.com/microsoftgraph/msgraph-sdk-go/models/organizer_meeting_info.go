package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OrganizerMeetingInfo struct {
    MeetingInfo
}
// NewOrganizerMeetingInfo instantiates a new OrganizerMeetingInfo and sets the default values.
func NewOrganizerMeetingInfo()(*OrganizerMeetingInfo) {
    m := &OrganizerMeetingInfo{
        MeetingInfo: *NewMeetingInfo(),
    }
    odataTypeValue := "#microsoft.graph.organizerMeetingInfo"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOrganizerMeetingInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOrganizerMeetingInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOrganizerMeetingInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OrganizerMeetingInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MeetingInfo.GetFieldDeserializers()
    res["organizer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizer(val.(IdentitySetable))
        }
        return nil
    }
    return res
}
// GetOrganizer gets the organizer property value. The organizer property
// returns a IdentitySetable when successful
func (m *OrganizerMeetingInfo) GetOrganizer()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("organizer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *OrganizerMeetingInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MeetingInfo.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("organizer", m.GetOrganizer())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOrganizer sets the organizer property value. The organizer property
func (m *OrganizerMeetingInfo) SetOrganizer(value IdentitySetable)() {
    err := m.GetBackingStore().Set("organizer", value)
    if err != nil {
        panic(err)
    }
}
type OrganizerMeetingInfoable interface {
    MeetingInfoable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetOrganizer()(IdentitySetable)
    SetOrganizer(value IdentitySetable)()
}
