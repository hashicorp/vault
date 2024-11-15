package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrivilegedAccessScheduleRequest struct {
    Request
}
// NewPrivilegedAccessScheduleRequest instantiates a new PrivilegedAccessScheduleRequest and sets the default values.
func NewPrivilegedAccessScheduleRequest()(*PrivilegedAccessScheduleRequest) {
    m := &PrivilegedAccessScheduleRequest{
        Request: *NewRequest(),
    }
    return m
}
// CreatePrivilegedAccessScheduleRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrivilegedAccessScheduleRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.privilegedAccessGroupAssignmentScheduleRequest":
                        return NewPrivilegedAccessGroupAssignmentScheduleRequest(), nil
                    case "#microsoft.graph.privilegedAccessGroupEligibilityScheduleRequest":
                        return NewPrivilegedAccessGroupEligibilityScheduleRequest(), nil
                }
            }
        }
    }
    return NewPrivilegedAccessScheduleRequest(), nil
}
// GetAction gets the action property value. Represents the type of operation on the group membership or ownership assignment request. The possible values are: adminAssign, adminUpdate, adminRemove, selfActivate, selfDeactivate, adminExtend, adminRenew. adminAssign: For administrators to assign group membership or ownership to principals.adminRemove: For administrators to remove principals from group membership or ownership. adminUpdate: For administrators to change existing group membership or ownership assignments.adminExtend: For administrators to extend expiring assignments.adminRenew: For administrators to renew expired assignments.selfActivate: For principals to activate their assignments.selfDeactivate: For principals to deactivate their active assignments.
// returns a *ScheduleRequestActions when successful
func (m *PrivilegedAccessScheduleRequest) GetAction()(*ScheduleRequestActions) {
    val, err := m.GetBackingStore().Get("action")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ScheduleRequestActions)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrivilegedAccessScheduleRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Request.GetFieldDeserializers()
    res["action"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseScheduleRequestActions)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAction(val.(*ScheduleRequestActions))
        }
        return nil
    }
    res["isValidationOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsValidationOnly(val)
        }
        return nil
    }
    res["justification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJustification(val)
        }
        return nil
    }
    res["scheduleInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRequestScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduleInfo(val.(RequestScheduleable))
        }
        return nil
    }
    res["ticketInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTicketInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTicketInfo(val.(TicketInfoable))
        }
        return nil
    }
    return res
}
// GetIsValidationOnly gets the isValidationOnly property value. Determines whether the call is a validation or an actual call. Only set this property if you want to check whether an activation is subject to additional rules like MFA before actually submitting the request.
// returns a *bool when successful
func (m *PrivilegedAccessScheduleRequest) GetIsValidationOnly()(*bool) {
    val, err := m.GetBackingStore().Get("isValidationOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJustification gets the justification property value. A message provided by users and administrators when create they create the privilegedAccessGroupAssignmentScheduleRequest object.
// returns a *string when successful
func (m *PrivilegedAccessScheduleRequest) GetJustification()(*string) {
    val, err := m.GetBackingStore().Get("justification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScheduleInfo gets the scheduleInfo property value. The period of the group membership or ownership assignment. Recurring schedules are currently unsupported.
// returns a RequestScheduleable when successful
func (m *PrivilegedAccessScheduleRequest) GetScheduleInfo()(RequestScheduleable) {
    val, err := m.GetBackingStore().Get("scheduleInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RequestScheduleable)
    }
    return nil
}
// GetTicketInfo gets the ticketInfo property value. Ticket details linked to the group membership or ownership assignment request including details of the ticket number and ticket system.
// returns a TicketInfoable when successful
func (m *PrivilegedAccessScheduleRequest) GetTicketInfo()(TicketInfoable) {
    val, err := m.GetBackingStore().Get("ticketInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TicketInfoable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrivilegedAccessScheduleRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Request.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAction() != nil {
        cast := (*m.GetAction()).String()
        err = writer.WriteStringValue("action", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isValidationOnly", m.GetIsValidationOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("justification", m.GetJustification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("scheduleInfo", m.GetScheduleInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("ticketInfo", m.GetTicketInfo())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAction sets the action property value. Represents the type of operation on the group membership or ownership assignment request. The possible values are: adminAssign, adminUpdate, adminRemove, selfActivate, selfDeactivate, adminExtend, adminRenew. adminAssign: For administrators to assign group membership or ownership to principals.adminRemove: For administrators to remove principals from group membership or ownership. adminUpdate: For administrators to change existing group membership or ownership assignments.adminExtend: For administrators to extend expiring assignments.adminRenew: For administrators to renew expired assignments.selfActivate: For principals to activate their assignments.selfDeactivate: For principals to deactivate their active assignments.
func (m *PrivilegedAccessScheduleRequest) SetAction(value *ScheduleRequestActions)() {
    err := m.GetBackingStore().Set("action", value)
    if err != nil {
        panic(err)
    }
}
// SetIsValidationOnly sets the isValidationOnly property value. Determines whether the call is a validation or an actual call. Only set this property if you want to check whether an activation is subject to additional rules like MFA before actually submitting the request.
func (m *PrivilegedAccessScheduleRequest) SetIsValidationOnly(value *bool)() {
    err := m.GetBackingStore().Set("isValidationOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetJustification sets the justification property value. A message provided by users and administrators when create they create the privilegedAccessGroupAssignmentScheduleRequest object.
func (m *PrivilegedAccessScheduleRequest) SetJustification(value *string)() {
    err := m.GetBackingStore().Set("justification", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleInfo sets the scheduleInfo property value. The period of the group membership or ownership assignment. Recurring schedules are currently unsupported.
func (m *PrivilegedAccessScheduleRequest) SetScheduleInfo(value RequestScheduleable)() {
    err := m.GetBackingStore().Set("scheduleInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetTicketInfo sets the ticketInfo property value. Ticket details linked to the group membership or ownership assignment request including details of the ticket number and ticket system.
func (m *PrivilegedAccessScheduleRequest) SetTicketInfo(value TicketInfoable)() {
    err := m.GetBackingStore().Set("ticketInfo", value)
    if err != nil {
        panic(err)
    }
}
type PrivilegedAccessScheduleRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Requestable
    GetAction()(*ScheduleRequestActions)
    GetIsValidationOnly()(*bool)
    GetJustification()(*string)
    GetScheduleInfo()(RequestScheduleable)
    GetTicketInfo()(TicketInfoable)
    SetAction(value *ScheduleRequestActions)()
    SetIsValidationOnly(value *bool)()
    SetJustification(value *string)()
    SetScheduleInfo(value RequestScheduleable)()
    SetTicketInfo(value TicketInfoable)()
}
