package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleEligibilityScheduleRequest struct {
    Request
}
// NewUnifiedRoleEligibilityScheduleRequest instantiates a new UnifiedRoleEligibilityScheduleRequest and sets the default values.
func NewUnifiedRoleEligibilityScheduleRequest()(*UnifiedRoleEligibilityScheduleRequest) {
    m := &UnifiedRoleEligibilityScheduleRequest{
        Request: *NewRequest(),
    }
    return m
}
// CreateUnifiedRoleEligibilityScheduleRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleEligibilityScheduleRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleEligibilityScheduleRequest(), nil
}
// GetAction gets the action property value. Represents the type of operation on the role eligibility request. The possible values are: adminAssign, adminUpdate, adminRemove, selfActivate, selfDeactivate, adminExtend, adminRenew, selfExtend, selfRenew, unknownFutureValue. adminAssign: For administrators to assign eligible roles to principals.adminRemove: For administrators to remove eligible roles from principals. adminUpdate: For administrators to change existing role eligibilities.adminExtend: For administrators to extend expiring role eligibilities.adminRenew: For administrators to renew expired eligibilities.selfActivate: For users to activate their assignments.selfDeactivate: For users to deactivate their active assignments.selfExtend: For users to request to extend their expiring assignments.selfRenew: For users to request to renew their expired assignments.
// returns a *UnifiedRoleScheduleRequestActions when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetAction()(*UnifiedRoleScheduleRequestActions) {
    val, err := m.GetBackingStore().Get("action")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*UnifiedRoleScheduleRequestActions)
    }
    return nil
}
// GetAppScope gets the appScope property value. Read-only property with details of the app-specific scope when the role eligibility is scoped to an app. Nullable. Supports $expand.
// returns a AppScopeable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetAppScope()(AppScopeable) {
    val, err := m.GetBackingStore().Get("appScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AppScopeable)
    }
    return nil
}
// GetAppScopeId gets the appScopeId property value. Identifier of the app-specific scope when the role eligibility is scoped to an app. The scope of a role eligibility determines the set of resources for which the principal is eligible to access. App scopes are scopes that are defined and understood by this application only. Use / for tenant-wide app scopes. Use directoryScopeId to limit the scope to particular directory objects, for example, administrative units. Supports $filter (eq, ne, and on null values).
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetAppScopeId()(*string) {
    val, err := m.GetBackingStore().Get("appScopeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDirectoryScope gets the directoryScope property value. The directory object that is the scope of the role eligibility. Read-only. Supports $expand.
// returns a DirectoryObjectable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetDirectoryScope()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("directoryScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetDirectoryScopeId gets the directoryScopeId property value. Identifier of the directory object representing the scope of the role eligibility. The scope of a role eligibility determines the set of resources for which the principal has been granted access. Directory scopes are shared scopes stored in the directory that are understood by multiple applications. Use / for tenant-wide scope. Use appScopeId to limit the scope to an application only. Supports $filter (eq, ne, and on null values).
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetDirectoryScopeId()(*string) {
    val, err := m.GetBackingStore().Get("directoryScopeId")
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
func (m *UnifiedRoleEligibilityScheduleRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Request.GetFieldDeserializers()
    res["action"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseUnifiedRoleScheduleRequestActions)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAction(val.(*UnifiedRoleScheduleRequestActions))
        }
        return nil
    }
    res["appScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAppScopeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppScope(val.(AppScopeable))
        }
        return nil
    }
    res["appScopeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppScopeId(val)
        }
        return nil
    }
    res["directoryScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirectoryScope(val.(DirectoryObjectable))
        }
        return nil
    }
    res["directoryScopeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDirectoryScopeId(val)
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
    res["principal"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipal(val.(DirectoryObjectable))
        }
        return nil
    }
    res["principalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipalId(val)
        }
        return nil
    }
    res["roleDefinition"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnifiedRoleDefinitionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleDefinition(val.(UnifiedRoleDefinitionable))
        }
        return nil
    }
    res["roleDefinitionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleDefinitionId(val)
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
    res["targetSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnifiedRoleEligibilityScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetSchedule(val.(UnifiedRoleEligibilityScheduleable))
        }
        return nil
    }
    res["targetScheduleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetScheduleId(val)
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
func (m *UnifiedRoleEligibilityScheduleRequest) GetIsValidationOnly()(*bool) {
    val, err := m.GetBackingStore().Get("isValidationOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetJustification gets the justification property value. A message provided by users and administrators when create they create the unifiedRoleEligibilityScheduleRequest object.
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetJustification()(*string) {
    val, err := m.GetBackingStore().Get("justification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrincipal gets the principal property value. The principal that's getting a role eligibility through the request. Supports $expand.
// returns a DirectoryObjectable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetPrincipal()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("principal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetPrincipalId gets the principalId property value. Identifier of the principal that has been granted the role eligibility. Can be a user or a role-assignable group. You can grant only active assignments service principals.Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetPrincipalId()(*string) {
    val, err := m.GetBackingStore().Get("principalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoleDefinition gets the roleDefinition property value. Detailed information for the unifiedRoleDefinition object that is referenced through the roleDefinitionId property. Supports $expand.
// returns a UnifiedRoleDefinitionable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetRoleDefinition()(UnifiedRoleDefinitionable) {
    val, err := m.GetBackingStore().Get("roleDefinition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnifiedRoleDefinitionable)
    }
    return nil
}
// GetRoleDefinitionId gets the roleDefinitionId property value. Identifier of the unifiedRoleDefinition object that is being assigned to the principal. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetRoleDefinitionId()(*string) {
    val, err := m.GetBackingStore().Get("roleDefinitionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScheduleInfo gets the scheduleInfo property value. The period of the role eligibility. Recurring schedules are currently unsupported.
// returns a RequestScheduleable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetScheduleInfo()(RequestScheduleable) {
    val, err := m.GetBackingStore().Get("scheduleInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RequestScheduleable)
    }
    return nil
}
// GetTargetSchedule gets the targetSchedule property value. The schedule for a role eligibility that is referenced through the targetScheduleId property. Supports $expand.
// returns a UnifiedRoleEligibilityScheduleable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetTargetSchedule()(UnifiedRoleEligibilityScheduleable) {
    val, err := m.GetBackingStore().Get("targetSchedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnifiedRoleEligibilityScheduleable)
    }
    return nil
}
// GetTargetScheduleId gets the targetScheduleId property value. Identifier of the schedule object that's linked to the eligibility request. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetTargetScheduleId()(*string) {
    val, err := m.GetBackingStore().Get("targetScheduleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTicketInfo gets the ticketInfo property value. Ticket details linked to the role eligibility request including details of the ticket number and ticket system. Optional.
// returns a TicketInfoable when successful
func (m *UnifiedRoleEligibilityScheduleRequest) GetTicketInfo()(TicketInfoable) {
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
func (m *UnifiedRoleEligibilityScheduleRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("appScope", m.GetAppScope())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appScopeId", m.GetAppScopeId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("directoryScope", m.GetDirectoryScope())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("directoryScopeId", m.GetDirectoryScopeId())
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
        err = writer.WriteObjectValue("principal", m.GetPrincipal())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalId", m.GetPrincipalId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("roleDefinition", m.GetRoleDefinition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleDefinitionId", m.GetRoleDefinitionId())
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
        err = writer.WriteObjectValue("targetSchedule", m.GetTargetSchedule())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("targetScheduleId", m.GetTargetScheduleId())
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
// SetAction sets the action property value. Represents the type of operation on the role eligibility request. The possible values are: adminAssign, adminUpdate, adminRemove, selfActivate, selfDeactivate, adminExtend, adminRenew, selfExtend, selfRenew, unknownFutureValue. adminAssign: For administrators to assign eligible roles to principals.adminRemove: For administrators to remove eligible roles from principals. adminUpdate: For administrators to change existing role eligibilities.adminExtend: For administrators to extend expiring role eligibilities.adminRenew: For administrators to renew expired eligibilities.selfActivate: For users to activate their assignments.selfDeactivate: For users to deactivate their active assignments.selfExtend: For users to request to extend their expiring assignments.selfRenew: For users to request to renew their expired assignments.
func (m *UnifiedRoleEligibilityScheduleRequest) SetAction(value *UnifiedRoleScheduleRequestActions)() {
    err := m.GetBackingStore().Set("action", value)
    if err != nil {
        panic(err)
    }
}
// SetAppScope sets the appScope property value. Read-only property with details of the app-specific scope when the role eligibility is scoped to an app. Nullable. Supports $expand.
func (m *UnifiedRoleEligibilityScheduleRequest) SetAppScope(value AppScopeable)() {
    err := m.GetBackingStore().Set("appScope", value)
    if err != nil {
        panic(err)
    }
}
// SetAppScopeId sets the appScopeId property value. Identifier of the app-specific scope when the role eligibility is scoped to an app. The scope of a role eligibility determines the set of resources for which the principal is eligible to access. App scopes are scopes that are defined and understood by this application only. Use / for tenant-wide app scopes. Use directoryScopeId to limit the scope to particular directory objects, for example, administrative units. Supports $filter (eq, ne, and on null values).
func (m *UnifiedRoleEligibilityScheduleRequest) SetAppScopeId(value *string)() {
    err := m.GetBackingStore().Set("appScopeId", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectoryScope sets the directoryScope property value. The directory object that is the scope of the role eligibility. Read-only. Supports $expand.
func (m *UnifiedRoleEligibilityScheduleRequest) SetDirectoryScope(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("directoryScope", value)
    if err != nil {
        panic(err)
    }
}
// SetDirectoryScopeId sets the directoryScopeId property value. Identifier of the directory object representing the scope of the role eligibility. The scope of a role eligibility determines the set of resources for which the principal has been granted access. Directory scopes are shared scopes stored in the directory that are understood by multiple applications. Use / for tenant-wide scope. Use appScopeId to limit the scope to an application only. Supports $filter (eq, ne, and on null values).
func (m *UnifiedRoleEligibilityScheduleRequest) SetDirectoryScopeId(value *string)() {
    err := m.GetBackingStore().Set("directoryScopeId", value)
    if err != nil {
        panic(err)
    }
}
// SetIsValidationOnly sets the isValidationOnly property value. Determines whether the call is a validation or an actual call. Only set this property if you want to check whether an activation is subject to additional rules like MFA before actually submitting the request.
func (m *UnifiedRoleEligibilityScheduleRequest) SetIsValidationOnly(value *bool)() {
    err := m.GetBackingStore().Set("isValidationOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetJustification sets the justification property value. A message provided by users and administrators when create they create the unifiedRoleEligibilityScheduleRequest object.
func (m *UnifiedRoleEligibilityScheduleRequest) SetJustification(value *string)() {
    err := m.GetBackingStore().Set("justification", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipal sets the principal property value. The principal that's getting a role eligibility through the request. Supports $expand.
func (m *UnifiedRoleEligibilityScheduleRequest) SetPrincipal(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("principal", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalId sets the principalId property value. Identifier of the principal that has been granted the role eligibility. Can be a user or a role-assignable group. You can grant only active assignments service principals.Supports $filter (eq, ne).
func (m *UnifiedRoleEligibilityScheduleRequest) SetPrincipalId(value *string)() {
    err := m.GetBackingStore().Set("principalId", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleDefinition sets the roleDefinition property value. Detailed information for the unifiedRoleDefinition object that is referenced through the roleDefinitionId property. Supports $expand.
func (m *UnifiedRoleEligibilityScheduleRequest) SetRoleDefinition(value UnifiedRoleDefinitionable)() {
    err := m.GetBackingStore().Set("roleDefinition", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleDefinitionId sets the roleDefinitionId property value. Identifier of the unifiedRoleDefinition object that is being assigned to the principal. Supports $filter (eq, ne).
func (m *UnifiedRoleEligibilityScheduleRequest) SetRoleDefinitionId(value *string)() {
    err := m.GetBackingStore().Set("roleDefinitionId", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleInfo sets the scheduleInfo property value. The period of the role eligibility. Recurring schedules are currently unsupported.
func (m *UnifiedRoleEligibilityScheduleRequest) SetScheduleInfo(value RequestScheduleable)() {
    err := m.GetBackingStore().Set("scheduleInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetSchedule sets the targetSchedule property value. The schedule for a role eligibility that is referenced through the targetScheduleId property. Supports $expand.
func (m *UnifiedRoleEligibilityScheduleRequest) SetTargetSchedule(value UnifiedRoleEligibilityScheduleable)() {
    err := m.GetBackingStore().Set("targetSchedule", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetScheduleId sets the targetScheduleId property value. Identifier of the schedule object that's linked to the eligibility request. Supports $filter (eq, ne).
func (m *UnifiedRoleEligibilityScheduleRequest) SetTargetScheduleId(value *string)() {
    err := m.GetBackingStore().Set("targetScheduleId", value)
    if err != nil {
        panic(err)
    }
}
// SetTicketInfo sets the ticketInfo property value. Ticket details linked to the role eligibility request including details of the ticket number and ticket system. Optional.
func (m *UnifiedRoleEligibilityScheduleRequest) SetTicketInfo(value TicketInfoable)() {
    err := m.GetBackingStore().Set("ticketInfo", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleEligibilityScheduleRequestable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    Requestable
    GetAction()(*UnifiedRoleScheduleRequestActions)
    GetAppScope()(AppScopeable)
    GetAppScopeId()(*string)
    GetDirectoryScope()(DirectoryObjectable)
    GetDirectoryScopeId()(*string)
    GetIsValidationOnly()(*bool)
    GetJustification()(*string)
    GetPrincipal()(DirectoryObjectable)
    GetPrincipalId()(*string)
    GetRoleDefinition()(UnifiedRoleDefinitionable)
    GetRoleDefinitionId()(*string)
    GetScheduleInfo()(RequestScheduleable)
    GetTargetSchedule()(UnifiedRoleEligibilityScheduleable)
    GetTargetScheduleId()(*string)
    GetTicketInfo()(TicketInfoable)
    SetAction(value *UnifiedRoleScheduleRequestActions)()
    SetAppScope(value AppScopeable)()
    SetAppScopeId(value *string)()
    SetDirectoryScope(value DirectoryObjectable)()
    SetDirectoryScopeId(value *string)()
    SetIsValidationOnly(value *bool)()
    SetJustification(value *string)()
    SetPrincipal(value DirectoryObjectable)()
    SetPrincipalId(value *string)()
    SetRoleDefinition(value UnifiedRoleDefinitionable)()
    SetRoleDefinitionId(value *string)()
    SetScheduleInfo(value RequestScheduleable)()
    SetTargetSchedule(value UnifiedRoleEligibilityScheduleable)()
    SetTargetScheduleId(value *string)()
    SetTicketInfo(value TicketInfoable)()
}
