package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TeamsAsyncOperation struct {
    Entity
}
// NewTeamsAsyncOperation instantiates a new TeamsAsyncOperation and sets the default values.
func NewTeamsAsyncOperation()(*TeamsAsyncOperation) {
    m := &TeamsAsyncOperation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateTeamsAsyncOperationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTeamsAsyncOperationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTeamsAsyncOperation(), nil
}
// GetAttemptsCount gets the attemptsCount property value. Number of times the operation was attempted before being marked successful or failed.
// returns a *int32 when successful
func (m *TeamsAsyncOperation) GetAttemptsCount()(*int32) {
    val, err := m.GetBackingStore().Get("attemptsCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Time when the operation was created.
// returns a *Time when successful
func (m *TeamsAsyncOperation) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetError gets the error property value. Any error that causes the async operation to fail.
// returns a OperationErrorable when successful
func (m *TeamsAsyncOperation) GetError()(OperationErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OperationErrorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TeamsAsyncOperation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attemptsCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttemptsCount(val)
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
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOperationErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(OperationErrorable))
        }
        return nil
    }
    res["lastActionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastActionDateTime(val)
        }
        return nil
    }
    res["operationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamsAsyncOperationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOperationType(val.(*TeamsAsyncOperationType))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamsAsyncOperationStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*TeamsAsyncOperationStatus))
        }
        return nil
    }
    res["targetResourceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetResourceId(val)
        }
        return nil
    }
    res["targetResourceLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetResourceLocation(val)
        }
        return nil
    }
    return res
}
// GetLastActionDateTime gets the lastActionDateTime property value. Time when the async operation was last updated.
// returns a *Time when successful
func (m *TeamsAsyncOperation) GetLastActionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastActionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetOperationType gets the operationType property value. The operationType property
// returns a *TeamsAsyncOperationType when successful
func (m *TeamsAsyncOperation) GetOperationType()(*TeamsAsyncOperationType) {
    val, err := m.GetBackingStore().Get("operationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamsAsyncOperationType)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a *TeamsAsyncOperationStatus when successful
func (m *TeamsAsyncOperation) GetStatus()(*TeamsAsyncOperationStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamsAsyncOperationStatus)
    }
    return nil
}
// GetTargetResourceId gets the targetResourceId property value. The ID of the object that's created or modified as result of this async operation, typically a team.
// returns a *string when successful
func (m *TeamsAsyncOperation) GetTargetResourceId()(*string) {
    val, err := m.GetBackingStore().Get("targetResourceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTargetResourceLocation gets the targetResourceLocation property value. The location of the object that's created or modified as result of this async operation. This URL should be treated as an opaque value and not parsed into its component paths.
// returns a *string when successful
func (m *TeamsAsyncOperation) GetTargetResourceLocation()(*string) {
    val, err := m.GetBackingStore().Get("targetResourceLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TeamsAsyncOperation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("attemptsCount", m.GetAttemptsCount())
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
    {
        err = writer.WriteObjectValue("error", m.GetError())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastActionDateTime", m.GetLastActionDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetOperationType() != nil {
        cast := (*m.GetOperationType()).String()
        err = writer.WriteStringValue("operationType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("targetResourceId", m.GetTargetResourceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("targetResourceLocation", m.GetTargetResourceLocation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttemptsCount sets the attemptsCount property value. Number of times the operation was attempted before being marked successful or failed.
func (m *TeamsAsyncOperation) SetAttemptsCount(value *int32)() {
    err := m.GetBackingStore().Set("attemptsCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Time when the operation was created.
func (m *TeamsAsyncOperation) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetError sets the error property value. Any error that causes the async operation to fail.
func (m *TeamsAsyncOperation) SetError(value OperationErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetLastActionDateTime sets the lastActionDateTime property value. Time when the async operation was last updated.
func (m *TeamsAsyncOperation) SetLastActionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastActionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetOperationType sets the operationType property value. The operationType property
func (m *TeamsAsyncOperation) SetOperationType(value *TeamsAsyncOperationType)() {
    err := m.GetBackingStore().Set("operationType", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *TeamsAsyncOperation) SetStatus(value *TeamsAsyncOperationStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetResourceId sets the targetResourceId property value. The ID of the object that's created or modified as result of this async operation, typically a team.
func (m *TeamsAsyncOperation) SetTargetResourceId(value *string)() {
    err := m.GetBackingStore().Set("targetResourceId", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetResourceLocation sets the targetResourceLocation property value. The location of the object that's created or modified as result of this async operation. This URL should be treated as an opaque value and not parsed into its component paths.
func (m *TeamsAsyncOperation) SetTargetResourceLocation(value *string)() {
    err := m.GetBackingStore().Set("targetResourceLocation", value)
    if err != nil {
        panic(err)
    }
}
type TeamsAsyncOperationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttemptsCount()(*int32)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetError()(OperationErrorable)
    GetLastActionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetOperationType()(*TeamsAsyncOperationType)
    GetStatus()(*TeamsAsyncOperationStatus)
    GetTargetResourceId()(*string)
    GetTargetResourceLocation()(*string)
    SetAttemptsCount(value *int32)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetError(value OperationErrorable)()
    SetLastActionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetOperationType(value *TeamsAsyncOperationType)()
    SetStatus(value *TeamsAsyncOperationStatus)()
    SetTargetResourceId(value *string)()
    SetTargetResourceLocation(value *string)()
}
