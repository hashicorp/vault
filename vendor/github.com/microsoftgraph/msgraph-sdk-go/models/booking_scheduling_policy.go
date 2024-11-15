package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// BookingSchedulingPolicy this type represents the set of policies that dictate how bookings can be created in a Booking Calendar.
type BookingSchedulingPolicy struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewBookingSchedulingPolicy instantiates a new BookingSchedulingPolicy and sets the default values.
func NewBookingSchedulingPolicy()(*BookingSchedulingPolicy) {
    m := &BookingSchedulingPolicy{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateBookingSchedulingPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBookingSchedulingPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingSchedulingPolicy(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *BookingSchedulingPolicy) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAllowStaffSelection gets the allowStaffSelection property value. True to allow customers to choose a specific person for the booking.
// returns a *bool when successful
func (m *BookingSchedulingPolicy) GetAllowStaffSelection()(*bool) {
    val, err := m.GetBackingStore().Get("allowStaffSelection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *BookingSchedulingPolicy) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCustomAvailabilities gets the customAvailabilities property value. Custom availability of the service in a given time frame.
// returns a []BookingsAvailabilityWindowable when successful
func (m *BookingSchedulingPolicy) GetCustomAvailabilities()([]BookingsAvailabilityWindowable) {
    val, err := m.GetBackingStore().Get("customAvailabilities")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BookingsAvailabilityWindowable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BookingSchedulingPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowStaffSelection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowStaffSelection(val)
        }
        return nil
    }
    res["customAvailabilities"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBookingsAvailabilityWindowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BookingsAvailabilityWindowable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BookingsAvailabilityWindowable)
                }
            }
            m.SetCustomAvailabilities(res)
        }
        return nil
    }
    res["generalAvailability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBookingsAvailabilityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGeneralAvailability(val.(BookingsAvailabilityable))
        }
        return nil
    }
    res["isMeetingInviteToCustomersEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMeetingInviteToCustomersEnabled(val)
        }
        return nil
    }
    res["maximumAdvance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaximumAdvance(val)
        }
        return nil
    }
    res["minimumLeadTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumLeadTime(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["sendConfirmationsToOwner"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSendConfirmationsToOwner(val)
        }
        return nil
    }
    res["timeSlotInterval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeSlotInterval(val)
        }
        return nil
    }
    return res
}
// GetGeneralAvailability gets the generalAvailability property value. General availability of the service defined by the scheduling policy.
// returns a BookingsAvailabilityable when successful
func (m *BookingSchedulingPolicy) GetGeneralAvailability()(BookingsAvailabilityable) {
    val, err := m.GetBackingStore().Get("generalAvailability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BookingsAvailabilityable)
    }
    return nil
}
// GetIsMeetingInviteToCustomersEnabled gets the isMeetingInviteToCustomersEnabled property value. Indicates whether the meeting invite is sent to the customers. The default value is false.
// returns a *bool when successful
func (m *BookingSchedulingPolicy) GetIsMeetingInviteToCustomersEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isMeetingInviteToCustomersEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaximumAdvance gets the maximumAdvance property value. Maximum number of days in advance that a booking can be made. It follows the ISO 8601 format.
// returns a *ISODuration when successful
func (m *BookingSchedulingPolicy) GetMaximumAdvance()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("maximumAdvance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetMinimumLeadTime gets the minimumLeadTime property value. The minimum amount of time before which bookings and cancellations must be made. It follows the ISO 8601 format.
// returns a *ISODuration when successful
func (m *BookingSchedulingPolicy) GetMinimumLeadTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("minimumLeadTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *BookingSchedulingPolicy) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSendConfirmationsToOwner gets the sendConfirmationsToOwner property value. True to notify the business via email when a booking is created or changed. Use the email address specified in the email property of the bookingBusiness entity for the business.
// returns a *bool when successful
func (m *BookingSchedulingPolicy) GetSendConfirmationsToOwner()(*bool) {
    val, err := m.GetBackingStore().Get("sendConfirmationsToOwner")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTimeSlotInterval gets the timeSlotInterval property value. Duration of each time slot, denoted in ISO 8601 format.
// returns a *ISODuration when successful
func (m *BookingSchedulingPolicy) GetTimeSlotInterval()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("timeSlotInterval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BookingSchedulingPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowStaffSelection", m.GetAllowStaffSelection())
        if err != nil {
            return err
        }
    }
    if m.GetCustomAvailabilities() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomAvailabilities()))
        for i, v := range m.GetCustomAvailabilities() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("customAvailabilities", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("generalAvailability", m.GetGeneralAvailability())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isMeetingInviteToCustomersEnabled", m.GetIsMeetingInviteToCustomersEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("maximumAdvance", m.GetMaximumAdvance())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("minimumLeadTime", m.GetMinimumLeadTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("sendConfirmationsToOwner", m.GetSendConfirmationsToOwner())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteISODurationValue("timeSlotInterval", m.GetTimeSlotInterval())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *BookingSchedulingPolicy) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowStaffSelection sets the allowStaffSelection property value. True to allow customers to choose a specific person for the booking.
func (m *BookingSchedulingPolicy) SetAllowStaffSelection(value *bool)() {
    err := m.GetBackingStore().Set("allowStaffSelection", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *BookingSchedulingPolicy) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCustomAvailabilities sets the customAvailabilities property value. Custom availability of the service in a given time frame.
func (m *BookingSchedulingPolicy) SetCustomAvailabilities(value []BookingsAvailabilityWindowable)() {
    err := m.GetBackingStore().Set("customAvailabilities", value)
    if err != nil {
        panic(err)
    }
}
// SetGeneralAvailability sets the generalAvailability property value. General availability of the service defined by the scheduling policy.
func (m *BookingSchedulingPolicy) SetGeneralAvailability(value BookingsAvailabilityable)() {
    err := m.GetBackingStore().Set("generalAvailability", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMeetingInviteToCustomersEnabled sets the isMeetingInviteToCustomersEnabled property value. Indicates whether the meeting invite is sent to the customers. The default value is false.
func (m *BookingSchedulingPolicy) SetIsMeetingInviteToCustomersEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isMeetingInviteToCustomersEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetMaximumAdvance sets the maximumAdvance property value. Maximum number of days in advance that a booking can be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) SetMaximumAdvance(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("maximumAdvance", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumLeadTime sets the minimumLeadTime property value. The minimum amount of time before which bookings and cancellations must be made. It follows the ISO 8601 format.
func (m *BookingSchedulingPolicy) SetMinimumLeadTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("minimumLeadTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *BookingSchedulingPolicy) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetSendConfirmationsToOwner sets the sendConfirmationsToOwner property value. True to notify the business via email when a booking is created or changed. Use the email address specified in the email property of the bookingBusiness entity for the business.
func (m *BookingSchedulingPolicy) SetSendConfirmationsToOwner(value *bool)() {
    err := m.GetBackingStore().Set("sendConfirmationsToOwner", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeSlotInterval sets the timeSlotInterval property value. Duration of each time slot, denoted in ISO 8601 format.
func (m *BookingSchedulingPolicy) SetTimeSlotInterval(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("timeSlotInterval", value)
    if err != nil {
        panic(err)
    }
}
type BookingSchedulingPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowStaffSelection()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCustomAvailabilities()([]BookingsAvailabilityWindowable)
    GetGeneralAvailability()(BookingsAvailabilityable)
    GetIsMeetingInviteToCustomersEnabled()(*bool)
    GetMaximumAdvance()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetMinimumLeadTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetOdataType()(*string)
    GetSendConfirmationsToOwner()(*bool)
    GetTimeSlotInterval()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    SetAllowStaffSelection(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCustomAvailabilities(value []BookingsAvailabilityWindowable)()
    SetGeneralAvailability(value BookingsAvailabilityable)()
    SetIsMeetingInviteToCustomersEnabled(value *bool)()
    SetMaximumAdvance(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetMinimumLeadTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetOdataType(value *string)()
    SetSendConfirmationsToOwner(value *bool)()
    SetTimeSlotInterval(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
}
