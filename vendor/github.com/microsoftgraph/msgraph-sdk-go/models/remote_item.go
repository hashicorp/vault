package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type RemoteItem struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewRemoteItem instantiates a new RemoteItem and sets the default values.
func NewRemoteItem()(*RemoteItem) {
    m := &RemoteItem{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateRemoteItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRemoteItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRemoteItem(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *RemoteItem) GetAdditionalData()(map[string]any) {
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
func (m *RemoteItem) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCreatedBy gets the createdBy property value. Identity of the user, device, and application which created the item. Read-only.
// returns a IdentitySetable when successful
func (m *RemoteItem) GetCreatedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time of item creation. Read-only.
// returns a *Time when successful
func (m *RemoteItem) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *RemoteItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["file"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFile(val.(Fileable))
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
    res["folder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFolderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFolder(val.(Folderable))
        }
        return nil
    }
    res["id"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetId(val)
        }
        return nil
    }
    res["image"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateImageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImage(val.(Imageable))
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(IdentitySetable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
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
    res["package"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePackageEscapedFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPackageEscaped(val.(PackageEscapedable))
        }
        return nil
    }
    res["parentReference"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentReference(val.(ItemReferenceable))
        }
        return nil
    }
    res["shared"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharedFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShared(val.(Sharedable))
        }
        return nil
    }
    res["sharepointIds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharepointIdsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharepointIds(val.(SharepointIdsable))
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
    res["specialFolder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSpecialFolderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpecialFolder(val.(SpecialFolderable))
        }
        return nil
    }
    res["video"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVideoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVideo(val.(Videoable))
        }
        return nil
    }
    res["webDavUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebDavUrl(val)
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetFile gets the file property value. Indicates that the remote item is a file. Read-only.
// returns a Fileable when successful
func (m *RemoteItem) GetFile()(Fileable) {
    val, err := m.GetBackingStore().Get("file")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Fileable)
    }
    return nil
}
// GetFileSystemInfo gets the fileSystemInfo property value. Information about the remote item from the local file system. Read-only.
// returns a FileSystemInfoable when successful
func (m *RemoteItem) GetFileSystemInfo()(FileSystemInfoable) {
    val, err := m.GetBackingStore().Get("fileSystemInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileSystemInfoable)
    }
    return nil
}
// GetFolder gets the folder property value. Indicates that the remote item is a folder. Read-only.
// returns a Folderable when successful
func (m *RemoteItem) GetFolder()(Folderable) {
    val, err := m.GetBackingStore().Get("folder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Folderable)
    }
    return nil
}
// GetId gets the id property value. Unique identifier for the remote item in its drive. Read-only.
// returns a *string when successful
func (m *RemoteItem) GetId()(*string) {
    val, err := m.GetBackingStore().Get("id")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetImage gets the image property value. Image metadata, if the item is an image. Read-only.
// returns a Imageable when successful
func (m *RemoteItem) GetImage()(Imageable) {
    val, err := m.GetBackingStore().Get("image")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Imageable)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. Identity of the user, device, and application which last modified the item. Read-only.
// returns a IdentitySetable when successful
func (m *RemoteItem) GetLastModifiedBy()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Date and time the item was last modified. Read-only.
// returns a *Time when successful
func (m *RemoteItem) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetName gets the name property value. Optional. Filename of the remote item. Read-only.
// returns a *string when successful
func (m *RemoteItem) GetName()(*string) {
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
func (m *RemoteItem) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPackageEscaped gets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
// returns a PackageEscapedable when successful
func (m *RemoteItem) GetPackageEscaped()(PackageEscapedable) {
    val, err := m.GetBackingStore().Get("packageEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PackageEscapedable)
    }
    return nil
}
// GetParentReference gets the parentReference property value. Properties of the parent of the remote item. Read-only.
// returns a ItemReferenceable when successful
func (m *RemoteItem) GetParentReference()(ItemReferenceable) {
    val, err := m.GetBackingStore().Get("parentReference")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemReferenceable)
    }
    return nil
}
// GetShared gets the shared property value. Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
// returns a Sharedable when successful
func (m *RemoteItem) GetShared()(Sharedable) {
    val, err := m.GetBackingStore().Get("shared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Sharedable)
    }
    return nil
}
// GetSharepointIds gets the sharepointIds property value. Provides interop between items in OneDrive for Business and SharePoint with the full set of item identifiers. Read-only.
// returns a SharepointIdsable when successful
func (m *RemoteItem) GetSharepointIds()(SharepointIdsable) {
    val, err := m.GetBackingStore().Get("sharepointIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharepointIdsable)
    }
    return nil
}
// GetSize gets the size property value. Size of the remote item. Read-only.
// returns a *int64 when successful
func (m *RemoteItem) GetSize()(*int64) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetSpecialFolder gets the specialFolder property value. If the current item is also available as a special folder, this facet is returned. Read-only.
// returns a SpecialFolderable when successful
func (m *RemoteItem) GetSpecialFolder()(SpecialFolderable) {
    val, err := m.GetBackingStore().Get("specialFolder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SpecialFolderable)
    }
    return nil
}
// GetVideo gets the video property value. Video metadata, if the item is a video. Read-only.
// returns a Videoable when successful
func (m *RemoteItem) GetVideo()(Videoable) {
    val, err := m.GetBackingStore().Get("video")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Videoable)
    }
    return nil
}
// GetWebDavUrl gets the webDavUrl property value. DAV compatible URL for the item.
// returns a *string when successful
func (m *RemoteItem) GetWebDavUrl()(*string) {
    val, err := m.GetBackingStore().Get("webDavUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. URL that displays the resource in the browser. Read-only.
// returns a *string when successful
func (m *RemoteItem) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RemoteItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("file", m.GetFile())
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
        err := writer.WriteObjectValue("folder", m.GetFolder())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("id", m.GetId())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("image", m.GetImage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
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
        err := writer.WriteObjectValue("package", m.GetPackageEscaped())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("parentReference", m.GetParentReference())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("shared", m.GetShared())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("specialFolder", m.GetSpecialFolder())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("video", m.GetVideo())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("webDavUrl", m.GetWebDavUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("webUrl", m.GetWebUrl())
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
func (m *RemoteItem) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *RemoteItem) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCreatedBy sets the createdBy property value. Identity of the user, device, and application which created the item. Read-only.
func (m *RemoteItem) SetCreatedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time of item creation. Read-only.
func (m *RemoteItem) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetFile sets the file property value. Indicates that the remote item is a file. Read-only.
func (m *RemoteItem) SetFile(value Fileable)() {
    err := m.GetBackingStore().Set("file", value)
    if err != nil {
        panic(err)
    }
}
// SetFileSystemInfo sets the fileSystemInfo property value. Information about the remote item from the local file system. Read-only.
func (m *RemoteItem) SetFileSystemInfo(value FileSystemInfoable)() {
    err := m.GetBackingStore().Set("fileSystemInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetFolder sets the folder property value. Indicates that the remote item is a folder. Read-only.
func (m *RemoteItem) SetFolder(value Folderable)() {
    err := m.GetBackingStore().Set("folder", value)
    if err != nil {
        panic(err)
    }
}
// SetId sets the id property value. Unique identifier for the remote item in its drive. Read-only.
func (m *RemoteItem) SetId(value *string)() {
    err := m.GetBackingStore().Set("id", value)
    if err != nil {
        panic(err)
    }
}
// SetImage sets the image property value. Image metadata, if the item is an image. Read-only.
func (m *RemoteItem) SetImage(value Imageable)() {
    err := m.GetBackingStore().Set("image", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. Identity of the user, device, and application which last modified the item. Read-only.
func (m *RemoteItem) SetLastModifiedBy(value IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Date and time the item was last modified. Read-only.
func (m *RemoteItem) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Optional. Filename of the remote item. Read-only.
func (m *RemoteItem) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *RemoteItem) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPackageEscaped sets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
func (m *RemoteItem) SetPackageEscaped(value PackageEscapedable)() {
    err := m.GetBackingStore().Set("packageEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetParentReference sets the parentReference property value. Properties of the parent of the remote item. Read-only.
func (m *RemoteItem) SetParentReference(value ItemReferenceable)() {
    err := m.GetBackingStore().Set("parentReference", value)
    if err != nil {
        panic(err)
    }
}
// SetShared sets the shared property value. Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
func (m *RemoteItem) SetShared(value Sharedable)() {
    err := m.GetBackingStore().Set("shared", value)
    if err != nil {
        panic(err)
    }
}
// SetSharepointIds sets the sharepointIds property value. Provides interop between items in OneDrive for Business and SharePoint with the full set of item identifiers. Read-only.
func (m *RemoteItem) SetSharepointIds(value SharepointIdsable)() {
    err := m.GetBackingStore().Set("sharepointIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. Size of the remote item. Read-only.
func (m *RemoteItem) SetSize(value *int64)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetSpecialFolder sets the specialFolder property value. If the current item is also available as a special folder, this facet is returned. Read-only.
func (m *RemoteItem) SetSpecialFolder(value SpecialFolderable)() {
    err := m.GetBackingStore().Set("specialFolder", value)
    if err != nil {
        panic(err)
    }
}
// SetVideo sets the video property value. Video metadata, if the item is a video. Read-only.
func (m *RemoteItem) SetVideo(value Videoable)() {
    err := m.GetBackingStore().Set("video", value)
    if err != nil {
        panic(err)
    }
}
// SetWebDavUrl sets the webDavUrl property value. DAV compatible URL for the item.
func (m *RemoteItem) SetWebDavUrl(value *string)() {
    err := m.GetBackingStore().Set("webDavUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. URL that displays the resource in the browser. Read-only.
func (m *RemoteItem) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type RemoteItemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCreatedBy()(IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetFile()(Fileable)
    GetFileSystemInfo()(FileSystemInfoable)
    GetFolder()(Folderable)
    GetId()(*string)
    GetImage()(Imageable)
    GetLastModifiedBy()(IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetName()(*string)
    GetOdataType()(*string)
    GetPackageEscaped()(PackageEscapedable)
    GetParentReference()(ItemReferenceable)
    GetShared()(Sharedable)
    GetSharepointIds()(SharepointIdsable)
    GetSize()(*int64)
    GetSpecialFolder()(SpecialFolderable)
    GetVideo()(Videoable)
    GetWebDavUrl()(*string)
    GetWebUrl()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCreatedBy(value IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetFile(value Fileable)()
    SetFileSystemInfo(value FileSystemInfoable)()
    SetFolder(value Folderable)()
    SetId(value *string)()
    SetImage(value Imageable)()
    SetLastModifiedBy(value IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetName(value *string)()
    SetOdataType(value *string)()
    SetPackageEscaped(value PackageEscapedable)()
    SetParentReference(value ItemReferenceable)()
    SetShared(value Sharedable)()
    SetSharepointIds(value SharepointIdsable)()
    SetSize(value *int64)()
    SetSpecialFolder(value SpecialFolderable)()
    SetVideo(value Videoable)()
    SetWebDavUrl(value *string)()
    SetWebUrl(value *string)()
}
