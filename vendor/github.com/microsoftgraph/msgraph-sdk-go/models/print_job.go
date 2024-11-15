package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrintJob struct {
    Entity
}
// NewPrintJob instantiates a new PrintJob and sets the default values.
func NewPrintJob()(*PrintJob) {
    m := &PrintJob{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePrintJobFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrintJobFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrintJob(), nil
}
// GetConfiguration gets the configuration property value. The configuration property
// returns a PrintJobConfigurationable when successful
func (m *PrintJob) GetConfiguration()(PrintJobConfigurationable) {
    val, err := m.GetBackingStore().Get("configuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintJobConfigurationable)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. The createdBy property
// returns a UserIdentityable when successful
func (m *PrintJob) GetCreatedBy()(UserIdentityable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserIdentityable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The DateTimeOffset when the job was created. Read-only.
// returns a *Time when successful
func (m *PrintJob) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDocuments gets the documents property value. The documents property
// returns a []PrintDocumentable when successful
func (m *PrintJob) GetDocuments()([]PrintDocumentable) {
    val, err := m.GetBackingStore().Get("documents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintDocumentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrintJob) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintJobConfigurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfiguration(val.(PrintJobConfigurationable))
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(UserIdentityable))
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
    res["documents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintDocumentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintDocumentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintDocumentable)
                }
            }
            m.SetDocuments(res)
        }
        return nil
    }
    res["isFetchable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsFetchable(val)
        }
        return nil
    }
    res["redirectedFrom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRedirectedFrom(val)
        }
        return nil
    }
    res["redirectedTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRedirectedTo(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrintJobStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(PrintJobStatusable))
        }
        return nil
    }
    res["tasks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintTaskFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintTaskable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintTaskable)
                }
            }
            m.SetTasks(res)
        }
        return nil
    }
    return res
}
// GetIsFetchable gets the isFetchable property value. If true, document can be fetched by printer.
// returns a *bool when successful
func (m *PrintJob) GetIsFetchable()(*bool) {
    val, err := m.GetBackingStore().Get("isFetchable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRedirectedFrom gets the redirectedFrom property value. Contains the source job URL, if the job has been redirected from another printer.
// returns a *string when successful
func (m *PrintJob) GetRedirectedFrom()(*string) {
    val, err := m.GetBackingStore().Get("redirectedFrom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRedirectedTo gets the redirectedTo property value. Contains the destination job URL, if the job has been redirected to another printer.
// returns a *string when successful
func (m *PrintJob) GetRedirectedTo()(*string) {
    val, err := m.GetBackingStore().Get("redirectedTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStatus gets the status property value. The status property
// returns a PrintJobStatusable when successful
func (m *PrintJob) GetStatus()(PrintJobStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrintJobStatusable)
    }
    return nil
}
// GetTasks gets the tasks property value. A list of printTasks that were triggered by this print job.
// returns a []PrintTaskable when successful
func (m *PrintJob) GetTasks()([]PrintTaskable) {
    val, err := m.GetBackingStore().Get("tasks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintTaskable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrintJob) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("configuration", m.GetConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
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
    if m.GetDocuments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDocuments()))
        for i, v := range m.GetDocuments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("documents", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isFetchable", m.GetIsFetchable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("redirectedFrom", m.GetRedirectedFrom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("redirectedTo", m.GetRedirectedTo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    if m.GetTasks() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTasks()))
        for i, v := range m.GetTasks() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tasks", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfiguration sets the configuration property value. The configuration property
func (m *PrintJob) SetConfiguration(value PrintJobConfigurationable)() {
    err := m.GetBackingStore().Set("configuration", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. The createdBy property
func (m *PrintJob) SetCreatedBy(value UserIdentityable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The DateTimeOffset when the job was created. Read-only.
func (m *PrintJob) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDocuments sets the documents property value. The documents property
func (m *PrintJob) SetDocuments(value []PrintDocumentable)() {
    err := m.GetBackingStore().Set("documents", value)
    if err != nil {
        panic(err)
    }
}
// SetIsFetchable sets the isFetchable property value. If true, document can be fetched by printer.
func (m *PrintJob) SetIsFetchable(value *bool)() {
    err := m.GetBackingStore().Set("isFetchable", value)
    if err != nil {
        panic(err)
    }
}
// SetRedirectedFrom sets the redirectedFrom property value. Contains the source job URL, if the job has been redirected from another printer.
func (m *PrintJob) SetRedirectedFrom(value *string)() {
    err := m.GetBackingStore().Set("redirectedFrom", value)
    if err != nil {
        panic(err)
    }
}
// SetRedirectedTo sets the redirectedTo property value. Contains the destination job URL, if the job has been redirected to another printer.
func (m *PrintJob) SetRedirectedTo(value *string)() {
    err := m.GetBackingStore().Set("redirectedTo", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status property
func (m *PrintJob) SetStatus(value PrintJobStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTasks sets the tasks property value. A list of printTasks that were triggered by this print job.
func (m *PrintJob) SetTasks(value []PrintTaskable)() {
    err := m.GetBackingStore().Set("tasks", value)
    if err != nil {
        panic(err)
    }
}
type PrintJobable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConfiguration()(PrintJobConfigurationable)
    GetCreatedBy()(UserIdentityable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDocuments()([]PrintDocumentable)
    GetIsFetchable()(*bool)
    GetRedirectedFrom()(*string)
    GetRedirectedTo()(*string)
    GetStatus()(PrintJobStatusable)
    GetTasks()([]PrintTaskable)
    SetConfiguration(value PrintJobConfigurationable)()
    SetCreatedBy(value UserIdentityable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDocuments(value []PrintDocumentable)()
    SetIsFetchable(value *bool)()
    SetRedirectedFrom(value *string)()
    SetRedirectedTo(value *string)()
    SetStatus(value PrintJobStatusable)()
    SetTasks(value []PrintTaskable)()
}
