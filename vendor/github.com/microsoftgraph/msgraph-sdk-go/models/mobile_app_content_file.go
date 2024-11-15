package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MobileAppContentFile contains properties for a single installer file that is associated with a given mobileAppContent version.
type MobileAppContentFile struct {
    Entity
}
// NewMobileAppContentFile instantiates a new MobileAppContentFile and sets the default values.
func NewMobileAppContentFile()(*MobileAppContentFile) {
    m := &MobileAppContentFile{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMobileAppContentFileFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMobileAppContentFileFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMobileAppContentFile(), nil
}
// GetAzureStorageUri gets the azureStorageUri property value. The Azure Storage URI.
// returns a *string when successful
func (m *MobileAppContentFile) GetAzureStorageUri()(*string) {
    val, err := m.GetBackingStore().Get("azureStorageUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAzureStorageUriExpirationDateTime gets the azureStorageUriExpirationDateTime property value. The time the Azure storage Uri expires.
// returns a *Time when successful
func (m *MobileAppContentFile) GetAzureStorageUriExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("azureStorageUriExpirationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The time the file was created.
// returns a *Time when successful
func (m *MobileAppContentFile) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *MobileAppContentFile) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["azureStorageUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureStorageUri(val)
        }
        return nil
    }
    res["azureStorageUriExpirationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureStorageUriExpirationDateTime(val)
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
    res["isCommitted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsCommitted(val)
        }
        return nil
    }
    res["isDependency"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDependency(val)
        }
        return nil
    }
    res["manifest"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManifest(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["size"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSize(val)
        }
        return nil
    }
    res["sizeEncrypted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSizeEncrypted(val)
        }
        return nil
    }
    res["uploadState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMobileAppContentFileUploadState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUploadState(val.(*MobileAppContentFileUploadState))
        }
        return nil
    }
    return res
}
// GetIsCommitted gets the isCommitted property value. A value indicating whether the file is committed.
// returns a *bool when successful
func (m *MobileAppContentFile) GetIsCommitted()(*bool) {
    val, err := m.GetBackingStore().Get("isCommitted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDependency gets the isDependency property value. Indicates whether this content file is a dependency for the main content file. TRUE means that the content file is a dependency, FALSE means that the content file is not a dependency and is the main content file. Defaults to FALSE.
// returns a *bool when successful
func (m *MobileAppContentFile) GetIsDependency()(*bool) {
    val, err := m.GetBackingStore().Get("isDependency")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetManifest gets the manifest property value. The manifest information.
// returns a []byte when successful
func (m *MobileAppContentFile) GetManifest()([]byte) {
    val, err := m.GetBackingStore().Get("manifest")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetName gets the name property value. the file name.
// returns a *string when successful
func (m *MobileAppContentFile) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSize gets the size property value. The size of the file prior to encryption.
// returns a *int64 when successful
func (m *MobileAppContentFile) GetSize()(*int64) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetSizeEncrypted gets the sizeEncrypted property value. The size of the file after encryption.
// returns a *int64 when successful
func (m *MobileAppContentFile) GetSizeEncrypted()(*int64) {
    val, err := m.GetBackingStore().Get("sizeEncrypted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUploadState gets the uploadState property value. Contains properties for upload request states.
// returns a *MobileAppContentFileUploadState when successful
func (m *MobileAppContentFile) GetUploadState()(*MobileAppContentFileUploadState) {
    val, err := m.GetBackingStore().Get("uploadState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MobileAppContentFileUploadState)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MobileAppContentFile) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isDependency", m.GetIsDependency())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("manifest", m.GetManifest())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("sizeEncrypted", m.GetSizeEncrypted())
        if err != nil {
            return err
        }
    }
    if m.GetUploadState() != nil {
        cast := (*m.GetUploadState()).String()
        err = writer.WriteStringValue("uploadState", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureStorageUri sets the azureStorageUri property value. The Azure Storage URI.
func (m *MobileAppContentFile) SetAzureStorageUri(value *string)() {
    err := m.GetBackingStore().Set("azureStorageUri", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureStorageUriExpirationDateTime sets the azureStorageUriExpirationDateTime property value. The time the Azure storage Uri expires.
func (m *MobileAppContentFile) SetAzureStorageUriExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("azureStorageUriExpirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The time the file was created.
func (m *MobileAppContentFile) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetIsCommitted sets the isCommitted property value. A value indicating whether the file is committed.
func (m *MobileAppContentFile) SetIsCommitted(value *bool)() {
    err := m.GetBackingStore().Set("isCommitted", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDependency sets the isDependency property value. Indicates whether this content file is a dependency for the main content file. TRUE means that the content file is a dependency, FALSE means that the content file is not a dependency and is the main content file. Defaults to FALSE.
func (m *MobileAppContentFile) SetIsDependency(value *bool)() {
    err := m.GetBackingStore().Set("isDependency", value)
    if err != nil {
        panic(err)
    }
}
// SetManifest sets the manifest property value. The manifest information.
func (m *MobileAppContentFile) SetManifest(value []byte)() {
    err := m.GetBackingStore().Set("manifest", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. the file name.
func (m *MobileAppContentFile) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. The size of the file prior to encryption.
func (m *MobileAppContentFile) SetSize(value *int64)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetSizeEncrypted sets the sizeEncrypted property value. The size of the file after encryption.
func (m *MobileAppContentFile) SetSizeEncrypted(value *int64)() {
    err := m.GetBackingStore().Set("sizeEncrypted", value)
    if err != nil {
        panic(err)
    }
}
// SetUploadState sets the uploadState property value. Contains properties for upload request states.
func (m *MobileAppContentFile) SetUploadState(value *MobileAppContentFileUploadState)() {
    err := m.GetBackingStore().Set("uploadState", value)
    if err != nil {
        panic(err)
    }
}
type MobileAppContentFileable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureStorageUri()(*string)
    GetAzureStorageUriExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetIsCommitted()(*bool)
    GetIsDependency()(*bool)
    GetManifest()([]byte)
    GetName()(*string)
    GetSize()(*int64)
    GetSizeEncrypted()(*int64)
    GetUploadState()(*MobileAppContentFileUploadState)
    SetAzureStorageUri(value *string)()
    SetAzureStorageUriExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetIsCommitted(value *bool)()
    SetIsDependency(value *bool)()
    SetManifest(value []byte)()
    SetName(value *string)()
    SetSize(value *int64)()
    SetSizeEncrypted(value *int64)()
    SetUploadState(value *MobileAppContentFileUploadState)()
}
