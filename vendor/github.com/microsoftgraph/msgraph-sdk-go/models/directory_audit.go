package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DirectoryAudit struct {
    Entity
}
// NewDirectoryAudit instantiates a new DirectoryAudit and sets the default values.
func NewDirectoryAudit()(*DirectoryAudit) {
    m := &DirectoryAudit{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDirectoryAuditFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectoryAuditFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDirectoryAudit(), nil
}
// GetActivityDateTime gets the activityDateTime property value. Indicates the date and time the activity was performed. The Timestamp type is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and $orderby.
// returns a *Time when successful
func (m *DirectoryAudit) GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("activityDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetActivityDisplayName gets the activityDisplayName property value. Indicates the activity name or the operation name (examples: 'Create User' and 'Add member to group'). For a list of activities logged, refer to Microsoft Entra audit log categories and activities. Supports $filter (eq, startswith).
// returns a *string when successful
func (m *DirectoryAudit) GetActivityDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("activityDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAdditionalDetails gets the additionalDetails property value. Indicates additional details on the activity.
// returns a []KeyValueable when successful
func (m *DirectoryAudit) GetAdditionalDetails()([]KeyValueable) {
    val, err := m.GetBackingStore().Get("additionalDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]KeyValueable)
    }
    return nil
}
// GetCategory gets the category property value. Indicates which resource category that's targeted by the activity. For example: UserManagement, GroupManagement, ApplicationManagement, RoleManagement. For a list of categories for activities logged, refer to Microsoft Entra audit log categories and activities.
// returns a *string when successful
func (m *DirectoryAudit) GetCategory()(*string) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCorrelationId gets the correlationId property value. Indicates a unique ID that helps correlate activities that span across various services. Can be used to trace logs across services. Supports $filter (eq).
// returns a *string when successful
func (m *DirectoryAudit) GetCorrelationId()(*string) {
    val, err := m.GetBackingStore().Get("correlationId")
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
func (m *DirectoryAudit) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activityDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityDateTime(val)
        }
        return nil
    }
    res["activityDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityDisplayName(val)
        }
        return nil
    }
    res["additionalDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateKeyValueFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]KeyValueable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(KeyValueable)
                }
            }
            m.SetAdditionalDetails(res)
        }
        return nil
    }
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val)
        }
        return nil
    }
    res["correlationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCorrelationId(val)
        }
        return nil
    }
    res["initiatedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuditActivityInitiatorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiatedBy(val.(AuditActivityInitiatorable))
        }
        return nil
    }
    res["loggedByService"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoggedByService(val)
        }
        return nil
    }
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val)
        }
        return nil
    }
    res["result"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseOperationResult)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResult(val.(*OperationResult))
        }
        return nil
    }
    res["resultReason"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResultReason(val)
        }
        return nil
    }
    res["targetResources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateTargetResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]TargetResourceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(TargetResourceable)
                }
            }
            m.SetTargetResources(res)
        }
        return nil
    }
    return res
}
// GetInitiatedBy gets the initiatedBy property value. The initiatedBy property
// returns a AuditActivityInitiatorable when successful
func (m *DirectoryAudit) GetInitiatedBy()(AuditActivityInitiatorable) {
    val, err := m.GetBackingStore().Get("initiatedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuditActivityInitiatorable)
    }
    return nil
}
// GetLoggedByService gets the loggedByService property value. Indicates information on which service initiated the activity (For example: Self-service Password Management, Core Directory, B2C, Invited Users, Microsoft Identity Manager, Privileged Identity Management. Supports $filter (eq).
// returns a *string when successful
func (m *DirectoryAudit) GetLoggedByService()(*string) {
    val, err := m.GetBackingStore().Get("loggedByService")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOperationType gets the operationType property value. Indicates the type of operation that was performed. The possible values include but are not limited to the following: Add, Assign, Update, Unassign, and Delete.
// returns a *string when successful
func (m *DirectoryAudit) GetOperationType()(*string) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResult gets the result property value. Indicates the result of the activity. Possible values are: success, failure, timeout, unknownFutureValue.
// returns a *OperationResult when successful
func (m *DirectoryAudit) GetResult()(*OperationResult) {
    val, err := m.GetBackingStore().Get("result")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*OperationResult)
    }
    return nil
}
// GetResultReason gets the resultReason property value. Indicates the reason for failure if the result is failure or timeout.
// returns a *string when successful
func (m *DirectoryAudit) GetResultReason()(*string) {
    val, err := m.GetBackingStore().Get("resultReason")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTargetResources gets the targetResources property value. Indicates information on which resource was changed due to the activity. Target Resource Type can be User, Device, Directory, App, Role, Group, Policy or Other. Supports $filter (eq) for id and displayName; and $filter (startswith) for displayName.
// returns a []TargetResourceable when successful
func (m *DirectoryAudit) GetTargetResources()([]TargetResourceable) {
    val, err := m.GetBackingStore().Get("targetResources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]TargetResourceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DirectoryAudit) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("activityDateTime", m.GetActivityDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("activityDisplayName", m.GetActivityDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetAdditionalDetails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAdditionalDetails()))
        for i, v := range m.GetAdditionalDetails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("additionalDetails", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("category", m.GetCategory())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("correlationId", m.GetCorrelationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiatedBy", m.GetInitiatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("loggedByService", m.GetLoggedByService())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("operationType", m.GetOperationType())
        if err != nil {
            return err
        }
    }
    if m.GetResult() != nil {
        cast := (*m.GetResult()).String()
        err = writer.WriteStringValue("result", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resultReason", m.GetResultReason())
        if err != nil {
            return err
        }
    }
    if m.GetTargetResources() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTargetResources()))
        for i, v := range m.GetTargetResources() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("targetResources", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivityDateTime sets the activityDateTime property value. Indicates the date and time the activity was performed. The Timestamp type is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Supports $filter (eq, ge, le) and $orderby.
func (m *DirectoryAudit) SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("activityDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetActivityDisplayName sets the activityDisplayName property value. Indicates the activity name or the operation name (examples: 'Create User' and 'Add member to group'). For a list of activities logged, refer to Microsoft Entra audit log categories and activities. Supports $filter (eq, startswith).
func (m *DirectoryAudit) SetActivityDisplayName(value *string)() {
    err := m.GetBackingStore().Set("activityDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalDetails sets the additionalDetails property value. Indicates additional details on the activity.
func (m *DirectoryAudit) SetAdditionalDetails(value []KeyValueable)() {
    err := m.GetBackingStore().Set("additionalDetails", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. Indicates which resource category that's targeted by the activity. For example: UserManagement, GroupManagement, ApplicationManagement, RoleManagement. For a list of categories for activities logged, refer to Microsoft Entra audit log categories and activities.
func (m *DirectoryAudit) SetCategory(value *string)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetCorrelationId sets the correlationId property value. Indicates a unique ID that helps correlate activities that span across various services. Can be used to trace logs across services. Supports $filter (eq).
func (m *DirectoryAudit) SetCorrelationId(value *string)() {
    err := m.GetBackingStore().Set("correlationId", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiatedBy sets the initiatedBy property value. The initiatedBy property
func (m *DirectoryAudit) SetInitiatedBy(value AuditActivityInitiatorable)() {
    err := m.GetBackingStore().Set("initiatedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLoggedByService sets the loggedByService property value. Indicates information on which service initiated the activity (For example: Self-service Password Management, Core Directory, B2C, Invited Users, Microsoft Identity Manager, Privileged Identity Management. Supports $filter (eq).
func (m *DirectoryAudit) SetLoggedByService(value *string)() {
    err := m.GetBackingStore().Set("loggedByService", value)
    if err != nil {
        panic(err)
    }
}
// SetOperationType sets the operationType property value. Indicates the type of operation that was performed. The possible values include but are not limited to the following: Add, Assign, Update, Unassign, and Delete.
func (m *DirectoryAudit) SetOperationType(value *string)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetResult sets the result property value. Indicates the result of the activity. Possible values are: success, failure, timeout, unknownFutureValue.
func (m *DirectoryAudit) SetResult(value *OperationResult)() {
    err := m.GetBackingStore().Set("result", value)
    if err != nil {
        panic(err)
    }
}
// SetResultReason sets the resultReason property value. Indicates the reason for failure if the result is failure or timeout.
func (m *DirectoryAudit) SetResultReason(value *string)() {
    err := m.GetBackingStore().Set("resultReason", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetResources sets the targetResources property value. Indicates information on which resource was changed due to the activity. Target Resource Type can be User, Device, Directory, App, Role, Group, Policy or Other. Supports $filter (eq) for id and displayName; and $filter (startswith) for displayName.
func (m *DirectoryAudit) SetTargetResources(value []TargetResourceable)() {
    err := m.GetBackingStore().Set("targetResources", value)
    if err != nil {
        panic(err)
    }
}
type DirectoryAuditable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivityDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetActivityDisplayName()(*string)
    GetAdditionalDetails()([]KeyValueable)
    GetCategory()(*string)
    GetCorrelationId()(*string)
    GetInitiatedBy()(AuditActivityInitiatorable)
    GetLoggedByService()(*string)
    GetOperationType()(*string)
    GetResult()(*OperationResult)
    GetResultReason()(*string)
    GetTargetResources()([]TargetResourceable)
    SetActivityDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetActivityDisplayName(value *string)()
    SetAdditionalDetails(value []KeyValueable)()
    SetCategory(value *string)()
    SetCorrelationId(value *string)()
    SetInitiatedBy(value AuditActivityInitiatorable)()
    SetLoggedByService(value *string)()
    SetOperationType(value *string)()
    SetResult(value *OperationResult)()
    SetResultReason(value *string)()
    SetTargetResources(value []TargetResourceable)()
}
