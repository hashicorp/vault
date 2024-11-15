package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudPcGalleryImage struct {
    Entity
}
// NewCloudPcGalleryImage instantiates a new CloudPcGalleryImage and sets the default values.
func NewCloudPcGalleryImage()(*CloudPcGalleryImage) {
    m := &CloudPcGalleryImage{
        Entity: *NewEntity(),
    }
    return m
}
// CreateCloudPcGalleryImageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudPcGalleryImageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudPcGalleryImage(), nil
}
// GetDisplayName gets the displayName property value. The display name of this gallery image. For example, Windows 11 Enterprise + Microsoft 365 Apps 22H2. Read-only.
// returns a *string when successful
func (m *CloudPcGalleryImage) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEndDate gets the endDate property value. The date when the status of the image becomes supportedWithWarning. Users can still provision new Cloud PCs if the current time is later than endDate and earlier than expirationDate. For example, assume the endDate of a gallery image is 2023-9-14 and expirationDate is 2024-3-14, users are able to provision new Cloud PCs if today is 2023-10-01. Read-only.
// returns a *DateOnly when successful
func (m *CloudPcGalleryImage) GetEndDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("endDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// GetExpirationDate gets the expirationDate property value. The date when the image is no longer available. Users are unable to provision new Cloud PCs if the current time is later than expirationDate. The value is usually endDate plus six months. For example, if the startDate is 2025-10-14, the expirationDate is usually 2026-04-14. Read-only.
// returns a *DateOnly when successful
func (m *CloudPcGalleryImage) GetExpirationDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("expirationDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CloudPcGalleryImage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["endDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDate(val)
        }
        return nil
    }
    res["expirationDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpirationDate(val)
        }
        return nil
    }
    res["offerName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOfferName(val)
        }
        return nil
    }
    res["publisherName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisherName(val)
        }
        return nil
    }
    res["sizeInGB"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSizeInGB(val)
        }
        return nil
    }
    res["skuName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSkuName(val)
        }
        return nil
    }
    res["startDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetDateOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDate(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCloudPcGalleryImageStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*CloudPcGalleryImageStatus))
        }
        return nil
    }
    return res
}
// GetOfferName gets the offerName property value. The offer name of this gallery image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
// returns a *string when successful
func (m *CloudPcGalleryImage) GetOfferName()(*string) {
    val, err := m.GetBackingStore().Get("offerName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublisherName gets the publisherName property value. The publisher name of this gallery image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
// returns a *string when successful
func (m *CloudPcGalleryImage) GetPublisherName()(*string) {
    val, err := m.GetBackingStore().Get("publisherName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSizeInGB gets the sizeInGB property value. Indicates the size of this image in gigabytes. For example, 64. Read-only.
// returns a *int32 when successful
func (m *CloudPcGalleryImage) GetSizeInGB()(*int32) {
    val, err := m.GetBackingStore().Get("sizeInGB")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSkuName gets the skuName property value. The SKU name of this image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
// returns a *string when successful
func (m *CloudPcGalleryImage) GetSkuName()(*string) {
    val, err := m.GetBackingStore().Get("skuName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartDate gets the startDate property value. The date when the Cloud PC image is available for provisioning new Cloud PCs. For example, 2022-09-20. Read-only.
// returns a *DateOnly when successful
func (m *CloudPcGalleryImage) GetStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly) {
    val, err := m.GetBackingStore().Get("startDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    }
    return nil
}
// GetStatus gets the status property value. The status of the gallery image on the Cloud PC. Possible values are: supported, supportedWithWarning, notSupported, unknownFutureValue. The default value is supported. Read-only.
// returns a *CloudPcGalleryImageStatus when successful
func (m *CloudPcGalleryImage) GetStatus()(*CloudPcGalleryImageStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CloudPcGalleryImageStatus)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudPcGalleryImage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteDateOnlyValue("endDate", m.GetEndDate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteDateOnlyValue("expirationDate", m.GetExpirationDate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("offerName", m.GetOfferName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("publisherName", m.GetPublisherName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("sizeInGB", m.GetSizeInGB())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("skuName", m.GetSkuName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteDateOnlyValue("startDate", m.GetStartDate())
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
// SetDisplayName sets the displayName property value. The display name of this gallery image. For example, Windows 11 Enterprise + Microsoft 365 Apps 22H2. Read-only.
func (m *CloudPcGalleryImage) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEndDate sets the endDate property value. The date when the status of the image becomes supportedWithWarning. Users can still provision new Cloud PCs if the current time is later than endDate and earlier than expirationDate. For example, assume the endDate of a gallery image is 2023-9-14 and expirationDate is 2024-3-14, users are able to provision new Cloud PCs if today is 2023-10-01. Read-only.
func (m *CloudPcGalleryImage) SetEndDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("endDate", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDate sets the expirationDate property value. The date when the image is no longer available. Users are unable to provision new Cloud PCs if the current time is later than expirationDate. The value is usually endDate plus six months. For example, if the startDate is 2025-10-14, the expirationDate is usually 2026-04-14. Read-only.
func (m *CloudPcGalleryImage) SetExpirationDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("expirationDate", value)
    if err != nil {
        panic(err)
    }
}
// SetOfferName sets the offerName property value. The offer name of this gallery image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
func (m *CloudPcGalleryImage) SetOfferName(value *string)() {
    err := m.GetBackingStore().Set("offerName", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisherName sets the publisherName property value. The publisher name of this gallery image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
func (m *CloudPcGalleryImage) SetPublisherName(value *string)() {
    err := m.GetBackingStore().Set("publisherName", value)
    if err != nil {
        panic(err)
    }
}
// SetSizeInGB sets the sizeInGB property value. Indicates the size of this image in gigabytes. For example, 64. Read-only.
func (m *CloudPcGalleryImage) SetSizeInGB(value *int32)() {
    err := m.GetBackingStore().Set("sizeInGB", value)
    if err != nil {
        panic(err)
    }
}
// SetSkuName sets the skuName property value. The SKU name of this image that is passed to Azure Resource Manager (ARM) to retrieve the image resource. Read-only.
func (m *CloudPcGalleryImage) SetSkuName(value *string)() {
    err := m.GetBackingStore().Set("skuName", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDate sets the startDate property value. The date when the Cloud PC image is available for provisioning new Cloud PCs. For example, 2022-09-20. Read-only.
func (m *CloudPcGalleryImage) SetStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)() {
    err := m.GetBackingStore().Set("startDate", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of the gallery image on the Cloud PC. Possible values are: supported, supportedWithWarning, notSupported, unknownFutureValue. The default value is supported. Read-only.
func (m *CloudPcGalleryImage) SetStatus(value *CloudPcGalleryImageStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type CloudPcGalleryImageable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetEndDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    GetExpirationDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    GetOfferName()(*string)
    GetPublisherName()(*string)
    GetSizeInGB()(*int32)
    GetSkuName()(*string)
    GetStartDate()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)
    GetStatus()(*CloudPcGalleryImageStatus)
    SetDisplayName(value *string)()
    SetEndDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
    SetExpirationDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
    SetOfferName(value *string)()
    SetPublisherName(value *string)()
    SetSizeInGB(value *int32)()
    SetSkuName(value *string)()
    SetStartDate(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)()
    SetStatus(value *CloudPcGalleryImageStatus)()
}
