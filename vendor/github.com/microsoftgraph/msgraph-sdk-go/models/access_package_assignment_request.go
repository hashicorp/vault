package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessPackageAssignmentRequest struct {
    Entity
}
// NewAccessPackageAssignmentRequest instantiates a new AccessPackageAssignmentRequest and sets the default values.
func NewAccessPackageAssignmentRequest()(*AccessPackageAssignmentRequest) {
    m := &AccessPackageAssignmentRequest{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessPackageAssignmentRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentRequest(), nil
}
// GetAccessPackage gets the accessPackage property value. The access package associated with the accessPackageAssignmentRequest. An access package defines the collections of resource roles and the policies for how one or more users can get access to those resources. Read-only. Nullable.  Supports $expand.
// returns a AccessPackageable when successful
func (m *AccessPackageAssignmentRequest) GetAccessPackage()(AccessPackageable) {
    val, err := m.GetBackingStore().Get("accessPackage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageable)
    }
    return nil
}
// GetAnswers gets the answers property value. Answers provided by the requestor to accessPackageQuestions asked of them at the time of request.
// returns a []AccessPackageAnswerable when successful
func (m *AccessPackageAssignmentRequest) GetAnswers()([]AccessPackageAnswerable) {
    val, err := m.GetBackingStore().Get("answers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageAnswerable)
    }
    return nil
}
// GetAssignment gets the assignment property value. For a requestType of userAdd or adminAdd, this is an access package assignment requested to be created. For a requestType of userRemove, adminRemove or systemRemove, this has the id property of an existing assignment to be removed.   Supports $expand.
// returns a AccessPackageAssignmentable when successful
func (m *AccessPackageAssignmentRequest) GetAssignment()(AccessPackageAssignmentable) {
    val, err := m.GetBackingStore().Get("assignment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageAssignmentable)
    }
    return nil
}
// GetCompletedDateTime gets the completedDateTime property value. The date of the end of processing, either successful or failure, of a request. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
// returns a *Time when successful
func (m *AccessPackageAssignmentRequest) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports $filter.
// returns a *Time when successful
func (m *AccessPackageAssignmentRequest) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustomExtensionCalloutInstances gets the customExtensionCalloutInstances property value. Information about all the custom extension calls that were made during the access package assignment workflow.
// returns a []CustomExtensionCalloutInstanceable when successful
func (m *AccessPackageAssignmentRequest) GetCustomExtensionCalloutInstances()([]CustomExtensionCalloutInstanceable) {
    val, err := m.GetBackingStore().Get("customExtensionCalloutInstances")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CustomExtensionCalloutInstanceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAssignmentRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["answers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageAnswerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageAnswerable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageAnswerable)
                }
            }
            m.SetAnswers(res)
        }
        return nil
    }
    res["assignment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignment(val.(AccessPackageAssignmentable))
        }
        return nil
    }
    res["completedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletedDateTime(val)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
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
    res["requestor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessPackageSubjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestor(val.(AccessPackageSubjectable))
        }
        return nil
    }
    res["requestType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAccessPackageRequestType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestType(val.(*AccessPackageRequestType))
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
        val, err := n.GetEnumValue(ParseAccessPackageRequestState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*AccessPackageRequestState))
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
    return res
}
// GetRequestor gets the requestor property value. The subject who requested or, if a direct assignment, was assigned. Read-only. Nullable. Supports $expand.
// returns a AccessPackageSubjectable when successful
func (m *AccessPackageAssignmentRequest) GetRequestor()(AccessPackageSubjectable) {
    val, err := m.GetBackingStore().Get("requestor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessPackageSubjectable)
    }
    return nil
}
// GetRequestType gets the requestType property value. The type of the request. The possible values are: notSpecified, userAdd, UserExtend, userUpdate, userRemove, adminAdd, adminUpdate, adminRemove, systemAdd, systemUpdate, systemRemove, onBehalfAdd (not supported), unknownFutureValue. Requests from the user have a requestType of userAdd, userUpdate, or userRemove. This property can't be changed once set.
// returns a *AccessPackageRequestType when successful
func (m *AccessPackageAssignmentRequest) GetRequestType()(*AccessPackageRequestType) {
    val, err := m.GetBackingStore().Get("requestType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageRequestType)
    }
    return nil
}
// GetSchedule gets the schedule property value. The range of dates that access is to be assigned to the requestor. This property can't be changed once set, but a new schedule for an assignment can be included in another userUpdate or UserExtend or adminUpdate assignment request.
// returns a EntitlementManagementScheduleable when successful
func (m *AccessPackageAssignmentRequest) GetSchedule()(EntitlementManagementScheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementScheduleable)
    }
    return nil
}
// GetState gets the state property value. The state of the request. The possible values are: submitted, pendingApproval, delivering, delivered, deliveryFailed, denied, scheduled, canceled, partiallyDelivered, unknownFutureValue. Read-only. Supports $filter (eq).
// returns a *AccessPackageRequestState when successful
func (m *AccessPackageAssignmentRequest) GetState()(*AccessPackageRequestState) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AccessPackageRequestState)
    }
    return nil
}
// GetStatus gets the status property value. More information on the request processing status. Read-only.
// returns a *string when successful
func (m *AccessPackageAssignmentRequest) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    if m.GetAnswers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAnswers()))
        for i, v := range m.GetAnswers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("answers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("assignment", m.GetAssignment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
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
        err = writer.WriteObjectValue("requestor", m.GetRequestor())
        if err != nil {
            return err
        }
    }
    if m.GetRequestType() != nil {
        cast := (*m.GetRequestType()).String()
        err = writer.WriteStringValue("requestType", &cast)
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
    return nil
}
// SetAccessPackage sets the accessPackage property value. The access package associated with the accessPackageAssignmentRequest. An access package defines the collections of resource roles and the policies for how one or more users can get access to those resources. Read-only. Nullable.  Supports $expand.
func (m *AccessPackageAssignmentRequest) SetAccessPackage(value AccessPackageable)() {
    err := m.GetBackingStore().Set("accessPackage", value)
    if err != nil {
        panic(err)
    }
}
// SetAnswers sets the answers property value. Answers provided by the requestor to accessPackageQuestions asked of them at the time of request.
func (m *AccessPackageAssignmentRequest) SetAnswers(value []AccessPackageAnswerable)() {
    err := m.GetBackingStore().Set("answers", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignment sets the assignment property value. For a requestType of userAdd or adminAdd, this is an access package assignment requested to be created. For a requestType of userRemove, adminRemove or systemRemove, this has the id property of an existing assignment to be removed.   Supports $expand.
func (m *AccessPackageAssignmentRequest) SetAssignment(value AccessPackageAssignmentable)() {
    err := m.GetBackingStore().Set("assignment", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletedDateTime sets the completedDateTime property value. The date of the end of processing, either successful or failure, of a request. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
func (m *AccessPackageAssignmentRequest) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports $filter.
func (m *AccessPackageAssignmentRequest) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomExtensionCalloutInstances sets the customExtensionCalloutInstances property value. Information about all the custom extension calls that were made during the access package assignment workflow.
func (m *AccessPackageAssignmentRequest) SetCustomExtensionCalloutInstances(value []CustomExtensionCalloutInstanceable)() {
    err := m.GetBackingStore().Set("customExtensionCalloutInstances", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestor sets the requestor property value. The subject who requested or, if a direct assignment, was assigned. Read-only. Nullable. Supports $expand.
func (m *AccessPackageAssignmentRequest) SetRequestor(value AccessPackageSubjectable)() {
    err := m.GetBackingStore().Set("requestor", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestType sets the requestType property value. The type of the request. The possible values are: notSpecified, userAdd, UserExtend, userUpdate, userRemove, adminAdd, adminUpdate, adminRemove, systemAdd, systemUpdate, systemRemove, onBehalfAdd (not supported), unknownFutureValue. Requests from the user have a requestType of userAdd, userUpdate, or userRemove. This property can't be changed once set.
func (m *AccessPackageAssignmentRequest) SetRequestType(value *AccessPackageRequestType)() {
    err := m.GetBackingStore().Set("requestType", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. The range of dates that access is to be assigned to the requestor. This property can't be changed once set, but a new schedule for an assignment can be included in another userUpdate or UserExtend or adminUpdate assignment request.
func (m *AccessPackageAssignmentRequest) SetSchedule(value EntitlementManagementScheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state of the request. The possible values are: submitted, pendingApproval, delivering, delivered, deliveryFailed, denied, scheduled, canceled, partiallyDelivered, unknownFutureValue. Read-only. Supports $filter (eq).
func (m *AccessPackageAssignmentRequest) SetState(value *AccessPackageRequestState)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. More information on the request processing status. Read-only.
func (m *AccessPackageAssignmentRequest) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentRequestable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessPackage()(AccessPackageable)
    GetAnswers()([]AccessPackageAnswerable)
    GetAssignment()(AccessPackageAssignmentable)
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustomExtensionCalloutInstances()([]CustomExtensionCalloutInstanceable)
    GetRequestor()(AccessPackageSubjectable)
    GetRequestType()(*AccessPackageRequestType)
    GetSchedule()(EntitlementManagementScheduleable)
    GetState()(*AccessPackageRequestState)
    GetStatus()(*string)
    SetAccessPackage(value AccessPackageable)()
    SetAnswers(value []AccessPackageAnswerable)()
    SetAssignment(value AccessPackageAssignmentable)()
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustomExtensionCalloutInstances(value []CustomExtensionCalloutInstanceable)()
    SetRequestor(value AccessPackageSubjectable)()
    SetRequestType(value *AccessPackageRequestType)()
    SetSchedule(value EntitlementManagementScheduleable)()
    SetState(value *AccessPackageRequestState)()
    SetStatus(value *string)()
}
