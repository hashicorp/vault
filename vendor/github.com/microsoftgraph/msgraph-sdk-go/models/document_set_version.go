package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DocumentSetVersion struct {
    ListItemVersion
}
// NewDocumentSetVersion instantiates a new DocumentSetVersion and sets the default values.
func NewDocumentSetVersion()(*DocumentSetVersion) {
    m := &DocumentSetVersion{
        ListItemVersion: *NewListItemVersion(),
    }
    odataTypeValue := "#microsoft.graph.documentSetVersion"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDocumentSetVersionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDocumentSetVersionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDocumentSetVersion(), nil
}
// GetComment gets the comment property value. Comment about the captured version.
// returns a *string when successful
func (m *DocumentSetVersion) GetComment()(*string) {
    val, err := m.GetBackingStore().Get("comment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. User who captured the version.
// returns a IdentitySetable when successful
func (m *DocumentSetVersion) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time when this version was created.
// returns a *Time when successful
func (m *DocumentSetVersion) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
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
func (m *DocumentSetVersion) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ListItemVersion.GetFieldDeserializers()
    res["comment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComment(val)
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(IdentitySetable))
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
    res["items"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDocumentSetVersionItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DocumentSetVersionItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DocumentSetVersionItemable)
                }
            }
            m.SetItems(res)
        }
        return nil
    }
    res["shouldCaptureMinorVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShouldCaptureMinorVersion(val)
        }
        return nil
    }
    return res
}
// GetItems gets the items property value. Items within the document set that are captured as part of this version.
// returns a []DocumentSetVersionItemable when successful
func (m *DocumentSetVersion) GetItems()([]DocumentSetVersionItemable) {
    val, err := m.GetBackingStore().Get("items")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DocumentSetVersionItemable)
    }
    return nil
}
// GetShouldCaptureMinorVersion gets the shouldCaptureMinorVersion property value. If true, minor versions of items are also captured; otherwise, only major versions are captured. The default value is false.
// returns a *bool when successful
func (m *DocumentSetVersion) GetShouldCaptureMinorVersion()(*bool) {
    val, err := m.GetBackingStore().Get("shouldCaptureMinorVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DocumentSetVersion) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ListItemVersion.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("comment", m.GetComment())
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
    if m.GetItems() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetItems()))
        for i, v := range m.GetItems() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("items", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("shouldCaptureMinorVersion", m.GetShouldCaptureMinorVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetComment sets the comment property value. Comment about the captured version.
func (m *DocumentSetVersion) SetComment(value *string)() {
    err := m.GetBackingStore().Set("comment", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. User who captured the version.
func (m *DocumentSetVersion) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time when this version was created.
func (m *DocumentSetVersion) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetItems sets the items property value. Items within the document set that are captured as part of this version.
func (m *DocumentSetVersion) SetItems(value []DocumentSetVersionItemable)() {
    err := m.GetBackingStore().Set("items", value)
    if err != nil {
        panic(err)
    }
}
// SetShouldCaptureMinorVersion sets the shouldCaptureMinorVersion property value. If true, minor versions of items are also captured; otherwise, only major versions are captured. The default value is false.
func (m *DocumentSetVersion) SetShouldCaptureMinorVersion(value *bool)() {
    err := m.GetBackingStore().Set("shouldCaptureMinorVersion", value)
    if err != nil {
        panic(err)
    }
}
type DocumentSetVersionable interface {
    ListItemVersionable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetComment()(*string)
    GetCreatedBy()(IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetItems()([]DocumentSetVersionItemable)
    GetShouldCaptureMinorVersion()(*bool)
    SetComment(value *string)()
    SetCreatedBy(value IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetItems(value []DocumentSetVersionItemable)()
    SetShouldCaptureMinorVersion(value *bool)()
}
