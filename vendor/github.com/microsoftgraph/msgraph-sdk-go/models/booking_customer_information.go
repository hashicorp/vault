package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type BookingCustomerInformation struct {
    BookingCustomerInformationBase
}
// NewBookingCustomerInformation instantiates a new BookingCustomerInformation and sets the default values.
func NewBookingCustomerInformation()(*BookingCustomerInformation) {
    m := &BookingCustomerInformation{
        BookingCustomerInformationBase: *NewBookingCustomerInformationBase(),
    }
    odataTypeValue := "#microsoft.graph.bookingCustomerInformation"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateBookingCustomerInformationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBookingCustomerInformationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBookingCustomerInformation(), nil
}
// GetCustomerId gets the customerId property value. The ID of the bookingCustomer for this appointment. If no ID is specified when an appointment is created, then a new bookingCustomer object is created. Once set, you should consider the customerId immutable.
// returns a *string when successful
func (m *BookingCustomerInformation) GetCustomerId()(*string) {
    val, err := m.GetBackingStore().Get("customerId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCustomQuestionAnswers gets the customQuestionAnswers property value. It consists of the list of custom questions and answers given by the customer as part of the appointment
// returns a []BookingQuestionAnswerable when successful
func (m *BookingCustomerInformation) GetCustomQuestionAnswers()([]BookingQuestionAnswerable) {
    val, err := m.GetBackingStore().Get("customQuestionAnswers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BookingQuestionAnswerable)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. The SMTP address of the bookingCustomer who is booking the appointment
// returns a *string when successful
func (m *BookingCustomerInformation) GetEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("emailAddress")
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
func (m *BookingCustomerInformation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BookingCustomerInformationBase.GetFieldDeserializers()
    res["customerId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomerId(val)
        }
        return nil
    }
    res["customQuestionAnswers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBookingQuestionAnswerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BookingQuestionAnswerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BookingQuestionAnswerable)
                }
            }
            m.SetCustomQuestionAnswers(res)
        }
        return nil
    }
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val)
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(Locationable))
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["notes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotes(val)
        }
        return nil
    }
    res["phone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhone(val)
        }
        return nil
    }
    res["timeZone"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeZone(val)
        }
        return nil
    }
    return res
}
// GetLocation gets the location property value. Represents location information for the bookingCustomer who is booking the appointment.
// returns a Locationable when successful
func (m *BookingCustomerInformation) GetLocation()(Locationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Locationable)
    }
    return nil
}
// GetName gets the name property value. The customer's name.
// returns a *string when successful
func (m *BookingCustomerInformation) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotes gets the notes property value. Notes from the customer associated with this appointment. You can get the value only when reading this bookingAppointment by its ID. You can set this property only when initially creating an appointment with a new customer. After that point, the value is computed from the customer represented by the customerId.
// returns a *string when successful
func (m *BookingCustomerInformation) GetNotes()(*string) {
    val, err := m.GetBackingStore().Get("notes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhone gets the phone property value. The customer's phone number.
// returns a *string when successful
func (m *BookingCustomerInformation) GetPhone()(*string) {
    val, err := m.GetBackingStore().Get("phone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTimeZone gets the timeZone property value. The time zone of the customer. For a list of possible values, see dateTimeTimeZone.
// returns a *string when successful
func (m *BookingCustomerInformation) GetTimeZone()(*string) {
    val, err := m.GetBackingStore().Get("timeZone")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BookingCustomerInformation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BookingCustomerInformationBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("customerId", m.GetCustomerId())
        if err != nil {
            return err
        }
    }
    if m.GetCustomQuestionAnswers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomQuestionAnswers()))
        for i, v := range m.GetCustomQuestionAnswers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customQuestionAnswers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notes", m.GetNotes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("phone", m.GetPhone())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("timeZone", m.GetTimeZone())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCustomerId sets the customerId property value. The ID of the bookingCustomer for this appointment. If no ID is specified when an appointment is created, then a new bookingCustomer object is created. Once set, you should consider the customerId immutable.
func (m *BookingCustomerInformation) SetCustomerId(value *string)() {
    err := m.GetBackingStore().Set("customerId", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomQuestionAnswers sets the customQuestionAnswers property value. It consists of the list of custom questions and answers given by the customer as part of the appointment
func (m *BookingCustomerInformation) SetCustomQuestionAnswers(value []BookingQuestionAnswerable)() {
    err := m.GetBackingStore().Set("customQuestionAnswers", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. The SMTP address of the bookingCustomer who is booking the appointment
func (m *BookingCustomerInformation) SetEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. Represents location information for the bookingCustomer who is booking the appointment.
func (m *BookingCustomerInformation) SetLocation(value Locationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The customer's name.
func (m *BookingCustomerInformation) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNotes sets the notes property value. Notes from the customer associated with this appointment. You can get the value only when reading this bookingAppointment by its ID. You can set this property only when initially creating an appointment with a new customer. After that point, the value is computed from the customer represented by the customerId.
func (m *BookingCustomerInformation) SetNotes(value *string)() {
    err := m.GetBackingStore().Set("notes", value)
    if err != nil {
        panic(err)
    }
}
// SetPhone sets the phone property value. The customer's phone number.
func (m *BookingCustomerInformation) SetPhone(value *string)() {
    err := m.GetBackingStore().Set("phone", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeZone sets the timeZone property value. The time zone of the customer. For a list of possible values, see dateTimeTimeZone.
func (m *BookingCustomerInformation) SetTimeZone(value *string)() {
    err := m.GetBackingStore().Set("timeZone", value)
    if err != nil {
        panic(err)
    }
}
type BookingCustomerInformationable interface {
    BookingCustomerInformationBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCustomerId()(*string)
    GetCustomQuestionAnswers()([]BookingQuestionAnswerable)
    GetEmailAddress()(*string)
    GetLocation()(Locationable)
    GetName()(*string)
    GetNotes()(*string)
    GetPhone()(*string)
    GetTimeZone()(*string)
    SetCustomerId(value *string)()
    SetCustomQuestionAnswers(value []BookingQuestionAnswerable)()
    SetEmailAddress(value *string)()
    SetLocation(value Locationable)()
    SetName(value *string)()
    SetNotes(value *string)()
    SetPhone(value *string)()
    SetTimeZone(value *string)()
}
