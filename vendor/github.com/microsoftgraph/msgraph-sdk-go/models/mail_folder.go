package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailFolder struct {
    Entity
}
// NewMailFolder instantiates a new MailFolder and sets the default values.
func NewMailFolder()(*MailFolder) {
    m := &MailFolder{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMailFolderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailFolderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.mailSearchFolder":
                        return NewMailSearchFolder(), nil
                }
            }
        }
    }
    return NewMailFolder(), nil
}
// GetChildFolderCount gets the childFolderCount property value. The number of immediate child mailFolders in the current mailFolder.
// returns a *int32 when successful
func (m *MailFolder) GetChildFolderCount()(*int32) {
    val, err := m.GetBackingStore().Get("childFolderCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetChildFolders gets the childFolders property value. The collection of child folders in the mailFolder.
// returns a []MailFolderable when successful
func (m *MailFolder) GetChildFolders()([]MailFolderable) {
    val, err := m.GetBackingStore().Get("childFolders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailFolderable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The mailFolder's display name.
// returns a *string when successful
func (m *MailFolder) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
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
func (m *MailFolder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["childFolderCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChildFolderCount(val)
        }
        return nil
    }
    res["childFolders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailFolderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailFolderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailFolderable)
                }
            }
            m.SetChildFolders(res)
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
    res["isHidden"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsHidden(val)
        }
        return nil
    }
    res["messageRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMessageRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MessageRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MessageRuleable)
                }
            }
            m.SetMessageRules(res)
        }
        return nil
    }
    res["messages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Messageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Messageable)
                }
            }
            m.SetMessages(res)
        }
        return nil
    }
    res["multiValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MultiValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MultiValueLegacyExtendedPropertyable)
                }
            }
            m.SetMultiValueExtendedProperties(res)
        }
        return nil
    }
    res["parentFolderId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentFolderId(val)
        }
        return nil
    }
    res["singleValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SingleValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SingleValueLegacyExtendedPropertyable)
                }
            }
            m.SetSingleValueExtendedProperties(res)
        }
        return nil
    }
    res["totalItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalItemCount(val)
        }
        return nil
    }
    res["unreadItemCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnreadItemCount(val)
        }
        return nil
    }
    return res
}
// GetIsHidden gets the isHidden property value. Indicates whether the mailFolder is hidden. This property can be set only when creating the folder. Find more information in Hidden mail folders.
// returns a *bool when successful
func (m *MailFolder) GetIsHidden()(*bool) {
    val, err := m.GetBackingStore().Get("isHidden")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMessageRules gets the messageRules property value. The collection of rules that apply to the user's Inbox folder.
// returns a []MessageRuleable when successful
func (m *MailFolder) GetMessageRules()([]MessageRuleable) {
    val, err := m.GetBackingStore().Get("messageRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MessageRuleable)
    }
    return nil
}
// GetMessages gets the messages property value. The collection of messages in the mailFolder.
// returns a []Messageable when successful
func (m *MailFolder) GetMessages()([]Messageable) {
    val, err := m.GetBackingStore().Get("messages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Messageable)
    }
    return nil
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the mailFolder. Read-only. Nullable.
// returns a []MultiValueLegacyExtendedPropertyable when successful
func (m *MailFolder) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("multiValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MultiValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetParentFolderId gets the parentFolderId property value. The unique identifier for the mailFolder's parent mailFolder.
// returns a *string when successful
func (m *MailFolder) GetParentFolderId()(*string) {
    val, err := m.GetBackingStore().Get("parentFolderId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the mailFolder. Read-only. Nullable.
// returns a []SingleValueLegacyExtendedPropertyable when successful
func (m *MailFolder) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("singleValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SingleValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetTotalItemCount gets the totalItemCount property value. The number of items in the mailFolder.
// returns a *int32 when successful
func (m *MailFolder) GetTotalItemCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnreadItemCount gets the unreadItemCount property value. The number of items in the mailFolder marked as unread.
// returns a *int32 when successful
func (m *MailFolder) GetUnreadItemCount()(*int32) {
    val, err := m.GetBackingStore().Get("unreadItemCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailFolder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("childFolderCount", m.GetChildFolderCount())
        if err != nil {
            return err
        }
    }
    if m.GetChildFolders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildFolders()))
        for i, v := range m.GetChildFolders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("childFolders", cast)
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
        err = writer.WriteBoolValue("isHidden", m.GetIsHidden())
        if err != nil {
            return err
        }
    }
    if m.GetMessageRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessageRules()))
        for i, v := range m.GetMessageRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messageRules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessages()))
        for i, v := range m.GetMessages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messages", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMultiValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMultiValueExtendedProperties()))
        for i, v := range m.GetMultiValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("multiValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentFolderId", m.GetParentFolderId())
        if err != nil {
            return err
        }
    }
    if m.GetSingleValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSingleValueExtendedProperties()))
        for i, v := range m.GetSingleValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("singleValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalItemCount", m.GetTotalItemCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unreadItemCount", m.GetUnreadItemCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChildFolderCount sets the childFolderCount property value. The number of immediate child mailFolders in the current mailFolder.
func (m *MailFolder) SetChildFolderCount(value *int32)() {
    err := m.GetBackingStore().Set("childFolderCount", value)
    if err != nil {
        panic(err)
    }
}
// SetChildFolders sets the childFolders property value. The collection of child folders in the mailFolder.
func (m *MailFolder) SetChildFolders(value []MailFolderable)() {
    err := m.GetBackingStore().Set("childFolders", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The mailFolder's display name.
func (m *MailFolder) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsHidden sets the isHidden property value. Indicates whether the mailFolder is hidden. This property can be set only when creating the folder. Find more information in Hidden mail folders.
func (m *MailFolder) SetIsHidden(value *bool)() {
    err := m.GetBackingStore().Set("isHidden", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageRules sets the messageRules property value. The collection of rules that apply to the user's Inbox folder.
func (m *MailFolder) SetMessageRules(value []MessageRuleable)() {
    err := m.GetBackingStore().Set("messageRules", value)
    if err != nil {
        panic(err)
    }
}
// SetMessages sets the messages property value. The collection of messages in the mailFolder.
func (m *MailFolder) SetMessages(value []Messageable)() {
    err := m.GetBackingStore().Set("messages", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("multiValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetParentFolderId sets the parentFolderId property value. The unique identifier for the mailFolder's parent mailFolder.
func (m *MailFolder) SetParentFolderId(value *string)() {
    err := m.GetBackingStore().Set("parentFolderId", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the mailFolder. Read-only. Nullable.
func (m *MailFolder) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("singleValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalItemCount sets the totalItemCount property value. The number of items in the mailFolder.
func (m *MailFolder) SetTotalItemCount(value *int32)() {
    err := m.GetBackingStore().Set("totalItemCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnreadItemCount sets the unreadItemCount property value. The number of items in the mailFolder marked as unread.
func (m *MailFolder) SetUnreadItemCount(value *int32)() {
    err := m.GetBackingStore().Set("unreadItemCount", value)
    if err != nil {
        panic(err)
    }
}
type MailFolderable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChildFolderCount()(*int32)
    GetChildFolders()([]MailFolderable)
    GetDisplayName()(*string)
    GetIsHidden()(*bool)
    GetMessageRules()([]MessageRuleable)
    GetMessages()([]Messageable)
    GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable)
    GetParentFolderId()(*string)
    GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable)
    GetTotalItemCount()(*int32)
    GetUnreadItemCount()(*int32)
    SetChildFolderCount(value *int32)()
    SetChildFolders(value []MailFolderable)()
    SetDisplayName(value *string)()
    SetIsHidden(value *bool)()
    SetMessageRules(value []MessageRuleable)()
    SetMessages(value []Messageable)()
    SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)()
    SetParentFolderId(value *string)()
    SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)()
    SetTotalItemCount(value *int32)()
    SetUnreadItemCount(value *int32)()
}
