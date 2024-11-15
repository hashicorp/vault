package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ScheduleChangeRequest struct {
    ChangeTrackedEntity
}
// NewScheduleChangeRequest instantiates a new ScheduleChangeRequest and sets the default values.
func NewScheduleChangeRequest()(*ScheduleChangeRequest) {
    m := &ScheduleChangeRequest{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.scheduleChangeRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateScheduleChangeRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateScheduleChangeRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.offerShiftRequest":
                        return NewOfferShiftRequest(), nil
                    case "#microsoft.graph.openShiftChangeRequest":
                        return NewOpenShiftChangeRequest(), nil
                    case "#microsoft.graph.swapShiftsChangeRequest":
                        return NewSwapShiftsChangeRequest(), nil
                    case "#microsoft.graph.timeOffRequest":
                        return NewTimeOffRequest(), nil
                }
            }
        }
    }
    return NewScheduleChangeRequest(), nil
}
// GetAssignedTo gets the assignedTo property value. The assignedTo property
// returns a *ScheduleChangeRequestActor when successful
func (m *ScheduleChangeRequest) GetAssignedTo()(*ScheduleChangeRequestActor) {
    val, err := m.GetBackingStore().Get("assignedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ScheduleChangeRequestActor)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ScheduleChangeRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
    res["assignedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseScheduleChangeRequestActor)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignedTo(val.(*ScheduleChangeRequestActor))
        }
        return nil
    }
    res["managerActionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagerActionDateTime(val)
        }
        return nil
    }
    res["managerActionMessage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagerActionMessage(val)
        }
        return nil
    }
    res["managerUserId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagerUserId(val)
        }
        return nil
    }
    res["senderDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSenderDateTime(val)
        }
        return nil
    }
    res["senderMessage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSenderMessage(val)
        }
        return nil
    }
    res["senderUserId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSenderUserId(val)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseScheduleChangeState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*ScheduleChangeState))
        }
        return nil
    }
    return res
}
// GetManagerActionDateTime gets the managerActionDateTime property value. The managerActionDateTime property
// returns a *Time when successful
func (m *ScheduleChangeRequest) GetManagerActionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("managerActionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetManagerActionMessage gets the managerActionMessage property value. The managerActionMessage property
// returns a *string when successful
func (m *ScheduleChangeRequest) GetManagerActionMessage()(*string) {
    val, err := m.GetBackingStore().Get("managerActionMessage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetManagerUserId gets the managerUserId property value. The managerUserId property
// returns a *string when successful
func (m *ScheduleChangeRequest) GetManagerUserId()(*string) {
    val, err := m.GetBackingStore().Get("managerUserId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSenderDateTime gets the senderDateTime property value. The senderDateTime property
// returns a *Time when successful
func (m *ScheduleChangeRequest) GetSenderDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("senderDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSenderMessage gets the senderMessage property value. The senderMessage property
// returns a *string when successful
func (m *ScheduleChangeRequest) GetSenderMessage()(*string) {
    val, err := m.GetBackingStore().Get("senderMessage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSenderUserId gets the senderUserId property value. The senderUserId property
// returns a *string when successful
func (m *ScheduleChangeRequest) GetSenderUserId()(*string) {
    val, err := m.GetBackingStore().Get("senderUserId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *ScheduleChangeState when successful
func (m *ScheduleChangeRequest) GetState()(*ScheduleChangeState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ScheduleChangeState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ScheduleChangeRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignedTo() != nil {
        cast := (*m.GetAssignedTo()).String()
        err = writer.WriteStringValue("assignedTo", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("managerActionMessage", m.GetManagerActionMessage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("senderMessage", m.GetSenderMessage())
        if err != nil {
            return err
        }
    }
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err = writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignedTo sets the assignedTo property value. The assignedTo property
func (m *ScheduleChangeRequest) SetAssignedTo(value *ScheduleChangeRequestActor)() {
    err := m.GetBackingStore().Set("assignedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetManagerActionDateTime sets the managerActionDateTime property value. The managerActionDateTime property
func (m *ScheduleChangeRequest) SetManagerActionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("managerActionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetManagerActionMessage sets the managerActionMessage property value. The managerActionMessage property
func (m *ScheduleChangeRequest) SetManagerActionMessage(value *string)() {
    err := m.GetBackingStore().Set("managerActionMessage", value)
    if err != nil {
        panic(err)
    }
}
// SetManagerUserId sets the managerUserId property value. The managerUserId property
func (m *ScheduleChangeRequest) SetManagerUserId(value *string)() {
    err := m.GetBackingStore().Set("managerUserId", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderDateTime sets the senderDateTime property value. The senderDateTime property
func (m *ScheduleChangeRequest) SetSenderDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("senderDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderMessage sets the senderMessage property value. The senderMessage property
func (m *ScheduleChangeRequest) SetSenderMessage(value *string)() {
    err := m.GetBackingStore().Set("senderMessage", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderUserId sets the senderUserId property value. The senderUserId property
func (m *ScheduleChangeRequest) SetSenderUserId(value *string)() {
    err := m.GetBackingStore().Set("senderUserId", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *ScheduleChangeRequest) SetState(value *ScheduleChangeState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
type ScheduleChangeRequestable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignedTo()(*ScheduleChangeRequestActor)
    GetManagerActionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetManagerActionMessage()(*string)
    GetManagerUserId()(*string)
    GetSenderDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSenderMessage()(*string)
    GetSenderUserId()(*string)
    GetState()(*ScheduleChangeState)
    SetAssignedTo(value *ScheduleChangeRequestActor)()
    SetManagerActionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetManagerActionMessage(value *string)()
    SetManagerUserId(value *string)()
    SetSenderDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSenderMessage(value *string)()
    SetSenderUserId(value *string)()
    SetState(value *ScheduleChangeState)()
}
