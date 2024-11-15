package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SynchronizationTaskExecution struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSynchronizationTaskExecution instantiates a new SynchronizationTaskExecution and sets the default values.
func NewSynchronizationTaskExecution()(*SynchronizationTaskExecution) {
    m := &SynchronizationTaskExecution{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSynchronizationTaskExecutionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSynchronizationTaskExecutionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSynchronizationTaskExecution(), nil
}
// GetActivityIdentifier gets the activityIdentifier property value. Identifier of the job run.
// returns a *string when successful
func (m *SynchronizationTaskExecution) GetActivityIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("activityIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SynchronizationTaskExecution) GetAdditionalData()(map[string]any) {
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
func (m *SynchronizationTaskExecution) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCountEntitled gets the countEntitled property value. Count of processed entries that were assigned for this application.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountEntitled()(*int64) {
    val, err := m.GetBackingStore().Get("countEntitled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountEntitledForProvisioning gets the countEntitledForProvisioning property value. Count of processed entries that were assigned for provisioning.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountEntitledForProvisioning()(*int64) {
    val, err := m.GetBackingStore().Get("countEntitledForProvisioning")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountEscrowed gets the countEscrowed property value. Count of entries that were escrowed (errors).
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountEscrowed()(*int64) {
    val, err := m.GetBackingStore().Get("countEscrowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountEscrowedRaw gets the countEscrowedRaw property value. Count of entries that were escrowed, including system-generated escrows.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountEscrowedRaw()(*int64) {
    val, err := m.GetBackingStore().Get("countEscrowedRaw")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountExported gets the countExported property value. Count of exported entries.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountExported()(*int64) {
    val, err := m.GetBackingStore().Get("countExported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountExports gets the countExports property value. Count of entries that were expected to be exported.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountExports()(*int64) {
    val, err := m.GetBackingStore().Get("countExports")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountImported gets the countImported property value. Count of imported entries.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountImported()(*int64) {
    val, err := m.GetBackingStore().Get("countImported")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountImportedDeltas gets the countImportedDeltas property value. Count of imported delta-changes.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountImportedDeltas()(*int64) {
    val, err := m.GetBackingStore().Get("countImportedDeltas")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetCountImportedReferenceDeltas gets the countImportedReferenceDeltas property value. Count of imported delta-changes pertaining to reference changes.
// returns a *int64 when successful
func (m *SynchronizationTaskExecution) GetCountImportedReferenceDeltas()(*int64) {
    val, err := m.GetBackingStore().Get("countImportedReferenceDeltas")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetError gets the error property value. If an error was encountered, contains a synchronizationError object with details.
// returns a SynchronizationErrorable when successful
func (m *SynchronizationTaskExecution) GetError()(SynchronizationErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SynchronizationErrorable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SynchronizationTaskExecution) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["activityIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivityIdentifier(val)
        }
        return nil
    }
    res["countEntitled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountEntitled(val)
        }
        return nil
    }
    res["countEntitledForProvisioning"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountEntitledForProvisioning(val)
        }
        return nil
    }
    res["countEscrowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountEscrowed(val)
        }
        return nil
    }
    res["countEscrowedRaw"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountEscrowedRaw(val)
        }
        return nil
    }
    res["countExported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountExported(val)
        }
        return nil
    }
    res["countExports"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountExports(val)
        }
        return nil
    }
    res["countImported"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountImported(val)
        }
        return nil
    }
    res["countImportedDeltas"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountImportedDeltas(val)
        }
        return nil
    }
    res["countImportedReferenceDeltas"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountImportedReferenceDeltas(val)
        }
        return nil
    }
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSynchronizationErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(SynchronizationErrorable))
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
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSynchronizationTaskExecutionResult)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val.(*SynchronizationTaskExecutionResult))
        }
        return nil
    }
    res["timeBegan"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeBegan(val)
        }
        return nil
    }
    res["timeEnded"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimeEnded(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SynchronizationTaskExecution) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetState gets the state property value. The state property
// returns a *SynchronizationTaskExecutionResult when successful
func (m *SynchronizationTaskExecution) GetState()(*SynchronizationTaskExecutionResult) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SynchronizationTaskExecutionResult)
    }
    return nil
}
// GetTimeBegan gets the timeBegan property value. Time when this job run began. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *SynchronizationTaskExecution) GetTimeBegan()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("timeBegan")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTimeEnded gets the timeEnded property value. Time when this job run ended. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *SynchronizationTaskExecution) GetTimeEnded()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("timeEnded")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SynchronizationTaskExecution) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("activityIdentifier", m.GetActivityIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countEntitled", m.GetCountEntitled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countEntitledForProvisioning", m.GetCountEntitledForProvisioning())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countEscrowed", m.GetCountEscrowed())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countEscrowedRaw", m.GetCountEscrowedRaw())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countExported", m.GetCountExported())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countExports", m.GetCountExports())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countImported", m.GetCountImported())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countImportedDeltas", m.GetCountImportedDeltas())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("countImportedReferenceDeltas", m.GetCountImportedReferenceDeltas())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("error", m.GetError())
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
    if m.GetState() != nil {
        cast := (*m.GetState()).String()
        err := writer.WriteStringValue("state", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("timeBegan", m.GetTimeBegan())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("timeEnded", m.GetTimeEnded())
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
// SetActivityIdentifier sets the activityIdentifier property value. Identifier of the job run.
func (m *SynchronizationTaskExecution) SetActivityIdentifier(value *string)() {
    err := m.GetBackingStore().Set("activityIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SynchronizationTaskExecution) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SynchronizationTaskExecution) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCountEntitled sets the countEntitled property value. Count of processed entries that were assigned for this application.
func (m *SynchronizationTaskExecution) SetCountEntitled(value *int64)() {
    err := m.GetBackingStore().Set("countEntitled", value)
    if err != nil {
        panic(err)
    }
}
// SetCountEntitledForProvisioning sets the countEntitledForProvisioning property value. Count of processed entries that were assigned for provisioning.
func (m *SynchronizationTaskExecution) SetCountEntitledForProvisioning(value *int64)() {
    err := m.GetBackingStore().Set("countEntitledForProvisioning", value)
    if err != nil {
        panic(err)
    }
}
// SetCountEscrowed sets the countEscrowed property value. Count of entries that were escrowed (errors).
func (m *SynchronizationTaskExecution) SetCountEscrowed(value *int64)() {
    err := m.GetBackingStore().Set("countEscrowed", value)
    if err != nil {
        panic(err)
    }
}
// SetCountEscrowedRaw sets the countEscrowedRaw property value. Count of entries that were escrowed, including system-generated escrows.
func (m *SynchronizationTaskExecution) SetCountEscrowedRaw(value *int64)() {
    err := m.GetBackingStore().Set("countEscrowedRaw", value)
    if err != nil {
        panic(err)
    }
}
// SetCountExported sets the countExported property value. Count of exported entries.
func (m *SynchronizationTaskExecution) SetCountExported(value *int64)() {
    err := m.GetBackingStore().Set("countExported", value)
    if err != nil {
        panic(err)
    }
}
// SetCountExports sets the countExports property value. Count of entries that were expected to be exported.
func (m *SynchronizationTaskExecution) SetCountExports(value *int64)() {
    err := m.GetBackingStore().Set("countExports", value)
    if err != nil {
        panic(err)
    }
}
// SetCountImported sets the countImported property value. Count of imported entries.
func (m *SynchronizationTaskExecution) SetCountImported(value *int64)() {
    err := m.GetBackingStore().Set("countImported", value)
    if err != nil {
        panic(err)
    }
}
// SetCountImportedDeltas sets the countImportedDeltas property value. Count of imported delta-changes.
func (m *SynchronizationTaskExecution) SetCountImportedDeltas(value *int64)() {
    err := m.GetBackingStore().Set("countImportedDeltas", value)
    if err != nil {
        panic(err)
    }
}
// SetCountImportedReferenceDeltas sets the countImportedReferenceDeltas property value. Count of imported delta-changes pertaining to reference changes.
func (m *SynchronizationTaskExecution) SetCountImportedReferenceDeltas(value *int64)() {
    err := m.GetBackingStore().Set("countImportedReferenceDeltas", value)
    if err != nil {
        panic(err)
    }
}
// SetError sets the error property value. If an error was encountered, contains a synchronizationError object with details.
func (m *SynchronizationTaskExecution) SetError(value SynchronizationErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SynchronizationTaskExecution) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. The state property
func (m *SynchronizationTaskExecution) SetState(value *SynchronizationTaskExecutionResult)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeBegan sets the timeBegan property value. Time when this job run began. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *SynchronizationTaskExecution) SetTimeBegan(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("timeBegan", value)
    if err != nil {
        panic(err)
    }
}
// SetTimeEnded sets the timeEnded property value. Time when this job run ended. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *SynchronizationTaskExecution) SetTimeEnded(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("timeEnded", value)
    if err != nil {
        panic(err)
    }
}
type SynchronizationTaskExecutionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivityIdentifier()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCountEntitled()(*int64)
    GetCountEntitledForProvisioning()(*int64)
    GetCountEscrowed()(*int64)
    GetCountEscrowedRaw()(*int64)
    GetCountExported()(*int64)
    GetCountExports()(*int64)
    GetCountImported()(*int64)
    GetCountImportedDeltas()(*int64)
    GetCountImportedReferenceDeltas()(*int64)
    GetError()(SynchronizationErrorable)
    GetOdataType()(*string)
    GetState()(*SynchronizationTaskExecutionResult)
    GetTimeBegan()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTimeEnded()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetActivityIdentifier(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCountEntitled(value *int64)()
    SetCountEntitledForProvisioning(value *int64)()
    SetCountEscrowed(value *int64)()
    SetCountEscrowedRaw(value *int64)()
    SetCountExported(value *int64)()
    SetCountExports(value *int64)()
    SetCountImported(value *int64)()
    SetCountImportedDeltas(value *int64)()
    SetCountImportedReferenceDeltas(value *int64)()
    SetError(value SynchronizationErrorable)()
    SetOdataType(value *string)()
    SetState(value *SynchronizationTaskExecutionResult)()
    SetTimeBegan(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTimeEnded(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
