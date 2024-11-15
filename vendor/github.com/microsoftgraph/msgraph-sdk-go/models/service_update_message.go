package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ServiceUpdateMessage struct {
    ServiceAnnouncementBase
}
// NewServiceUpdateMessage instantiates a new ServiceUpdateMessage and sets the default values.
func NewServiceUpdateMessage()(*ServiceUpdateMessage) {
    m := &ServiceUpdateMessage{
        ServiceAnnouncementBase: *NewServiceAnnouncementBase(),
    }
    odataTypeValue := "#microsoft.graph.serviceUpdateMessage"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateServiceUpdateMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateServiceUpdateMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewServiceUpdateMessage(), nil
}
// GetActionRequiredByDateTime gets the actionRequiredByDateTime property value. The expected deadline of the action for the message.
// returns a *Time when successful
func (m *ServiceUpdateMessage) GetActionRequiredByDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("actionRequiredByDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetAttachments gets the attachments property value. A collection of serviceAnnouncementAttachments.
// returns a []ServiceAnnouncementAttachmentable when successful
func (m *ServiceUpdateMessage) GetAttachments()([]ServiceAnnouncementAttachmentable) {
    val, err := m.GetBackingStore().Get("attachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceAnnouncementAttachmentable)
    }
    return nil
}
// GetAttachmentsArchive gets the attachmentsArchive property value. The zip file that contains all attachments for a message.
// returns a []byte when successful
func (m *ServiceUpdateMessage) GetAttachmentsArchive()([]byte) {
    val, err := m.GetBackingStore().Get("attachmentsArchive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetBody gets the body property value. The body property
// returns a ItemBodyable when successful
func (m *ServiceUpdateMessage) GetBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("body")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetCategory gets the category property value. The category property
// returns a *ServiceUpdateCategory when successful
func (m *ServiceUpdateMessage) GetCategory()(*ServiceUpdateCategory) {
    val, err := m.GetBackingStore().Get("category")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceUpdateCategory)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ServiceUpdateMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ServiceAnnouncementBase.GetFieldDeserializers()
    res["actionRequiredByDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActionRequiredByDateTime(val)
        }
        return nil
    }
    res["attachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceAnnouncementAttachmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceAnnouncementAttachmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceAnnouncementAttachmentable)
                }
            }
            m.SetAttachments(res)
        }
        return nil
    }
    res["attachmentsArchive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttachmentsArchive(val)
        }
        return nil
    }
    res["body"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBody(val.(ItemBodyable))
        }
        return nil
    }
    res["category"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceUpdateCategory)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory(val.(*ServiceUpdateCategory))
        }
        return nil
    }
    res["hasAttachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasAttachments(val)
        }
        return nil
    }
    res["isMajorChange"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMajorChange(val)
        }
        return nil
    }
    res["services"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetServices(res)
        }
        return nil
    }
    res["severity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseServiceUpdateSeverity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeverity(val.(*ServiceUpdateSeverity))
        }
        return nil
    }
    res["tags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetTags(res)
        }
        return nil
    }
    res["viewPoint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateServiceUpdateMessageViewpointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetViewPoint(val.(ServiceUpdateMessageViewpointable))
        }
        return nil
    }
    return res
}
// GetHasAttachments gets the hasAttachments property value. Indicates whether the message has any attachment.
// returns a *bool when successful
func (m *ServiceUpdateMessage) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMajorChange gets the isMajorChange property value. Indicates whether the message describes a major update for the service.
// returns a *bool when successful
func (m *ServiceUpdateMessage) GetIsMajorChange()(*bool) {
    val, err := m.GetBackingStore().Get("isMajorChange")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetServices gets the services property value. The affected services by the service message.
// returns a []string when successful
func (m *ServiceUpdateMessage) GetServices()([]string) {
    val, err := m.GetBackingStore().Get("services")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSeverity gets the severity property value. The severity property
// returns a *ServiceUpdateSeverity when successful
func (m *ServiceUpdateMessage) GetSeverity()(*ServiceUpdateSeverity) {
    val, err := m.GetBackingStore().Get("severity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ServiceUpdateSeverity)
    }
    return nil
}
// GetTags gets the tags property value. A collection of tags for the service message. Tags are provided by the service team/support team who post the message to tell whether this message contains privacy data, or whether this message is for a service new feature update, and so on.
// returns a []string when successful
func (m *ServiceUpdateMessage) GetTags()([]string) {
    val, err := m.GetBackingStore().Get("tags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetViewPoint gets the viewPoint property value. Represents user viewpoints data of the service message. This data includes message status such as whether the user has archived, read, or marked the message as favorite. This property is null when accessed with application permissions.
// returns a ServiceUpdateMessageViewpointable when successful
func (m *ServiceUpdateMessage) GetViewPoint()(ServiceUpdateMessageViewpointable) {
    val, err := m.GetBackingStore().Get("viewPoint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ServiceUpdateMessageViewpointable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ServiceUpdateMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ServiceAnnouncementBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("actionRequiredByDateTime", m.GetActionRequiredByDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetAttachments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttachments()))
        for i, v := range m.GetAttachments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("attachments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("attachmentsArchive", m.GetAttachmentsArchive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("body", m.GetBody())
        if err != nil {
            return err
        }
    }
    if m.GetCategory() != nil {
        cast := (*m.GetCategory()).String()
        err = writer.WriteStringValue("category", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isMajorChange", m.GetIsMajorChange())
        if err != nil {
            return err
        }
    }
    if m.GetServices() != nil {
        err = writer.WriteCollectionOfStringValues("services", m.GetServices())
        if err != nil {
            return err
        }
    }
    if m.GetSeverity() != nil {
        cast := (*m.GetSeverity()).String()
        err = writer.WriteStringValue("severity", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetTags() != nil {
        err = writer.WriteCollectionOfStringValues("tags", m.GetTags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("viewPoint", m.GetViewPoint())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActionRequiredByDateTime sets the actionRequiredByDateTime property value. The expected deadline of the action for the message.
func (m *ServiceUpdateMessage) SetActionRequiredByDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("actionRequiredByDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetAttachments sets the attachments property value. A collection of serviceAnnouncementAttachments.
func (m *ServiceUpdateMessage) SetAttachments(value []ServiceAnnouncementAttachmentable)() {
    err := m.GetBackingStore().Set("attachments", value)
    if err != nil {
        panic(err)
    }
}
// SetAttachmentsArchive sets the attachmentsArchive property value. The zip file that contains all attachments for a message.
func (m *ServiceUpdateMessage) SetAttachmentsArchive(value []byte)() {
    err := m.GetBackingStore().Set("attachmentsArchive", value)
    if err != nil {
        panic(err)
    }
}
// SetBody sets the body property value. The body property
func (m *ServiceUpdateMessage) SetBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("body", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory sets the category property value. The category property
func (m *ServiceUpdateMessage) SetCategory(value *ServiceUpdateCategory)() {
    err := m.GetBackingStore().Set("category", value)
    if err != nil {
        panic(err)
    }
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether the message has any attachment.
func (m *ServiceUpdateMessage) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMajorChange sets the isMajorChange property value. Indicates whether the message describes a major update for the service.
func (m *ServiceUpdateMessage) SetIsMajorChange(value *bool)() {
    err := m.GetBackingStore().Set("isMajorChange", value)
    if err != nil {
        panic(err)
    }
}
// SetServices sets the services property value. The affected services by the service message.
func (m *ServiceUpdateMessage) SetServices(value []string)() {
    err := m.GetBackingStore().Set("services", value)
    if err != nil {
        panic(err)
    }
}
// SetSeverity sets the severity property value. The severity property
func (m *ServiceUpdateMessage) SetSeverity(value *ServiceUpdateSeverity)() {
    err := m.GetBackingStore().Set("severity", value)
    if err != nil {
        panic(err)
    }
}
// SetTags sets the tags property value. A collection of tags for the service message. Tags are provided by the service team/support team who post the message to tell whether this message contains privacy data, or whether this message is for a service new feature update, and so on.
func (m *ServiceUpdateMessage) SetTags(value []string)() {
    err := m.GetBackingStore().Set("tags", value)
    if err != nil {
        panic(err)
    }
}
// SetViewPoint sets the viewPoint property value. Represents user viewpoints data of the service message. This data includes message status such as whether the user has archived, read, or marked the message as favorite. This property is null when accessed with application permissions.
func (m *ServiceUpdateMessage) SetViewPoint(value ServiceUpdateMessageViewpointable)() {
    err := m.GetBackingStore().Set("viewPoint", value)
    if err != nil {
        panic(err)
    }
}
type ServiceUpdateMessageable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ServiceAnnouncementBaseable
    GetActionRequiredByDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetAttachments()([]ServiceAnnouncementAttachmentable)
    GetAttachmentsArchive()([]byte)
    GetBody()(ItemBodyable)
    GetCategory()(*ServiceUpdateCategory)
    GetHasAttachments()(*bool)
    GetIsMajorChange()(*bool)
    GetServices()([]string)
    GetSeverity()(*ServiceUpdateSeverity)
    GetTags()([]string)
    GetViewPoint()(ServiceUpdateMessageViewpointable)
    SetActionRequiredByDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetAttachments(value []ServiceAnnouncementAttachmentable)()
    SetAttachmentsArchive(value []byte)()
    SetBody(value ItemBodyable)()
    SetCategory(value *ServiceUpdateCategory)()
    SetHasAttachments(value *bool)()
    SetIsMajorChange(value *bool)()
    SetServices(value []string)()
    SetSeverity(value *ServiceUpdateSeverity)()
    SetTags(value []string)()
    SetViewPoint(value ServiceUpdateMessageViewpointable)()
}
