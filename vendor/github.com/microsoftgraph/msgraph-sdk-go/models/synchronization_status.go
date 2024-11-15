package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SynchronizationStatus struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSynchronizationStatus instantiates a new SynchronizationStatus and sets the default values.
func NewSynchronizationStatus()(*SynchronizationStatus) {
    m := &SynchronizationStatus{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSynchronizationStatusFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationStatusFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationStatus(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SynchronizationStatus) GetAdditionalData()(map[string]any) {
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
func (m *SynchronizationStatus) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCode gets the code property value. The code property
// returns a *SynchronizationStatusCode when successful
func (m *SynchronizationStatus) GetCode()(*SynchronizationStatusCode) {
    val, err := m.GetBackingStore().Get("code")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SynchronizationStatusCode)
    }
    return nil
}
// GetCountSuccessiveCompleteFailures gets the countSuccessiveCompleteFailures property value. Number of consecutive times this job failed.
// returns a *int64 when successful
func (m *SynchronizationStatus) GetCountSuccessiveCompleteFailures()(*int64) {
    val, err := m.GetBackingStore().Get("countSuccessiveCompleteFailures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetEscrowsPruned gets the escrowsPruned property value. true if the job's escrows (object-level errors) were pruned during initial synchronization. Escrows can be pruned if during the initial synchronization, you reach the threshold of errors that would normally put the job in quarantine. Instead of going into quarantine, the synchronization process clears the job's errors and continues until the initial synchronization is completed. When the initial synchronization is completed, the job will pause and wait for the customer to clean up the errors.
// returns a *bool when successful
func (m *SynchronizationStatus) GetEscrowsPruned()(*bool) {
    val, err := m.GetBackingStore().Get("escrowsPruned")
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
func (m *SynchronizationStatus) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["code"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSynchronizationStatusCode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCode(val.(*SynchronizationStatusCode))
        }
        return nil
    }
    res["countSuccessiveCompleteFailures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountSuccessiveCompleteFailures(val)
        }
        return nil
    }
    res["escrowsPruned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEscrowsPruned(val)
        }
        return nil
    }
    res["lastExecution"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationTaskExecutionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastExecution(val.(SynchronizationTaskExecutionable))
        }
        return nil
    }
    res["lastSuccessfulExecution"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationTaskExecutionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSuccessfulExecution(val.(SynchronizationTaskExecutionable))
        }
        return nil
    }
    res["lastSuccessfulExecutionWithExports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationTaskExecutionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSuccessfulExecutionWithExports(val.(SynchronizationTaskExecutionable))
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
    res["progress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSynchronizationProgressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SynchronizationProgressable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SynchronizationProgressable)
                }
            }
            m.SetProgress(res)
        }
        return nil
    }
    res["quarantine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationQuarantineFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuarantine(val.(SynchronizationQuarantineable))
        }
        return nil
    }
    res["steadyStateFirstAchievedTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSteadyStateFirstAchievedTime(val)
        }
        return nil
    }
    res["steadyStateLastAchievedTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSteadyStateLastAchievedTime(val)
        }
        return nil
    }
    res["synchronizedEntryCountByType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateStringKeyLongValuePairFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]StringKeyLongValuePairable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(StringKeyLongValuePairable)
                }
            }
            m.SetSynchronizedEntryCountByType(res)
        }
        return nil
    }
    res["troubleshootingUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTroubleshootingUrl(val)
        }
        return nil
    }
    return res
}
// GetLastExecution gets the lastExecution property value. Details of the last execution of the job.
// returns a SynchronizationTaskExecutionable when successful
func (m *SynchronizationStatus) GetLastExecution()(SynchronizationTaskExecutionable) {
    val, err := m.GetBackingStore().Get("lastExecution")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationTaskExecutionable)
    }
    return nil
}
// GetLastSuccessfulExecution gets the lastSuccessfulExecution property value. Details of the last execution of this job, which didn't have any errors.
// returns a SynchronizationTaskExecutionable when successful
func (m *SynchronizationStatus) GetLastSuccessfulExecution()(SynchronizationTaskExecutionable) {
    val, err := m.GetBackingStore().Get("lastSuccessfulExecution")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationTaskExecutionable)
    }
    return nil
}
// GetLastSuccessfulExecutionWithExports gets the lastSuccessfulExecutionWithExports property value. Details of the last execution of the job, which exported objects into the target directory.
// returns a SynchronizationTaskExecutionable when successful
func (m *SynchronizationStatus) GetLastSuccessfulExecutionWithExports()(SynchronizationTaskExecutionable) {
    val, err := m.GetBackingStore().Get("lastSuccessfulExecutionWithExports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationTaskExecutionable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SynchronizationStatus) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProgress gets the progress property value. Details of the progress of a job toward completion.
// returns a []SynchronizationProgressable when successful
func (m *SynchronizationStatus) GetProgress()([]SynchronizationProgressable) {
    val, err := m.GetBackingStore().Get("progress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SynchronizationProgressable)
    }
    return nil
}
// GetQuarantine gets the quarantine property value. If job is in quarantine, quarantine details.
// returns a SynchronizationQuarantineable when successful
func (m *SynchronizationStatus) GetQuarantine()(SynchronizationQuarantineable) {
    val, err := m.GetBackingStore().Get("quarantine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationQuarantineable)
    }
    return nil
}
// GetSteadyStateFirstAchievedTime gets the steadyStateFirstAchievedTime property value. The time when steady state (no more changes to the process) was first achieved. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *SynchronizationStatus) GetSteadyStateFirstAchievedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("steadyStateFirstAchievedTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSteadyStateLastAchievedTime gets the steadyStateLastAchievedTime property value. The time when steady state (no more changes to the process) was last achieved. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *SynchronizationStatus) GetSteadyStateLastAchievedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("steadyStateLastAchievedTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSynchronizedEntryCountByType gets the synchronizedEntryCountByType property value. Count of synchronized objects, listed by object type.
// returns a []StringKeyLongValuePairable when successful
func (m *SynchronizationStatus) GetSynchronizedEntryCountByType()([]StringKeyLongValuePairable) {
    val, err := m.GetBackingStore().Get("synchronizedEntryCountByType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]StringKeyLongValuePairable)
    }
    return nil
}
// GetTroubleshootingUrl gets the troubleshootingUrl property value. In the event of an error, the URL with the troubleshooting steps for the issue.
// returns a *string when successful
func (m *SynchronizationStatus) GetTroubleshootingUrl()(*string) {
    val, err := m.GetBackingStore().Get("troubleshootingUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationStatus) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetCode() != nil {
        cast := (*m.GetCode()).String()
        err := writer.WriteStringValue("code", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countSuccessiveCompleteFailures", m.GetCountSuccessiveCompleteFailures())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("escrowsPruned", m.GetEscrowsPruned())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("lastExecution", m.GetLastExecution())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("lastSuccessfulExecution", m.GetLastSuccessfulExecution())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("lastSuccessfulExecutionWithExports", m.GetLastSuccessfulExecutionWithExports())
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
    if m.GetProgress() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProgress()))
        for i, v := range m.GetProgress() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("progress", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("quarantine", m.GetQuarantine())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("steadyStateFirstAchievedTime", m.GetSteadyStateFirstAchievedTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("steadyStateLastAchievedTime", m.GetSteadyStateLastAchievedTime())
        if err != nil {
            return err
        }
    }
    if m.GetSynchronizedEntryCountByType() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSynchronizedEntryCountByType()))
        for i, v := range m.GetSynchronizedEntryCountByType() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("synchronizedEntryCountByType", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("troubleshootingUrl", m.GetTroubleshootingUrl())
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
func (m *SynchronizationStatus) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SynchronizationStatus) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCode sets the code property value. The code property
func (m *SynchronizationStatus) SetCode(value *SynchronizationStatusCode)() {
    err := m.GetBackingStore().Set("code", value)
    if err != nil {
        panic(err)
    }
}
// SetCountSuccessiveCompleteFailures sets the countSuccessiveCompleteFailures property value. Number of consecutive times this job failed.
func (m *SynchronizationStatus) SetCountSuccessiveCompleteFailures(value *int64)() {
    err := m.GetBackingStore().Set("countSuccessiveCompleteFailures", value)
    if err != nil {
        panic(err)
    }
}
// SetEscrowsPruned sets the escrowsPruned property value. true if the job's escrows (object-level errors) were pruned during initial synchronization. Escrows can be pruned if during the initial synchronization, you reach the threshold of errors that would normally put the job in quarantine. Instead of going into quarantine, the synchronization process clears the job's errors and continues until the initial synchronization is completed. When the initial synchronization is completed, the job will pause and wait for the customer to clean up the errors.
func (m *SynchronizationStatus) SetEscrowsPruned(value *bool)() {
    err := m.GetBackingStore().Set("escrowsPruned", value)
    if err != nil {
        panic(err)
    }
}
// SetLastExecution sets the lastExecution property value. Details of the last execution of the job.
func (m *SynchronizationStatus) SetLastExecution(value SynchronizationTaskExecutionable)() {
    err := m.GetBackingStore().Set("lastExecution", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSuccessfulExecution sets the lastSuccessfulExecution property value. Details of the last execution of this job, which didn't have any errors.
func (m *SynchronizationStatus) SetLastSuccessfulExecution(value SynchronizationTaskExecutionable)() {
    err := m.GetBackingStore().Set("lastSuccessfulExecution", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSuccessfulExecutionWithExports sets the lastSuccessfulExecutionWithExports property value. Details of the last execution of the job, which exported objects into the target directory.
func (m *SynchronizationStatus) SetLastSuccessfulExecutionWithExports(value SynchronizationTaskExecutionable)() {
    err := m.GetBackingStore().Set("lastSuccessfulExecutionWithExports", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SynchronizationStatus) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetProgress sets the progress property value. Details of the progress of a job toward completion.
func (m *SynchronizationStatus) SetProgress(value []SynchronizationProgressable)() {
    err := m.GetBackingStore().Set("progress", value)
    if err != nil {
        panic(err)
    }
}
// SetQuarantine sets the quarantine property value. If job is in quarantine, quarantine details.
func (m *SynchronizationStatus) SetQuarantine(value SynchronizationQuarantineable)() {
    err := m.GetBackingStore().Set("quarantine", value)
    if err != nil {
        panic(err)
    }
}
// SetSteadyStateFirstAchievedTime sets the steadyStateFirstAchievedTime property value. The time when steady state (no more changes to the process) was first achieved. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *SynchronizationStatus) SetSteadyStateFirstAchievedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("steadyStateFirstAchievedTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSteadyStateLastAchievedTime sets the steadyStateLastAchievedTime property value. The time when steady state (no more changes to the process) was last achieved. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *SynchronizationStatus) SetSteadyStateLastAchievedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("steadyStateLastAchievedTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSynchronizedEntryCountByType sets the synchronizedEntryCountByType property value. Count of synchronized objects, listed by object type.
func (m *SynchronizationStatus) SetSynchronizedEntryCountByType(value []StringKeyLongValuePairable)() {
    err := m.GetBackingStore().Set("synchronizedEntryCountByType", value)
    if err != nil {
        panic(err)
    }
}
// SetTroubleshootingUrl sets the troubleshootingUrl property value. In the event of an error, the URL with the troubleshooting steps for the issue.
func (m *SynchronizationStatus) SetTroubleshootingUrl(value *string)() {
    err := m.GetBackingStore().Set("troubleshootingUrl", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationStatusable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCode()(*SynchronizationStatusCode)
    GetCountSuccessiveCompleteFailures()(*int64)
    GetEscrowsPruned()(*bool)
    GetLastExecution()(SynchronizationTaskExecutionable)
    GetLastSuccessfulExecution()(SynchronizationTaskExecutionable)
    GetLastSuccessfulExecutionWithExports()(SynchronizationTaskExecutionable)
    GetOdataType()(*string)
    GetProgress()([]SynchronizationProgressable)
    GetQuarantine()(SynchronizationQuarantineable)
    GetSteadyStateFirstAchievedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSteadyStateLastAchievedTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSynchronizedEntryCountByType()([]StringKeyLongValuePairable)
    GetTroubleshootingUrl()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCode(value *SynchronizationStatusCode)()
    SetCountSuccessiveCompleteFailures(value *int64)()
    SetEscrowsPruned(value *bool)()
    SetLastExecution(value SynchronizationTaskExecutionable)()
    SetLastSuccessfulExecution(value SynchronizationTaskExecutionable)()
    SetLastSuccessfulExecutionWithExports(value SynchronizationTaskExecutionable)()
    SetOdataType(value *string)()
    SetProgress(value []SynchronizationProgressable)()
    SetQuarantine(value SynchronizationQuarantineable)()
    SetSteadyStateFirstAchievedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSteadyStateLastAchievedTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSynchronizedEntryCountByType(value []StringKeyLongValuePairable)()
    SetTroubleshootingUrl(value *string)()
}
