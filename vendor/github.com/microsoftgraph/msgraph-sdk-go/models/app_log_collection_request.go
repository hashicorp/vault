package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AppLogCollectionRequest entity for AppLogCollectionRequest contains all logs values.
type AppLogCollectionRequest struct {
    Entity
}
// NewAppLogCollectionRequest instantiates a new AppLogCollectionRequest and sets the default values.
func NewAppLogCollectionRequest()(*AppLogCollectionRequest) {
    m := &AppLogCollectionRequest{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAppLogCollectionRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppLogCollectionRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppLogCollectionRequest(), nil
}
// GetCompletedDateTime gets the completedDateTime property value. Time at which the upload log request reached a completed state if not completed yet NULL will be returned.
// returns a *Time when successful
func (m *AppLogCollectionRequest) GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustomLogFolders gets the customLogFolders property value. List of log folders.
// returns a []string when successful
func (m *AppLogCollectionRequest) GetCustomLogFolders()([]string) {
    val, err := m.GetBackingStore().Get("customLogFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetErrorMessage gets the errorMessage property value. Indicates error message if any during the upload process.
// returns a *string when successful
func (m *AppLogCollectionRequest) GetErrorMessage()(*string) {
    val, err := m.GetBackingStore().Get("errorMessage")
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
func (m *AppLogCollectionRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["customLogFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetCustomLogFolders(res)
        }
        return nil
    }
    res["errorMessage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetErrorMessage(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAppLogUploadState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*AppLogUploadState))
        }
        return nil
    }
    return res
}
// GetStatus gets the status property value. AppLogUploadStatus
// returns a *AppLogUploadState when successful
func (m *AppLogCollectionRequest) GetStatus()(*AppLogUploadState) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AppLogUploadState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AppLogCollectionRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("completedDateTime", m.GetCompletedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetCustomLogFolders() != nil {
        err = writer.WriteCollectionOfStringValues("customLogFolders", m.GetCustomLogFolders())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("errorMessage", m.GetErrorMessage())
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
    return nil
}
// SetCompletedDateTime sets the completedDateTime property value. Time at which the upload log request reached a completed state if not completed yet NULL will be returned.
func (m *AppLogCollectionRequest) SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomLogFolders sets the customLogFolders property value. List of log folders.
func (m *AppLogCollectionRequest) SetCustomLogFolders(value []string)() {
    err := m.GetBackingStore().Set("customLogFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorMessage sets the errorMessage property value. Indicates error message if any during the upload process.
func (m *AppLogCollectionRequest) SetErrorMessage(value *string)() {
    err := m.GetBackingStore().Set("errorMessage", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. AppLogUploadStatus
func (m *AppLogCollectionRequest) SetStatus(value *AppLogUploadState)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type AppLogCollectionRequestable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustomLogFolders()([]string)
    GetErrorMessage()(*string)
    GetStatus()(*AppLogUploadState)
    SetCompletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustomLogFolders(value []string)()
    SetErrorMessage(value *string)()
    SetStatus(value *AppLogUploadState)()
}
