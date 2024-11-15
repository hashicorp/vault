package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MeetingTimeSuggestion struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMeetingTimeSuggestion instantiates a new MeetingTimeSuggestion and sets the default values.
func NewMeetingTimeSuggestion()(*MeetingTimeSuggestion) {
    m := &MeetingTimeSuggestion{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMeetingTimeSuggestionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMeetingTimeSuggestionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingTimeSuggestion(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MeetingTimeSuggestion) GetAdditionalData()(map[string]any) {
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
// GetAttendeeAvailability gets the attendeeAvailability property value. An array that shows the availability status of each attendee for this meeting suggestion.
// returns a []AttendeeAvailabilityable when successful
func (m *MeetingTimeSuggestion) GetAttendeeAvailability()([]AttendeeAvailabilityable) {
    val, err := m.GetBackingStore().Get("attendeeAvailability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AttendeeAvailabilityable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MeetingTimeSuggestion) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetConfidence gets the confidence property value. A percentage that represents the likelhood of all the attendees attending.
// returns a *float64 when successful
func (m *MeetingTimeSuggestion) GetConfidence()(*float64) {
    val, err := m.GetBackingStore().Get("confidence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MeetingTimeSuggestion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["attendeeAvailability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttendeeAvailabilityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AttendeeAvailabilityable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AttendeeAvailabilityable)
                }
            }
            m.SetAttendeeAvailability(res)
        }
        return nil
    }
    res["confidence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfidence(val)
        }
        return nil
    }
    res["locations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Locationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Locationable)
                }
            }
            m.SetLocations(res)
        }
        return nil
    }
    res["meetingTimeSlot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTimeSlotFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingTimeSlot(val.(TimeSlotable))
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
    res["order"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrder(val)
        }
        return nil
    }
    res["organizerAvailability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFreeBusyStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizerAvailability(val.(*FreeBusyStatus))
        }
        return nil
    }
    res["suggestionReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuggestionReason(val)
        }
        return nil
    }
    return res
}
// GetLocations gets the locations property value. An array that specifies the name and geographic location of each meeting location for this meeting suggestion.
// returns a []Locationable when successful
func (m *MeetingTimeSuggestion) GetLocations()([]Locationable) {
    val, err := m.GetBackingStore().Get("locations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Locationable)
    }
    return nil
}
// GetMeetingTimeSlot gets the meetingTimeSlot property value. A time period suggested for the meeting.
// returns a TimeSlotable when successful
func (m *MeetingTimeSuggestion) GetMeetingTimeSlot()(TimeSlotable) {
    val, err := m.GetBackingStore().Get("meetingTimeSlot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TimeSlotable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MeetingTimeSuggestion) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrder gets the order property value. Order of meeting time suggestions sorted by their computed confidence value from high to low, then by chronology if there are suggestions with the same confidence.
// returns a *int32 when successful
func (m *MeetingTimeSuggestion) GetOrder()(*int32) {
    val, err := m.GetBackingStore().Get("order")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOrganizerAvailability gets the organizerAvailability property value. Availability of the meeting organizer for this meeting suggestion. The possible values are: free, tentative, busy, oof, workingElsewhere, unknown.
// returns a *FreeBusyStatus when successful
func (m *MeetingTimeSuggestion) GetOrganizerAvailability()(*FreeBusyStatus) {
    val, err := m.GetBackingStore().Get("organizerAvailability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FreeBusyStatus)
    }
    return nil
}
// GetSuggestionReason gets the suggestionReason property value. Reason for suggesting the meeting time.
// returns a *string when successful
func (m *MeetingTimeSuggestion) GetSuggestionReason()(*string) {
    val, err := m.GetBackingStore().Get("suggestionReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MeetingTimeSuggestion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAttendeeAvailability() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttendeeAvailability()))
        for i, v := range m.GetAttendeeAvailability() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("attendeeAvailability", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("confidence", m.GetConfidence())
        if err != nil {
            return err
        }
    }
    if m.GetLocations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetLocations()))
        for i, v := range m.GetLocations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("locations", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("meetingTimeSlot", m.GetMeetingTimeSlot())
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
        err := writer.WriteInt32Value("order", m.GetOrder())
        if err != nil {
            return err
        }
    }
    if m.GetOrganizerAvailability() != nil {
        cast := (*m.GetOrganizerAvailability()).String()
        err := writer.WriteStringValue("organizerAvailability", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("suggestionReason", m.GetSuggestionReason())
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
func (m *MeetingTimeSuggestion) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAttendeeAvailability sets the attendeeAvailability property value. An array that shows the availability status of each attendee for this meeting suggestion.
func (m *MeetingTimeSuggestion) SetAttendeeAvailability(value []AttendeeAvailabilityable)() {
    err := m.GetBackingStore().Set("attendeeAvailability", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MeetingTimeSuggestion) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetConfidence sets the confidence property value. A percentage that represents the likelhood of all the attendees attending.
func (m *MeetingTimeSuggestion) SetConfidence(value *float64)() {
    err := m.GetBackingStore().Set("confidence", value)
    if err != nil {
        panic(err)
    }
}
// SetLocations sets the locations property value. An array that specifies the name and geographic location of each meeting location for this meeting suggestion.
func (m *MeetingTimeSuggestion) SetLocations(value []Locationable)() {
    err := m.GetBackingStore().Set("locations", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingTimeSlot sets the meetingTimeSlot property value. A time period suggested for the meeting.
func (m *MeetingTimeSuggestion) SetMeetingTimeSlot(value TimeSlotable)() {
    err := m.GetBackingStore().Set("meetingTimeSlot", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MeetingTimeSuggestion) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrder sets the order property value. Order of meeting time suggestions sorted by their computed confidence value from high to low, then by chronology if there are suggestions with the same confidence.
func (m *MeetingTimeSuggestion) SetOrder(value *int32)() {
    err := m.GetBackingStore().Set("order", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganizerAvailability sets the organizerAvailability property value. Availability of the meeting organizer for this meeting suggestion. The possible values are: free, tentative, busy, oof, workingElsewhere, unknown.
func (m *MeetingTimeSuggestion) SetOrganizerAvailability(value *FreeBusyStatus)() {
    err := m.GetBackingStore().Set("organizerAvailability", value)
    if err != nil {
        panic(err)
    }
}
// SetSuggestionReason sets the suggestionReason property value. Reason for suggesting the meeting time.
func (m *MeetingTimeSuggestion) SetSuggestionReason(value *string)() {
    err := m.GetBackingStore().Set("suggestionReason", value)
    if err != nil {
        panic(err)
    }
}
type MeetingTimeSuggestionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttendeeAvailability()([]AttendeeAvailabilityable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetConfidence()(*float64)
    GetLocations()([]Locationable)
    GetMeetingTimeSlot()(TimeSlotable)
    GetOdataType()(*string)
    GetOrder()(*int32)
    GetOrganizerAvailability()(*FreeBusyStatus)
    GetSuggestionReason()(*string)
    SetAttendeeAvailability(value []AttendeeAvailabilityable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetConfidence(value *float64)()
    SetLocations(value []Locationable)()
    SetMeetingTimeSlot(value TimeSlotable)()
    SetOdataType(value *string)()
    SetOrder(value *int32)()
    SetOrganizerAvailability(value *FreeBusyStatus)()
    SetSuggestionReason(value *string)()
}
