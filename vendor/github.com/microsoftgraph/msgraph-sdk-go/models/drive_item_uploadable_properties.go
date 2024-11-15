package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type DriveItemUploadableProperties struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDriveItemUploadableProperties instantiates a new DriveItemUploadableProperties and sets the default values.
func NewDriveItemUploadableProperties()(*DriveItemUploadableProperties) {
    m := &DriveItemUploadableProperties{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDriveItemUploadablePropertiesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDriveItemUploadablePropertiesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDriveItemUploadableProperties(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DriveItemUploadableProperties) GetAdditionalData()(map[string]any) {
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
func (m *DriveItemUploadableProperties) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDescription gets the description property value. Provides a user-visible description of the item. Read-write. Only on OneDrive Personal.
// returns a *string when successful
func (m *DriveItemUploadableProperties) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDriveItemSource gets the driveItemSource property value. Information about the drive item source. Read-write. Only on OneDrive for Business and SharePoint.
// returns a DriveItemSourceable when successful
func (m *DriveItemUploadableProperties) GetDriveItemSource()(DriveItemSourceable) {
    val, err := m.GetBackingStore().Get("driveItemSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DriveItemSourceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DriveItemUploadableProperties) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["driveItemSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveItemSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDriveItemSource(val.(DriveItemSourceable))
        }
        return nil
    }
    res["fileSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileSize(val)
        }
        return nil
    }
    res["fileSystemInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileSystemInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFileSystemInfo(val.(FileSystemInfoable))
        }
        return nil
    }
    res["mediaSource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaSourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaSource(val.(MediaSourceable))
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
    return res
}
// GetFileSize gets the fileSize property value. Provides an expected file size to perform a quota check before uploading. Only on OneDrive Personal.
// returns a *int64 when successful
func (m *DriveItemUploadableProperties) GetFileSize()(*int64) {
    val, err := m.GetBackingStore().Get("fileSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFileSystemInfo gets the fileSystemInfo property value. File system information on client. Read-write.
// returns a FileSystemInfoable when successful
func (m *DriveItemUploadableProperties) GetFileSystemInfo()(FileSystemInfoable) {
    val, err := m.GetBackingStore().Get("fileSystemInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileSystemInfoable)
    }
    return nil
}
// GetMediaSource gets the mediaSource property value. Media source information. Read-write. Only on OneDrive for Business and SharePoint.
// returns a MediaSourceable when successful
func (m *DriveItemUploadableProperties) GetMediaSource()(MediaSourceable) {
    val, err := m.GetBackingStore().Get("mediaSource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaSourceable)
    }
    return nil
}
// GetName gets the name property value. The name of the item (filename and extension). Read-write.
// returns a *string when successful
func (m *DriveItemUploadableProperties) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DriveItemUploadableProperties) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DriveItemUploadableProperties) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("driveItemSource", m.GetDriveItemSource())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("fileSize", m.GetFileSize())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("fileSystemInfo", m.GetFileSystemInfo())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("mediaSource", m.GetMediaSource())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("name", m.GetName())
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
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *DriveItemUploadableProperties) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DriveItemUploadableProperties) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDescription sets the description property value. Provides a user-visible description of the item. Read-write. Only on OneDrive Personal.
func (m *DriveItemUploadableProperties) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDriveItemSource sets the driveItemSource property value. Information about the drive item source. Read-write. Only on OneDrive for Business and SharePoint.
func (m *DriveItemUploadableProperties) SetDriveItemSource(value DriveItemSourceable)() {
    err := m.GetBackingStore().Set("driveItemSource", value)
    if err != nil {
        panic(err)
    }
}
// SetFileSize sets the fileSize property value. Provides an expected file size to perform a quota check before uploading. Only on OneDrive Personal.
func (m *DriveItemUploadableProperties) SetFileSize(value *int64)() {
    err := m.GetBackingStore().Set("fileSize", value)
    if err != nil {
        panic(err)
    }
}
// SetFileSystemInfo sets the fileSystemInfo property value. File system information on client. Read-write.
func (m *DriveItemUploadableProperties) SetFileSystemInfo(value FileSystemInfoable)() {
    err := m.GetBackingStore().Set("fileSystemInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaSource sets the mediaSource property value. Media source information. Read-write. Only on OneDrive for Business and SharePoint.
func (m *DriveItemUploadableProperties) SetMediaSource(value MediaSourceable)() {
    err := m.GetBackingStore().Set("mediaSource", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name of the item (filename and extension). Read-write.
func (m *DriveItemUploadableProperties) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DriveItemUploadableProperties) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type DriveItemUploadablePropertiesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDescription()(*string)
    GetDriveItemSource()(DriveItemSourceable)
    GetFileSize()(*int64)
    GetFileSystemInfo()(FileSystemInfoable)
    GetMediaSource()(MediaSourceable)
    GetName()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDescription(value *string)()
    SetDriveItemSource(value DriveItemSourceable)()
    SetFileSize(value *int64)()
    SetFileSystemInfo(value FileSystemInfoable)()
    SetMediaSource(value MediaSourceable)()
    SetName(value *string)()
    SetOdataType(value *string)()
}
