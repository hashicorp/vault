package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DriveItem struct {
    BaseItem
}
// NewDriveItem instantiates a new DriveItem and sets the default values.
func NewDriveItem()(*DriveItem) {
    m := &DriveItem{
        BaseItem: *NewBaseItem(),
    }
    odataTypeValue := "#microsoft.graph.driveItem"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDriveItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDriveItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDriveItem(), nil
}
// GetAnalytics gets the analytics property value. Analytics about the view activities that took place on this item.
// returns a ItemAnalyticsable when successful
func (m *DriveItem) GetAnalytics()(ItemAnalyticsable) {
    val, err := m.GetBackingStore().Get("analytics")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemAnalyticsable)
    }
    return nil
}
// GetAudio gets the audio property value. Audio metadata, if the item is an audio file. Read-only. Read-only. Only on OneDrive Personal.
// returns a Audioable when successful
func (m *DriveItem) GetAudio()(Audioable) {
    val, err := m.GetBackingStore().Get("audio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Audioable)
    }
    return nil
}
// GetBundle gets the bundle property value. Bundle metadata, if the item is a bundle. Read-only.
// returns a Bundleable when successful
func (m *DriveItem) GetBundle()(Bundleable) {
    val, err := m.GetBackingStore().Get("bundle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Bundleable)
    }
    return nil
}
// GetChildren gets the children property value. Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable.
// returns a []DriveItemable when successful
func (m *DriveItem) GetChildren()([]DriveItemable) {
    val, err := m.GetBackingStore().Get("children")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemable)
    }
    return nil
}
// GetContent gets the content property value. The content stream, if the item represents a file.
// returns a []byte when successful
func (m *DriveItem) GetContent()([]byte) {
    val, err := m.GetBackingStore().Get("content")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetCTag gets the cTag property value. An eTag for the content of the item. This eTag isn't changed if only the metadata is changed. Note This property isn't returned if the item is a folder. Read-only.
// returns a *string when successful
func (m *DriveItem) GetCTag()(*string) {
    val, err := m.GetBackingStore().Get("cTag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeleted gets the deleted property value. Information about the deleted state of the item. Read-only.
// returns a Deletedable when successful
func (m *DriveItem) GetDeleted()(Deletedable) {
    val, err := m.GetBackingStore().Get("deleted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Deletedable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DriveItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseItem.GetFieldDeserializers()
    res["analytics"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemAnalyticsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAnalytics(val.(ItemAnalyticsable))
        }
        return nil
    }
    res["audio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAudioFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAudio(val.(Audioable))
        }
        return nil
    }
    res["bundle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBundleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBundle(val.(Bundleable))
        }
        return nil
    }
    res["children"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveItemable)
                }
            }
            m.SetChildren(res)
        }
        return nil
    }
    res["content"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContent(val)
        }
        return nil
    }
    res["cTag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCTag(val)
        }
        return nil
    }
    res["deleted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeletedFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeleted(val.(Deletedable))
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
    res["listItem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetListItem(val.(ListItemable))
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGeoCoordinatesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(GeoCoordinatesable))
        }
        return nil
    }
    res["malware"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMalwareFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMalware(val.(Malwareable))
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
    res["pendingOperations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePendingOperationsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingOperations(val.(PendingOperationsable))
        }
        return nil
    }
    res["permissions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePermissionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Permissionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Permissionable)
                }
            }
            m.SetPermissions(res)
        }
        return nil
    }
    res["photo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePhotoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoto(val.(Photoable))
        }
        return nil
    }
    res["publication"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePublicationFacetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublication(val.(PublicationFacetable))
        }
        return nil
    }
    res["remoteItem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRemoteItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoteItem(val.(RemoteItemable))
        }
        return nil
    }
    res["retentionLabel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemRetentionLabelFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRetentionLabel(val.(ItemRetentionLabelable))
        }
        return nil
    }
    res["root"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoot(val.(Rootable))
        }
        return nil
    }
    res["searchResult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSearchResultFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchResult(val.(SearchResultable))
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
    res["subscriptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSubscriptionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Subscriptionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Subscriptionable)
                }
            }
            m.SetSubscriptions(res)
        }
        return nil
    }
    res["thumbnails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateThumbnailSetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ThumbnailSetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ThumbnailSetable)
                }
            }
            m.SetThumbnails(res)
        }
        return nil
    }
    res["versions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveItemVersionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveItemVersionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveItemVersionable)
                }
            }
            m.SetVersions(res)
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
    res["workbook"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkbook(val.(Workbookable))
        }
        return nil
    }
    return res
}
// GetFile gets the file property value. File metadata, if the item is a file. Read-only.
// returns a Fileable when successful
func (m *DriveItem) GetFile()(Fileable) {
    val, err := m.GetBackingStore().Get("file")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Fileable)
    }
    return nil
}
// GetFileSystemInfo gets the fileSystemInfo property value. File system information on client. Read-write.
// returns a FileSystemInfoable when successful
func (m *DriveItem) GetFileSystemInfo()(FileSystemInfoable) {
    val, err := m.GetBackingStore().Get("fileSystemInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileSystemInfoable)
    }
    return nil
}
// GetFolder gets the folder property value. Folder metadata, if the item is a folder. Read-only.
// returns a Folderable when successful
func (m *DriveItem) GetFolder()(Folderable) {
    val, err := m.GetBackingStore().Get("folder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Folderable)
    }
    return nil
}
// GetImage gets the image property value. Image metadata, if the item is an image. Read-only.
// returns a Imageable when successful
func (m *DriveItem) GetImage()(Imageable) {
    val, err := m.GetBackingStore().Get("image")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Imageable)
    }
    return nil
}
// GetListItem gets the listItem property value. For drives in SharePoint, the associated document library list item. Read-only. Nullable.
// returns a ListItemable when successful
func (m *DriveItem) GetListItem()(ListItemable) {
    val, err := m.GetBackingStore().Get("listItem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ListItemable)
    }
    return nil
}
// GetLocation gets the location property value. Location metadata, if the item has location data. Read-only.
// returns a GeoCoordinatesable when successful
func (m *DriveItem) GetLocation()(GeoCoordinatesable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(GeoCoordinatesable)
    }
    return nil
}
// GetMalware gets the malware property value. Malware metadata, if the item was detected to contain malware. Read-only.
// returns a Malwareable when successful
func (m *DriveItem) GetMalware()(Malwareable) {
    val, err := m.GetBackingStore().Get("malware")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Malwareable)
    }
    return nil
}
// GetPackageEscaped gets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
// returns a PackageEscapedable when successful
func (m *DriveItem) GetPackageEscaped()(PackageEscapedable) {
    val, err := m.GetBackingStore().Get("packageEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PackageEscapedable)
    }
    return nil
}
// GetPendingOperations gets the pendingOperations property value. If present, indicates that one or more operations that might affect the state of the driveItem are pending completion. Read-only.
// returns a PendingOperationsable when successful
func (m *DriveItem) GetPendingOperations()(PendingOperationsable) {
    val, err := m.GetBackingStore().Get("pendingOperations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PendingOperationsable)
    }
    return nil
}
// GetPermissions gets the permissions property value. The set of permissions for the item. Read-only. Nullable.
// returns a []Permissionable when successful
func (m *DriveItem) GetPermissions()([]Permissionable) {
    val, err := m.GetBackingStore().Get("permissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Permissionable)
    }
    return nil
}
// GetPhoto gets the photo property value. Photo metadata, if the item is a photo. Read-only.
// returns a Photoable when successful
func (m *DriveItem) GetPhoto()(Photoable) {
    val, err := m.GetBackingStore().Get("photo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Photoable)
    }
    return nil
}
// GetPublication gets the publication property value. Provides information about the published or checked-out state of an item, in locations that support such actions. This property isn't returned by default. Read-only.
// returns a PublicationFacetable when successful
func (m *DriveItem) GetPublication()(PublicationFacetable) {
    val, err := m.GetBackingStore().Get("publication")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PublicationFacetable)
    }
    return nil
}
// GetRemoteItem gets the remoteItem property value. Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
// returns a RemoteItemable when successful
func (m *DriveItem) GetRemoteItem()(RemoteItemable) {
    val, err := m.GetBackingStore().Get("remoteItem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RemoteItemable)
    }
    return nil
}
// GetRetentionLabel gets the retentionLabel property value. Information about retention label and settings enforced on the driveItem. Read-write.
// returns a ItemRetentionLabelable when successful
func (m *DriveItem) GetRetentionLabel()(ItemRetentionLabelable) {
    val, err := m.GetBackingStore().Get("retentionLabel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemRetentionLabelable)
    }
    return nil
}
// GetRoot gets the root property value. If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.
// returns a Rootable when successful
func (m *DriveItem) GetRoot()(Rootable) {
    val, err := m.GetBackingStore().Get("root")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Rootable)
    }
    return nil
}
// GetSearchResult gets the searchResult property value. Search metadata, if the item is from a search result. Read-only.
// returns a SearchResultable when successful
func (m *DriveItem) GetSearchResult()(SearchResultable) {
    val, err := m.GetBackingStore().Get("searchResult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SearchResultable)
    }
    return nil
}
// GetShared gets the shared property value. Indicates that the item was shared with others and provides information about the shared state of the item. Read-only.
// returns a Sharedable when successful
func (m *DriveItem) GetShared()(Sharedable) {
    val, err := m.GetBackingStore().Get("shared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Sharedable)
    }
    return nil
}
// GetSharepointIds gets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
// returns a SharepointIdsable when successful
func (m *DriveItem) GetSharepointIds()(SharepointIdsable) {
    val, err := m.GetBackingStore().Get("sharepointIds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharepointIdsable)
    }
    return nil
}
// GetSize gets the size property value. Size of the item in bytes. Read-only.
// returns a *int64 when successful
func (m *DriveItem) GetSize()(*int64) {
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
func (m *DriveItem) GetSpecialFolder()(SpecialFolderable) {
    val, err := m.GetBackingStore().Get("specialFolder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SpecialFolderable)
    }
    return nil
}
// GetSubscriptions gets the subscriptions property value. The set of subscriptions on the item. Only supported on the root of a drive.
// returns a []Subscriptionable when successful
func (m *DriveItem) GetSubscriptions()([]Subscriptionable) {
    val, err := m.GetBackingStore().Get("subscriptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Subscriptionable)
    }
    return nil
}
// GetThumbnails gets the thumbnails property value. Collection of thumbnailSet objects associated with the item. For more information, see getting thumbnails. Read-only. Nullable.
// returns a []ThumbnailSetable when successful
func (m *DriveItem) GetThumbnails()([]ThumbnailSetable) {
    val, err := m.GetBackingStore().Get("thumbnails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ThumbnailSetable)
    }
    return nil
}
// GetVersions gets the versions property value. The list of previous versions of the item. For more info, see getting previous versions. Read-only. Nullable.
// returns a []DriveItemVersionable when successful
func (m *DriveItem) GetVersions()([]DriveItemVersionable) {
    val, err := m.GetBackingStore().Get("versions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveItemVersionable)
    }
    return nil
}
// GetVideo gets the video property value. Video metadata, if the item is a video. Read-only.
// returns a Videoable when successful
func (m *DriveItem) GetVideo()(Videoable) {
    val, err := m.GetBackingStore().Get("video")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Videoable)
    }
    return nil
}
// GetWebDavUrl gets the webDavUrl property value. WebDAV compatible URL for the item.
// returns a *string when successful
func (m *DriveItem) GetWebDavUrl()(*string) {
    val, err := m.GetBackingStore().Get("webDavUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWorkbook gets the workbook property value. For files that are Excel spreadsheets, access to the workbook API to work with the spreadsheet's contents. Nullable.
// returns a Workbookable when successful
func (m *DriveItem) GetWorkbook()(Workbookable) {
    val, err := m.GetBackingStore().Get("workbook")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Workbookable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DriveItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("analytics", m.GetAnalytics())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("audio", m.GetAudio())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("bundle", m.GetBundle())
        if err != nil {
            return err
        }
    }
    if m.GetChildren() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildren()))
        for i, v := range m.GetChildren() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("children", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("content", m.GetContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("cTag", m.GetCTag())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deleted", m.GetDeleted())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("file", m.GetFile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("fileSystemInfo", m.GetFileSystemInfo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("folder", m.GetFolder())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("image", m.GetImage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("listItem", m.GetListItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("malware", m.GetMalware())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("package", m.GetPackageEscaped())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("pendingOperations", m.GetPendingOperations())
        if err != nil {
            return err
        }
    }
    if m.GetPermissions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPermissions()))
        for i, v := range m.GetPermissions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("permissions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("publication", m.GetPublication())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("remoteItem", m.GetRemoteItem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("retentionLabel", m.GetRetentionLabel())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("root", m.GetRoot())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("searchResult", m.GetSearchResult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("shared", m.GetShared())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharepointIds", m.GetSharepointIds())
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
        err = writer.WriteObjectValue("specialFolder", m.GetSpecialFolder())
        if err != nil {
            return err
        }
    }
    if m.GetSubscriptions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSubscriptions()))
        for i, v := range m.GetSubscriptions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("subscriptions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetThumbnails() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetThumbnails()))
        for i, v := range m.GetThumbnails() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("thumbnails", cast)
        if err != nil {
            return err
        }
    }
    if m.GetVersions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetVersions()))
        for i, v := range m.GetVersions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("versions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("video", m.GetVideo())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webDavUrl", m.GetWebDavUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("workbook", m.GetWorkbook())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAnalytics sets the analytics property value. Analytics about the view activities that took place on this item.
func (m *DriveItem) SetAnalytics(value ItemAnalyticsable)() {
    err := m.GetBackingStore().Set("analytics", value)
    if err != nil {
        panic(err)
    }
}
// SetAudio sets the audio property value. Audio metadata, if the item is an audio file. Read-only. Read-only. Only on OneDrive Personal.
func (m *DriveItem) SetAudio(value Audioable)() {
    err := m.GetBackingStore().Set("audio", value)
    if err != nil {
        panic(err)
    }
}
// SetBundle sets the bundle property value. Bundle metadata, if the item is a bundle. Read-only.
func (m *DriveItem) SetBundle(value Bundleable)() {
    err := m.GetBackingStore().Set("bundle", value)
    if err != nil {
        panic(err)
    }
}
// SetChildren sets the children property value. Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable.
func (m *DriveItem) SetChildren(value []DriveItemable)() {
    err := m.GetBackingStore().Set("children", value)
    if err != nil {
        panic(err)
    }
}
// SetContent sets the content property value. The content stream, if the item represents a file.
func (m *DriveItem) SetContent(value []byte)() {
    err := m.GetBackingStore().Set("content", value)
    if err != nil {
        panic(err)
    }
}
// SetCTag sets the cTag property value. An eTag for the content of the item. This eTag isn't changed if only the metadata is changed. Note This property isn't returned if the item is a folder. Read-only.
func (m *DriveItem) SetCTag(value *string)() {
    err := m.GetBackingStore().Set("cTag", value)
    if err != nil {
        panic(err)
    }
}
// SetDeleted sets the deleted property value. Information about the deleted state of the item. Read-only.
func (m *DriveItem) SetDeleted(value Deletedable)() {
    err := m.GetBackingStore().Set("deleted", value)
    if err != nil {
        panic(err)
    }
}
// SetFile sets the file property value. File metadata, if the item is a file. Read-only.
func (m *DriveItem) SetFile(value Fileable)() {
    err := m.GetBackingStore().Set("file", value)
    if err != nil {
        panic(err)
    }
}
// SetFileSystemInfo sets the fileSystemInfo property value. File system information on client. Read-write.
func (m *DriveItem) SetFileSystemInfo(value FileSystemInfoable)() {
    err := m.GetBackingStore().Set("fileSystemInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetFolder sets the folder property value. Folder metadata, if the item is a folder. Read-only.
func (m *DriveItem) SetFolder(value Folderable)() {
    err := m.GetBackingStore().Set("folder", value)
    if err != nil {
        panic(err)
    }
}
// SetImage sets the image property value. Image metadata, if the item is an image. Read-only.
func (m *DriveItem) SetImage(value Imageable)() {
    err := m.GetBackingStore().Set("image", value)
    if err != nil {
        panic(err)
    }
}
// SetListItem sets the listItem property value. For drives in SharePoint, the associated document library list item. Read-only. Nullable.
func (m *DriveItem) SetListItem(value ListItemable)() {
    err := m.GetBackingStore().Set("listItem", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. Location metadata, if the item has location data. Read-only.
func (m *DriveItem) SetLocation(value GeoCoordinatesable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetMalware sets the malware property value. Malware metadata, if the item was detected to contain malware. Read-only.
func (m *DriveItem) SetMalware(value Malwareable)() {
    err := m.GetBackingStore().Set("malware", value)
    if err != nil {
        panic(err)
    }
}
// SetPackageEscaped sets the package property value. If present, indicates that this item is a package instead of a folder or file. Packages are treated like files in some contexts and folders in others. Read-only.
func (m *DriveItem) SetPackageEscaped(value PackageEscapedable)() {
    err := m.GetBackingStore().Set("packageEscaped", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingOperations sets the pendingOperations property value. If present, indicates that one or more operations that might affect the state of the driveItem are pending completion. Read-only.
func (m *DriveItem) SetPendingOperations(value PendingOperationsable)() {
    err := m.GetBackingStore().Set("pendingOperations", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissions sets the permissions property value. The set of permissions for the item. Read-only. Nullable.
func (m *DriveItem) SetPermissions(value []Permissionable)() {
    err := m.GetBackingStore().Set("permissions", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoto sets the photo property value. Photo metadata, if the item is a photo. Read-only.
func (m *DriveItem) SetPhoto(value Photoable)() {
    err := m.GetBackingStore().Set("photo", value)
    if err != nil {
        panic(err)
    }
}
// SetPublication sets the publication property value. Provides information about the published or checked-out state of an item, in locations that support such actions. This property isn't returned by default. Read-only.
func (m *DriveItem) SetPublication(value PublicationFacetable)() {
    err := m.GetBackingStore().Set("publication", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoteItem sets the remoteItem property value. Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
func (m *DriveItem) SetRemoteItem(value RemoteItemable)() {
    err := m.GetBackingStore().Set("remoteItem", value)
    if err != nil {
        panic(err)
    }
}
// SetRetentionLabel sets the retentionLabel property value. Information about retention label and settings enforced on the driveItem. Read-write.
func (m *DriveItem) SetRetentionLabel(value ItemRetentionLabelable)() {
    err := m.GetBackingStore().Set("retentionLabel", value)
    if err != nil {
        panic(err)
    }
}
// SetRoot sets the root property value. If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.
func (m *DriveItem) SetRoot(value Rootable)() {
    err := m.GetBackingStore().Set("root", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchResult sets the searchResult property value. Search metadata, if the item is from a search result. Read-only.
func (m *DriveItem) SetSearchResult(value SearchResultable)() {
    err := m.GetBackingStore().Set("searchResult", value)
    if err != nil {
        panic(err)
    }
}
// SetShared sets the shared property value. Indicates that the item was shared with others and provides information about the shared state of the item. Read-only.
func (m *DriveItem) SetShared(value Sharedable)() {
    err := m.GetBackingStore().Set("shared", value)
    if err != nil {
        panic(err)
    }
}
// SetSharepointIds sets the sharepointIds property value. Returns identifiers useful for SharePoint REST compatibility. Read-only.
func (m *DriveItem) SetSharepointIds(value SharepointIdsable)() {
    err := m.GetBackingStore().Set("sharepointIds", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. Size of the item in bytes. Read-only.
func (m *DriveItem) SetSize(value *int64)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetSpecialFolder sets the specialFolder property value. If the current item is also available as a special folder, this facet is returned. Read-only.
func (m *DriveItem) SetSpecialFolder(value SpecialFolderable)() {
    err := m.GetBackingStore().Set("specialFolder", value)
    if err != nil {
        panic(err)
    }
}
// SetSubscriptions sets the subscriptions property value. The set of subscriptions on the item. Only supported on the root of a drive.
func (m *DriveItem) SetSubscriptions(value []Subscriptionable)() {
    err := m.GetBackingStore().Set("subscriptions", value)
    if err != nil {
        panic(err)
    }
}
// SetThumbnails sets the thumbnails property value. Collection of thumbnailSet objects associated with the item. For more information, see getting thumbnails. Read-only. Nullable.
func (m *DriveItem) SetThumbnails(value []ThumbnailSetable)() {
    err := m.GetBackingStore().Set("thumbnails", value)
    if err != nil {
        panic(err)
    }
}
// SetVersions sets the versions property value. The list of previous versions of the item. For more info, see getting previous versions. Read-only. Nullable.
func (m *DriveItem) SetVersions(value []DriveItemVersionable)() {
    err := m.GetBackingStore().Set("versions", value)
    if err != nil {
        panic(err)
    }
}
// SetVideo sets the video property value. Video metadata, if the item is a video. Read-only.
func (m *DriveItem) SetVideo(value Videoable)() {
    err := m.GetBackingStore().Set("video", value)
    if err != nil {
        panic(err)
    }
}
// SetWebDavUrl sets the webDavUrl property value. WebDAV compatible URL for the item.
func (m *DriveItem) SetWebDavUrl(value *string)() {
    err := m.GetBackingStore().Set("webDavUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkbook sets the workbook property value. For files that are Excel spreadsheets, access to the workbook API to work with the spreadsheet's contents. Nullable.
func (m *DriveItem) SetWorkbook(value Workbookable)() {
    err := m.GetBackingStore().Set("workbook", value)
    if err != nil {
        panic(err)
    }
}
type DriveItemable interface {
    BaseItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAnalytics()(ItemAnalyticsable)
    GetAudio()(Audioable)
    GetBundle()(Bundleable)
    GetChildren()([]DriveItemable)
    GetContent()([]byte)
    GetCTag()(*string)
    GetDeleted()(Deletedable)
    GetFile()(Fileable)
    GetFileSystemInfo()(FileSystemInfoable)
    GetFolder()(Folderable)
    GetImage()(Imageable)
    GetListItem()(ListItemable)
    GetLocation()(GeoCoordinatesable)
    GetMalware()(Malwareable)
    GetPackageEscaped()(PackageEscapedable)
    GetPendingOperations()(PendingOperationsable)
    GetPermissions()([]Permissionable)
    GetPhoto()(Photoable)
    GetPublication()(PublicationFacetable)
    GetRemoteItem()(RemoteItemable)
    GetRetentionLabel()(ItemRetentionLabelable)
    GetRoot()(Rootable)
    GetSearchResult()(SearchResultable)
    GetShared()(Sharedable)
    GetSharepointIds()(SharepointIdsable)
    GetSize()(*int64)
    GetSpecialFolder()(SpecialFolderable)
    GetSubscriptions()([]Subscriptionable)
    GetThumbnails()([]ThumbnailSetable)
    GetVersions()([]DriveItemVersionable)
    GetVideo()(Videoable)
    GetWebDavUrl()(*string)
    GetWorkbook()(Workbookable)
    SetAnalytics(value ItemAnalyticsable)()
    SetAudio(value Audioable)()
    SetBundle(value Bundleable)()
    SetChildren(value []DriveItemable)()
    SetContent(value []byte)()
    SetCTag(value *string)()
    SetDeleted(value Deletedable)()
    SetFile(value Fileable)()
    SetFileSystemInfo(value FileSystemInfoable)()
    SetFolder(value Folderable)()
    SetImage(value Imageable)()
    SetListItem(value ListItemable)()
    SetLocation(value GeoCoordinatesable)()
    SetMalware(value Malwareable)()
    SetPackageEscaped(value PackageEscapedable)()
    SetPendingOperations(value PendingOperationsable)()
    SetPermissions(value []Permissionable)()
    SetPhoto(value Photoable)()
    SetPublication(value PublicationFacetable)()
    SetRemoteItem(value RemoteItemable)()
    SetRetentionLabel(value ItemRetentionLabelable)()
    SetRoot(value Rootable)()
    SetSearchResult(value SearchResultable)()
    SetShared(value Sharedable)()
    SetSharepointIds(value SharepointIdsable)()
    SetSize(value *int64)()
    SetSpecialFolder(value SpecialFolderable)()
    SetSubscriptions(value []Subscriptionable)()
    SetThumbnails(value []ThumbnailSetable)()
    SetVersions(value []DriveItemVersionable)()
    SetVideo(value Videoable)()
    SetWebDavUrl(value *string)()
    SetWorkbook(value Workbookable)()
}
