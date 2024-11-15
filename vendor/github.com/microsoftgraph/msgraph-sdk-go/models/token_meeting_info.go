package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TokenMeetingInfo struct {
    MeetingInfo
}
// NewTokenMeetingInfo instantiates a new TokenMeetingInfo and sets the default values.
func NewTokenMeetingInfo()(*TokenMeetingInfo) {
    m := &TokenMeetingInfo{
        MeetingInfo: *NewMeetingInfo(),
    }
    odataTypeValue := "#microsoft.graph.tokenMeetingInfo"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTokenMeetingInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTokenMeetingInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTokenMeetingInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TokenMeetingInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MeetingInfo.GetFieldDeserializers()
    res["token"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetToken(val)
        }
        return nil
    }
    return res
}
// GetToken gets the token property value. The token used to join the call.
// returns a *string when successful
func (m *TokenMeetingInfo) GetToken()(*string) {
    val, err := m.GetBackingStore().Get("token")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TokenMeetingInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MeetingInfo.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("token", m.GetToken())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetToken sets the token property value. The token used to join the call.
func (m *TokenMeetingInfo) SetToken(value *string)() {
    err := m.GetBackingStore().Set("token", value)
    if err != nil {
        panic(err)
    }
}
type TokenMeetingInfoable interface {
    MeetingInfoable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetToken()(*string)
    SetToken(value *string)()
}
