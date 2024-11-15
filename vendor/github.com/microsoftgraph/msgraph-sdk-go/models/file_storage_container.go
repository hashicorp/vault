package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FileStorageContainer struct {
    Entity
}
// NewFileStorageContainer instantiates a new FileStorageContainer and sets the default values.
func NewFileStorageContainer()(*FileStorageContainer) {
    m := &FileStorageContainer{
        Entity: *NewEntity(),
    }
    return m
}
// CreateFileStorageContainerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileStorageContainerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileStorageContainer(), nil
}
// GetContainerTypeId gets the containerTypeId property value. Container type ID of the fileStorageContainer. For details about container types, see Container Types. Each container must have only one container type. Read-only.
// returns a *UUID when successful
func (m *FileStorageContainer) GetContainerTypeId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("containerTypeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time of the fileStorageContainer creation. Read-only.
// returns a *Time when successful
func (m *FileStorageContainer) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCustomProperties gets the customProperties property value. Custom property collection for the fileStorageContainer. Read-write.
// returns a FileStorageContainerCustomPropertyDictionaryable when successful
func (m *FileStorageContainer) GetCustomProperties()(FileStorageContainerCustomPropertyDictionaryable) {
    val, err := m.GetBackingStore().Get("customProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileStorageContainerCustomPropertyDictionaryable)
    }
    return nil
}
// GetDescription gets the description property value. Provides a user-visible description of the fileStorageContainer. Read-write.
// returns a *string when successful
func (m *FileStorageContainer) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the fileStorageContainer. Read-write.
// returns a *string when successful
func (m *FileStorageContainer) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDrive gets the drive property value. The drive of the resource fileStorageContainer. Read-only.
// returns a Driveable when successful
func (m *FileStorageContainer) GetDrive()(Driveable) {
    val, err := m.GetBackingStore().Get("drive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Driveable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FileStorageContainer) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["containerTypeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContainerTypeId(val)
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
    res["customProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileStorageContainerCustomPropertyDictionaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomProperties(val.(FileStorageContainerCustomPropertyDictionaryable))
        }
        return nil
    }
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
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["drive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDriveFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDrive(val.(Driveable))
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
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFileStorageContainerStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*FileStorageContainerStatus))
        }
        return nil
    }
    res["viewpoint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFileStorageContainerViewpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetViewpoint(val.(FileStorageContainerViewpointable))
        }
        return nil
    }
    return res
}
// GetPermissions gets the permissions property value. The set of permissions for users in the fileStorageContainer. Permission for each user is set by the roles property. The possible values are: reader, writer, manager, and owner. Read-write.
// returns a []Permissionable when successful
func (m *FileStorageContainer) GetPermissions()([]Permissionable) {
    val, err := m.GetBackingStore().Get("permissions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Permissionable)
    }
    return nil
}
// GetStatus gets the status property value. Status of the fileStorageContainer. Containers are created as inactive and require activation. Inactive containers are subjected to automatic deletion in 24 hours. The possible values are: inactive, active. Read-only.
// returns a *FileStorageContainerStatus when successful
func (m *FileStorageContainer) GetStatus()(*FileStorageContainerStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FileStorageContainerStatus)
    }
    return nil
}
// GetViewpoint gets the viewpoint property value. Data specific to the current user. Read-only.
// returns a FileStorageContainerViewpointable when successful
func (m *FileStorageContainer) GetViewpoint()(FileStorageContainerViewpointable) {
    val, err := m.GetBackingStore().Get("viewpoint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FileStorageContainerViewpointable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FileStorageContainer) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteUUIDValue("containerTypeId", m.GetContainerTypeId())
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
        err = writer.WriteObjectValue("customProperties", m.GetCustomProperties())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("drive", m.GetDrive())
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
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("viewpoint", m.GetViewpoint())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContainerTypeId sets the containerTypeId property value. Container type ID of the fileStorageContainer. For details about container types, see Container Types. Each container must have only one container type. Read-only.
func (m *FileStorageContainer) SetContainerTypeId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("containerTypeId", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time of the fileStorageContainer creation. Read-only.
func (m *FileStorageContainer) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCustomProperties sets the customProperties property value. Custom property collection for the fileStorageContainer. Read-write.
func (m *FileStorageContainer) SetCustomProperties(value FileStorageContainerCustomPropertyDictionaryable)() {
    err := m.GetBackingStore().Set("customProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Provides a user-visible description of the fileStorageContainer. Read-write.
func (m *FileStorageContainer) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the fileStorageContainer. Read-write.
func (m *FileStorageContainer) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDrive sets the drive property value. The drive of the resource fileStorageContainer. Read-only.
func (m *FileStorageContainer) SetDrive(value Driveable)() {
    err := m.GetBackingStore().Set("drive", value)
    if err != nil {
        panic(err)
    }
}
// SetPermissions sets the permissions property value. The set of permissions for users in the fileStorageContainer. Permission for each user is set by the roles property. The possible values are: reader, writer, manager, and owner. Read-write.
func (m *FileStorageContainer) SetPermissions(value []Permissionable)() {
    err := m.GetBackingStore().Set("permissions", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Status of the fileStorageContainer. Containers are created as inactive and require activation. Inactive containers are subjected to automatic deletion in 24 hours. The possible values are: inactive, active. Read-only.
func (m *FileStorageContainer) SetStatus(value *FileStorageContainerStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetViewpoint sets the viewpoint property value. Data specific to the current user. Read-only.
func (m *FileStorageContainer) SetViewpoint(value FileStorageContainerViewpointable)() {
    err := m.GetBackingStore().Set("viewpoint", value)
    if err != nil {
        panic(err)
    }
}
type FileStorageContainerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContainerTypeId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCustomProperties()(FileStorageContainerCustomPropertyDictionaryable)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetDrive()(Driveable)
    GetPermissions()([]Permissionable)
    GetStatus()(*FileStorageContainerStatus)
    GetViewpoint()(FileStorageContainerViewpointable)
    SetContainerTypeId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCustomProperties(value FileStorageContainerCustomPropertyDictionaryable)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetDrive(value Driveable)()
    SetPermissions(value []Permissionable)()
    SetStatus(value *FileStorageContainerStatus)()
    SetViewpoint(value FileStorageContainerViewpointable)()
}
