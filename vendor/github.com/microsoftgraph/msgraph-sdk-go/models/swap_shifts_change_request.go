package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SwapShiftsChangeRequest struct {
    OfferShiftRequest
}
// NewSwapShiftsChangeRequest instantiates a new SwapShiftsChangeRequest and sets the default values.
func NewSwapShiftsChangeRequest()(*SwapShiftsChangeRequest) {
    m := &SwapShiftsChangeRequest{
        OfferShiftRequest: *NewOfferShiftRequest(),
    }
    odataTypeValue := "#microsoft.graph.swapShiftsChangeRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateSwapShiftsChangeRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSwapShiftsChangeRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSwapShiftsChangeRequest(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SwapShiftsChangeRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OfferShiftRequest.GetFieldDeserializers()
    res["recipientShiftId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipientShiftId(val)
        }
        return nil
    }
    return res
}
// GetRecipientShiftId gets the recipientShiftId property value. ShiftId for the recipient user with whom the request is to swap.
// returns a *string when successful
func (m *SwapShiftsChangeRequest) GetRecipientShiftId()(*string) {
    val, err := m.GetBackingStore().Get("recipientShiftId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SwapShiftsChangeRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OfferShiftRequest.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("recipientShiftId", m.GetRecipientShiftId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetRecipientShiftId sets the recipientShiftId property value. ShiftId for the recipient user with whom the request is to swap.
func (m *SwapShiftsChangeRequest) SetRecipientShiftId(value *string)() {
    err := m.GetBackingStore().Set("recipientShiftId", value)
    if err != nil {
        panic(err)
    }
}
type SwapShiftsChangeRequestable interface {
    OfferShiftRequestable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetRecipientShiftId()(*string)
    SetRecipientShiftId(value *string)()
}
