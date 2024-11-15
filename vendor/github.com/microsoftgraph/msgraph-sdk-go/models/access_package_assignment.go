package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageAssignment struct {
    Entity
}
// NewAccessPackageAssignment instantiates a new AccessPackageAssignment and sets the default values.
func NewAccessPackageAssignment()(*AccessPackageAssignment) {
    m := &AccessPackageAssignment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignment(), nil
}
// GetAccessPackage gets the accessPackage property value. Read-only. Nullable. Supports $filter (eq) on the id property and $expand query parameters.
// returns a AccessPackageable when successful
func (m *AccessPackageAssignment) GetAccessPackage()(AccessPackageable) {
    val, err := m.GetBackingStore().Get("accessPackage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageable)
    }
    return nil
}
// GetAssignmentPolicy gets the assignmentPolicy property value. Read-only. Supports $filter (eq) on the id property and $expand query parameters.
// returns a AccessPackageAssignmentPolicyable when successful
func (m *AccessPackageAssignment) GetAssignmentPolicy()(AccessPackageAssignmentPolicyable) {
    val, err := m.GetBackingStore().Get("assignmentPolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAssignmentPolicyable)
    }
    return nil
}
// GetCustomExtensionCalloutInstances gets the customExtensionCalloutInstances property value. Information about all the custom extension calls that were made during the access package assignment workflow.
// returns a []CustomExtensionCalloutInstanceable when successful
func (m *AccessPackageAssignment) GetCustomExtensionCalloutInstances()([]CustomExtensionCalloutInstanceable) {
    val, err := m.GetBackingStore().Get("customExtensionCalloutInstances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CustomExtensionCalloutInstanceable)
    }
    return nil
}
// GetExpiredDateTime gets the expiredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *AccessPackageAssignment) GetExpiredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("expiredDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessPackage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessPackage(val.(AccessPackageable))
        }
        return nil
    }
    res["assignmentPolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAssignmentPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentPolicy(val.(AccessPackageAssignmentPolicyable))
        }
        return nil
    }
    res["customExtensionCalloutInstances"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCustomExtensionCalloutInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CustomExtensionCalloutInstanceable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CustomExtensionCalloutInstanceable)
                }
            }
            m.SetCustomExtensionCalloutInstances(res)
        }
        return nil
    }
    res["expiredDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpiredDateTime(val)
        }
        return nil
    }
    res["schedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntitlementManagementScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedule(val.(EntitlementManagementScheduleable))
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessPackageAssignmentState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*AccessPackageAssignmentState))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val)
        }
        return nil
    }
    res["target"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageSubjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTarget(val.(AccessPackageSubjectable))
        }
        return nil
    }
    return res
}
// GetSchedule gets the schedule property value. When the access assignment is to be in place. Read-only.
// returns a EntitlementManagementScheduleable when successful
func (m *AccessPackageAssignment) GetSchedule()(EntitlementManagementScheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementScheduleable)
    }
    return nil
}
// GetState gets the state property value. The state of the access package assignment. The possible values are: delivering, partiallyDelivered, delivered, expired, deliveryFailed, unknownFutureValue. Read-only. Supports $filter (eq).
// returns a *AccessPackageAssignmentState when successful
func (m *AccessPackageAssignment) GetState()(*AccessPackageAssignmentState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageAssignmentState)
    }
    return nil
}
// GetStatus gets the status property value. More information about the assignment lifecycle. Possible values include Delivering, Delivered, NearExpiry1DayNotificationTriggered, or ExpiredNotificationTriggered. Read-only.
// returns a *string when successful
func (m *AccessPackageAssignment) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTarget gets the target property value. The subject of the access package assignment. Read-only. Nullable. Supports $expand. Supports $filter (eq) on objectId.
// returns a AccessPackageSubjectable when successful
func (m *AccessPackageAssignment) GetTarget()(AccessPackageSubjectable) {
    val, err := m.GetBackingStore().Get("target")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageSubjectable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("accessPackage", m.GetAccessPackage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignmentPolicy", m.GetAssignmentPolicy())
        if err != nil {
            return err
        }
    }
    if m.GetCustomExtensionCalloutInstances() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCustomExtensionCalloutInstances()))
        for i, v := range m.GetCustomExtensionCalloutInstances() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("customExtensionCalloutInstances", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("expiredDateTime", m.GetExpiredDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("schedule", m.GetSchedule())
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
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("target", m.GetTarget())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessPackage sets the accessPackage property value. Read-only. Nullable. Supports $filter (eq) on the id property and $expand query parameters.
func (m *AccessPackageAssignment) SetAccessPackage(value AccessPackageable)() {
    err := m.GetBackingStore().Set("accessPackage", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentPolicy sets the assignmentPolicy property value. Read-only. Supports $filter (eq) on the id property and $expand query parameters.
func (m *AccessPackageAssignment) SetAssignmentPolicy(value AccessPackageAssignmentPolicyable)() {
    err := m.GetBackingStore().Set("assignmentPolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomExtensionCalloutInstances sets the customExtensionCalloutInstances property value. Information about all the custom extension calls that were made during the access package assignment workflow.
func (m *AccessPackageAssignment) SetCustomExtensionCalloutInstances(value []CustomExtensionCalloutInstanceable)() {
    err := m.GetBackingStore().Set("customExtensionCalloutInstances", value)
    if err != nil {
        panic(err)
    }
}
// SetExpiredDateTime sets the expiredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AccessPackageAssignment) SetExpiredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expiredDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. When the access assignment is to be in place. Read-only.
func (m *AccessPackageAssignment) SetSchedule(value EntitlementManagementScheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state of the access package assignment. The possible values are: delivering, partiallyDelivered, delivered, expired, deliveryFailed, unknownFutureValue. Read-only. Supports $filter (eq).
func (m *AccessPackageAssignment) SetState(value *AccessPackageAssignmentState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. More information about the assignment lifecycle. Possible values include Delivering, Delivered, NearExpiry1DayNotificationTriggered, or ExpiredNotificationTriggered. Read-only.
func (m *AccessPackageAssignment) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTarget sets the target property value. The subject of the access package assignment. Read-only. Nullable. Supports $expand. Supports $filter (eq) on objectId.
func (m *AccessPackageAssignment) SetTarget(value AccessPackageSubjectable)() {
    err := m.GetBackingStore().Set("target", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackage()(AccessPackageable)
    GetAssignmentPolicy()(AccessPackageAssignmentPolicyable)
    GetCustomExtensionCalloutInstances()([]CustomExtensionCalloutInstanceable)
    GetExpiredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSchedule()(EntitlementManagementScheduleable)
    GetState()(*AccessPackageAssignmentState)
    GetStatus()(*string)
    GetTarget()(AccessPackageSubjectable)
    SetAccessPackage(value AccessPackageable)()
    SetAssignmentPolicy(value AccessPackageAssignmentPolicyable)()
    SetCustomExtensionCalloutInstances(value []CustomExtensionCalloutInstanceable)()
    SetExpiredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSchedule(value EntitlementManagementScheduleable)()
    SetState(value *AccessPackageAssignmentState)()
    SetStatus(value *string)()
    SetTarget(value AccessPackageSubjectable)()
}
