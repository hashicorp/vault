package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Reminder struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewReminder instantiates a new Reminder and sets the default values.
func NewReminder()(*Reminder) {
    m := &Reminder{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateReminderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateReminderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewReminder(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Reminder) GetAdditionalData()(map[string]any) {
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
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *Reminder) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetChangeKey gets the changeKey property value. Identifies the version of the reminder. Every time the reminder is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object.
// returns a *string when successful
func (m *Reminder) GetChangeKey()(*string) {
    val, err := m.GetBackingStore().Get("changeKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEventEndTime gets the eventEndTime property value. The date, time and time zone that the event ends.
// returns a DateTimeTimeZoneable when successful
func (m *Reminder) GetEventEndTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("eventEndTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetEventId gets the eventId property value. The unique ID of the event. Read only.
// returns a *string when successful
func (m *Reminder) GetEventId()(*string) {
    val, err := m.GetBackingStore().Get("eventId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEventLocation gets the eventLocation property value. The location of the event.
// returns a Locationable when successful
func (m *Reminder) GetEventLocation()(Locationable) {
    val, err := m.GetBackingStore().Get("eventLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Locationable)
    }
    return nil
}
// GetEventStartTime gets the eventStartTime property value. The date, time, and time zone that the event starts.
// returns a DateTimeTimeZoneable when successful
func (m *Reminder) GetEventStartTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("eventStartTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetEventSubject gets the eventSubject property value. The text of the event's subject line.
// returns a *string when successful
func (m *Reminder) GetEventSubject()(*string) {
    val, err := m.GetBackingStore().Get("eventSubject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEventWebLink gets the eventWebLink property value. The URL to open the event in Outlook on the web.The event opens in the browser if you're logged in to your mailbox via Outlook on the web. You're prompted to log in if you aren't already logged in with the browser.This URL can't be accessed from within an iFrame.
// returns a *string when successful
func (m *Reminder) GetEventWebLink()(*string) {
    val, err := m.GetBackingStore().Get("eventWebLink")
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
func (m *Reminder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["changeKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChangeKey(val)
        }
        return nil
    }
    res["eventEndTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventEndTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["eventId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventId(val)
        }
        return nil
    }
    res["eventLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventLocation(val.(Locationable))
        }
        return nil
    }
    res["eventStartTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventStartTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["eventSubject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventSubject(val)
        }
        return nil
    }
    res["eventWebLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventWebLink(val)
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
    res["reminderFireTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReminderFireTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Reminder) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReminderFireTime gets the reminderFireTime property value. The date, time, and time zone that the reminder is set to occur.
// returns a DateTimeTimeZoneable when successful
func (m *Reminder) GetReminderFireTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("reminderFireTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Reminder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("changeKey", m.GetChangeKey())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventEndTime", m.GetEventEndTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventId", m.GetEventId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventLocation", m.GetEventLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("eventStartTime", m.GetEventStartTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventSubject", m.GetEventSubject())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("eventWebLink", m.GetEventWebLink())
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
        err := writer.WriteObjectValue("reminderFireTime", m.GetReminderFireTime())
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
func (m *Reminder) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Reminder) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetChangeKey sets the changeKey property value. Identifies the version of the reminder. Every time the reminder is changed, changeKey changes as well. This allows Exchange to apply changes to the correct version of the object.
func (m *Reminder) SetChangeKey(value *string)() {
    err := m.GetBackingStore().Set("changeKey", value)
    if err != nil {
        panic(err)
    }
}
// SetEventEndTime sets the eventEndTime property value. The date, time and time zone that the event ends.
func (m *Reminder) SetEventEndTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("eventEndTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEventId sets the eventId property value. The unique ID of the event. Read only.
func (m *Reminder) SetEventId(value *string)() {
    err := m.GetBackingStore().Set("eventId", value)
    if err != nil {
        panic(err)
    }
}
// SetEventLocation sets the eventLocation property value. The location of the event.
func (m *Reminder) SetEventLocation(value Locationable)() {
    err := m.GetBackingStore().Set("eventLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetEventStartTime sets the eventStartTime property value. The date, time, and time zone that the event starts.
func (m *Reminder) SetEventStartTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("eventStartTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEventSubject sets the eventSubject property value. The text of the event's subject line.
func (m *Reminder) SetEventSubject(value *string)() {
    err := m.GetBackingStore().Set("eventSubject", value)
    if err != nil {
        panic(err)
    }
}
// SetEventWebLink sets the eventWebLink property value. The URL to open the event in Outlook on the web.The event opens in the browser if you're logged in to your mailbox via Outlook on the web. You're prompted to log in if you aren't already logged in with the browser.This URL can't be accessed from within an iFrame.
func (m *Reminder) SetEventWebLink(value *string)() {
    err := m.GetBackingStore().Set("eventWebLink", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Reminder) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetReminderFireTime sets the reminderFireTime property value. The date, time, and time zone that the reminder is set to occur.
func (m *Reminder) SetReminderFireTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("reminderFireTime", value)
    if err != nil {
        panic(err)
    }
}
type Reminderable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetChangeKey()(*string)
    GetEventEndTime()(DateTimeTimeZoneable)
    GetEventId()(*string)
    GetEventLocation()(Locationable)
    GetEventStartTime()(DateTimeTimeZoneable)
    GetEventSubject()(*string)
    GetEventWebLink()(*string)
    GetOdataType()(*string)
    GetReminderFireTime()(DateTimeTimeZoneable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetChangeKey(value *string)()
    SetEventEndTime(value DateTimeTimeZoneable)()
    SetEventId(value *string)()
    SetEventLocation(value Locationable)()
    SetEventStartTime(value DateTimeTimeZoneable)()
    SetEventSubject(value *string)()
    SetEventWebLink(value *string)()
    SetOdataType(value *string)()
    SetReminderFireTime(value DateTimeTimeZoneable)()
}
