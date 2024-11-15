package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CalendarPermission struct {
    Entity
}
// NewCalendarPermission instantiates a new CalendarPermission and sets the default values.
func NewCalendarPermission()(*CalendarPermission) {
    m := &CalendarPermission{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCalendarPermissionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCalendarPermissionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCalendarPermission(), nil
}
// GetAllowedRoles gets the allowedRoles property value. List of allowed sharing or delegating permission levels for the calendar. Possible values are: none, freeBusyRead, limitedRead, read, write, delegateWithoutPrivateEventAccess, delegateWithPrivateEventAccess, custom.
// returns a []CalendarRoleType when successful
func (m *CalendarPermission) GetAllowedRoles()([]CalendarRoleType) {
    val, err := m.GetBackingStore().Get("allowedRoles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CalendarRoleType)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. Represents a share recipient or delegate who has access to the calendar. For the 'My Organization' share recipient, the address property is null. Read-only.
// returns a EmailAddressable when successful
func (m *CalendarPermission) GetEmailAddress()(EmailAddressable) {
    val, err := m.GetBackingStore().Get("emailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailAddressable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CalendarPermission) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowedRoles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseCalendarRoleType)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CalendarRoleType, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*CalendarRoleType))
                }
            }
            m.SetAllowedRoles(res)
        }
        return nil
    }
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val.(EmailAddressable))
        }
        return nil
    }
    res["isInsideOrganization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInsideOrganization(val)
        }
        return nil
    }
    res["isRemovable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRemovable(val)
        }
        return nil
    }
    res["role"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCalendarRoleType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRole(val.(*CalendarRoleType))
        }
        return nil
    }
    return res
}
// GetIsInsideOrganization gets the isInsideOrganization property value. True if the user in context (recipient or delegate) is inside the same organization as the calendar owner.
// returns a *bool when successful
func (m *CalendarPermission) GetIsInsideOrganization()(*bool) {
    val, err := m.GetBackingStore().Get("isInsideOrganization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRemovable gets the isRemovable property value. True if the user can be removed from the list of recipients or delegates for the specified calendar, false otherwise. The 'My organization' user determines the permissions other people within your organization have to the given calendar. You can't remove 'My organization' as a share recipient to a calendar.
// returns a *bool when successful
func (m *CalendarPermission) GetIsRemovable()(*bool) {
    val, err := m.GetBackingStore().Get("isRemovable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRole gets the role property value. Current permission level of the calendar share recipient or delegate.
// returns a *CalendarRoleType when successful
func (m *CalendarPermission) GetRole()(*CalendarRoleType) {
    val, err := m.GetBackingStore().Get("role")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CalendarRoleType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CalendarPermission) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedRoles() != nil {
        err = writer.WriteCollectionOfStringValues("allowedRoles", SerializeCalendarRoleType(m.GetAllowedRoles()))
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isInsideOrganization", m.GetIsInsideOrganization())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isRemovable", m.GetIsRemovable())
        if err != nil {
            return err
        }
    }
    if m.GetRole() != nil {
        cast := (*m.GetRole()).String()
        err = writer.WriteStringValue("role", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedRoles sets the allowedRoles property value. List of allowed sharing or delegating permission levels for the calendar. Possible values are: none, freeBusyRead, limitedRead, read, write, delegateWithoutPrivateEventAccess, delegateWithPrivateEventAccess, custom.
func (m *CalendarPermission) SetAllowedRoles(value []CalendarRoleType)() {
    err := m.GetBackingStore().Set("allowedRoles", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. Represents a share recipient or delegate who has access to the calendar. For the 'My Organization' share recipient, the address property is null. Read-only.
func (m *CalendarPermission) SetEmailAddress(value EmailAddressable)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInsideOrganization sets the isInsideOrganization property value. True if the user in context (recipient or delegate) is inside the same organization as the calendar owner.
func (m *CalendarPermission) SetIsInsideOrganization(value *bool)() {
    err := m.GetBackingStore().Set("isInsideOrganization", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRemovable sets the isRemovable property value. True if the user can be removed from the list of recipients or delegates for the specified calendar, false otherwise. The 'My organization' user determines the permissions other people within your organization have to the given calendar. You can't remove 'My organization' as a share recipient to a calendar.
func (m *CalendarPermission) SetIsRemovable(value *bool)() {
    err := m.GetBackingStore().Set("isRemovable", value)
    if err != nil {
        panic(err)
    }
}
// SetRole sets the role property value. Current permission level of the calendar share recipient or delegate.
func (m *CalendarPermission) SetRole(value *CalendarRoleType)() {
    err := m.GetBackingStore().Set("role", value)
    if err != nil {
        panic(err)
    }
}
type CalendarPermissionable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedRoles()([]CalendarRoleType)
    GetEmailAddress()(EmailAddressable)
    GetIsInsideOrganization()(*bool)
    GetIsRemovable()(*bool)
    GetRole()(*CalendarRoleType)
    SetAllowedRoles(value []CalendarRoleType)()
    SetEmailAddress(value EmailAddressable)()
    SetIsInsideOrganization(value *bool)()
    SetIsRemovable(value *bool)()
    SetRole(value *CalendarRoleType)()
}
