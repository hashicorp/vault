package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type VirtualEventPresenterInfo struct {
    MeetingParticipantInfo
}
// NewVirtualEventPresenterInfo instantiates a new VirtualEventPresenterInfo and sets the default values.
func NewVirtualEventPresenterInfo()(*VirtualEventPresenterInfo) {
    m := &VirtualEventPresenterInfo{
        MeetingParticipantInfo: *NewMeetingParticipantInfo(),
    }
    odataTypeValue := "#microsoft.graph.virtualEventPresenterInfo"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateVirtualEventPresenterInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventPresenterInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventPresenterInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *VirtualEventPresenterInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MeetingParticipantInfo.GetFieldDeserializers()
    res["presenterDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVirtualEventPresenterDetailsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPresenterDetails(val.(VirtualEventPresenterDetailsable))
        }
        return nil
    }
    return res
}
// GetPresenterDetails gets the presenterDetails property value. The presenterDetails property
// returns a VirtualEventPresenterDetailsable when successful
func (m *VirtualEventPresenterInfo) GetPresenterDetails()(VirtualEventPresenterDetailsable) {
    val, err := m.GetBackingStore().Get("presenterDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VirtualEventPresenterDetailsable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventPresenterInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MeetingParticipantInfo.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("presenterDetails", m.GetPresenterDetails())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPresenterDetails sets the presenterDetails property value. The presenterDetails property
func (m *VirtualEventPresenterInfo) SetPresenterDetails(value VirtualEventPresenterDetailsable)() {
    err := m.GetBackingStore().Set("presenterDetails", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventPresenterInfoable interface {
    MeetingParticipantInfoable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPresenterDetails()(VirtualEventPresenterDetailsable)
    SetPresenterDetails(value VirtualEventPresenterDetailsable)()
}
