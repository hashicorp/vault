package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CalendarSharingMessage struct {
    Message
}
// NewCalendarSharingMessage instantiates a new CalendarSharingMessage and sets the default values.
func NewCalendarSharingMessage()(*CalendarSharingMessage) {
    m := &CalendarSharingMessage{
        Message: *NewMessage(),
    }
    odataTypeValue := "#microsoft.graph.calendarSharingMessage"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCalendarSharingMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCalendarSharingMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCalendarSharingMessage(), nil
}
// GetCanAccept gets the canAccept property value. The canAccept property
// returns a *bool when successful
func (m *CalendarSharingMessage) GetCanAccept()(*bool) {
    val, err := m.GetBackingStore().Get("canAccept")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CalendarSharingMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Message.GetFieldDeserializers()
    res["canAccept"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCanAccept(val)
        }
        return nil
    }
    res["sharingMessageAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCalendarSharingMessageActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharingMessageAction(val.(CalendarSharingMessageActionable))
        }
        return nil
    }
    res["sharingMessageActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCalendarSharingMessageActionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CalendarSharingMessageActionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CalendarSharingMessageActionable)
                }
            }
            m.SetSharingMessageActions(res)
        }
        return nil
    }
    res["suggestedCalendarName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuggestedCalendarName(val)
        }
        return nil
    }
    return res
}
// GetSharingMessageAction gets the sharingMessageAction property value. The sharingMessageAction property
// returns a CalendarSharingMessageActionable when successful
func (m *CalendarSharingMessage) GetSharingMessageAction()(CalendarSharingMessageActionable) {
    val, err := m.GetBackingStore().Get("sharingMessageAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(CalendarSharingMessageActionable)
    }
    return nil
}
// GetSharingMessageActions gets the sharingMessageActions property value. The sharingMessageActions property
// returns a []CalendarSharingMessageActionable when successful
func (m *CalendarSharingMessage) GetSharingMessageActions()([]CalendarSharingMessageActionable) {
    val, err := m.GetBackingStore().Get("sharingMessageActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CalendarSharingMessageActionable)
    }
    return nil
}
// GetSuggestedCalendarName gets the suggestedCalendarName property value. The suggestedCalendarName property
// returns a *string when successful
func (m *CalendarSharingMessage) GetSuggestedCalendarName()(*string) {
    val, err := m.GetBackingStore().Get("suggestedCalendarName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CalendarSharingMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Message.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("canAccept", m.GetCanAccept())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharingMessageAction", m.GetSharingMessageAction())
        if err != nil {
            return err
        }
    }
    if m.GetSharingMessageActions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharingMessageActions()))
        for i, v := range m.GetSharingMessageActions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharingMessageActions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("suggestedCalendarName", m.GetSuggestedCalendarName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCanAccept sets the canAccept property value. The canAccept property
func (m *CalendarSharingMessage) SetCanAccept(value *bool)() {
    err := m.GetBackingStore().Set("canAccept", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingMessageAction sets the sharingMessageAction property value. The sharingMessageAction property
func (m *CalendarSharingMessage) SetSharingMessageAction(value CalendarSharingMessageActionable)() {
    err := m.GetBackingStore().Set("sharingMessageAction", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingMessageActions sets the sharingMessageActions property value. The sharingMessageActions property
func (m *CalendarSharingMessage) SetSharingMessageActions(value []CalendarSharingMessageActionable)() {
    err := m.GetBackingStore().Set("sharingMessageActions", value)
    if err != nil {
        panic(err)
    }
}
// SetSuggestedCalendarName sets the suggestedCalendarName property value. The suggestedCalendarName property
func (m *CalendarSharingMessage) SetSuggestedCalendarName(value *string)() {
    err := m.GetBackingStore().Set("suggestedCalendarName", value)
    if err != nil {
        panic(err)
    }
}
type CalendarSharingMessageable interface {
    Messageable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCanAccept()(*bool)
    GetSharingMessageAction()(CalendarSharingMessageActionable)
    GetSharingMessageActions()([]CalendarSharingMessageActionable)
    GetSuggestedCalendarName()(*string)
    SetCanAccept(value *bool)()
    SetSharingMessageAction(value CalendarSharingMessageActionable)()
    SetSharingMessageActions(value []CalendarSharingMessageActionable)()
    SetSuggestedCalendarName(value *string)()
}
